package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/telehash/interoper/pkg/test"
)

type Execer interface {
	Exec(wdir, name string, proc *test.Process) (out <-chan Event, kill chan struct{}, done <-chan struct{})
}

type ID uint64
type Info map[string]interface{}

type Event struct {
	ID   ID     `json:"id"`
	Type string `json:"ty"`
	Time int64  `json:"ti"`
	Role string `json:"ro"`
	Info Info   `json:"in,omitempty"`
}

func (e *Event) UnmarshalJSON(p []byte) error {
	var input struct {
		ID   ID     `json:"id"`
		Type string `json:"ty"`
		Info Info   `json:"in,omitempty"`
	}

	err := json.Unmarshal(p, &input)
	if err != nil {
		return err
	}

	e.Type = input.Type
	e.ID = input.ID
	e.Info = input.Info
	e.RecordTime()

	return nil
}

func (e *Event) RecordTime() {
	if e != nil && e.Time == 0 {
		e.Time = time.Now().UTC().UnixNano() / (1000 * 1000)
	}
}

func (e *Event) AssertParams() (id int, value string) {
	if e.Type == "assert" && e.Info != nil {
		v1, ok1 := e.Info["assert_id"].(int)
		v2, ok2 := e.Info["value"].(string)
		if ok1 && ok2 {
			return v1, v2
		}
	}
	return
}

func (e *Event) ExitCode() (code int) {
	if e.Info != nil {
		v, ok := e.Info["exit_code"].(int)
		if ok {
			return v
		}
	}
	return
}

func (e *Event) Success() bool {
	if e.Info != nil {
		v, ok := e.Info["success"].(bool)
		if ok {
			return v
		}
	}
	return false
}

type Context struct {
	SUT     string
	Drivers []string
	Execer  Execer
	Report  *Report
}

type Report struct {
	SUT   string                 `json:"sut"`
	Tests map[string]*TestReport `json:"tests"`
}

type TestReport struct {
	Description *test.Description     `json:"desc"`
	Runs        map[string]*RunReport `json:"runs"`
}

type RunReport struct {
	Driver string  `json:"driver"`
	Log    []Event `json:"log"`
}

func Run(t *test.Description, ctx *Context) {
	fmt.Printf("--> %s\n", t.Name)

	if ctx.Report == nil {
		ctx.Report = &Report{SUT: ctx.SUT, Tests: make(map[string]*TestReport)}
	}

	testReport := ctx.Report.Tests[t.Name]
	if testReport == nil {
		testReport = &TestReport{Description: t, Runs: make(map[string]*RunReport)}
		ctx.Report.Tests[t.Name] = testReport
	}

	for _, driver := range ctx.Drivers {
		var success bool
		var report = &RunReport{}
		report.Driver = driver
		testReport.Runs[driver] = report

		fmt.Printf("    %s%s", driver, strings.Repeat(" ", 10-len(driver)))

		r := newRunner(t, ctx, ctx.SUT, driver)
		for event := range r.run() {
			if event.Type == "status" {
				success = event.Success()
			}
			report.Log = append(report.Log, event)
		}

		if success {
			fmt.Println(" DONE")
		} else {
			fmt.Println(" FAIL")
		}
	}
}

type runner struct {
	mtx     sync.Mutex
	timeout time.Duration
	execer  Execer
	wdir    string
	sut     string
	driver  string
	procs   map[string]*test.Process
	kills   map[string]chan struct{}
	stopped map[string]<-chan struct{}
}

func newRunner(t *test.Description, ctx *Context, sut, driver string) *runner {
	r := &runner{
		timeout: t.Timeout(),
		execer:  ctx.Execer,
		wdir:    path.Join(os.TempDir(), fmt.Sprintf("interop-%d", time.Now().UnixNano())),
		sut:     sut,
		driver:  driver,
		procs:   make(map[string]*test.Process),
		kills:   make(map[string]chan struct{}),
		stopped: make(map[string]<-chan struct{}),
	}

	for name, proc := range t.Containers {
		if name == "sut" {
			r.procs[name] = proc.ForImplementation(sut)
		} else {
			r.procs[name] = proc.ForImplementation(driver)
		}
	}

	return r
}

func (r *runner) kill(role string) (ok bool) {
	if role == "*" {
		r.mtx.Lock()
		roles := make([]string, 0, len(r.kills))
		for role := range r.kills {
			roles = append(roles, role)
		}
		r.mtx.Unlock()

		for _, role := range roles {
			r.kill(role)
		}
		return true
	}

	r.mtx.Lock()
	cKill := r.kills[role]
	r.mtx.Unlock()

	if cKill != nil {
		defer func() { recover() }()
		close(cKill)
		ok = true
	}

	return
}

func (r *runner) close(role string) (ok bool) {
	r.mtx.Lock()
	in := r.kills[role]
	delete(r.kills, role)
	r.mtx.Unlock()

	if in != nil {
		func() {
			defer func() { recover() }()
			close(in)
			ok = true
		}()
	}

	return
}

func (r *runner) controller(in <-chan Event) <-chan Event {
	out := make(chan Event)

	go func() {
		var (
			asserts = map[int][]Event{}
			failed  bool
			done    bool
		)

		defer func() { close(out) }()

		for event := range in {
			out <- event

			if event.Type == "done" && event.Role == "driver" {
				done = true
			}

			if event.Type == "exited" {
				r.kill("*")
			}

			if event.Type == "assert" {
				id, value := event.AssertParams()

				l := asserts[id]
				asserts[id] = append(l, event)

				if len(l) > 0 {
					_, expected := l[0].AssertParams()
					if expected != value {
						fail := event
						fail.Type = "assert-failed"
						failed = true
						out <- fail
					}
				}
			}

			if event.Type == "exited" && event.ExitCode() != 0 && !done {
				failed = true
			}
		}

		if failed {
			out <- r.traceStatus(false)
		} else {
			out <- r.traceStatus(true)
		}
	}()

	return out
}

func (r *runner) run() <-chan Event {
	out := make(chan Event)

	go func() {
		defer func() { close(out) }()

		out <- r.traceExec()

		os.MkdirAll(r.wdir, 0755)
		defer os.RemoveAll(r.wdir)

		t := time.AfterFunc(r.timeout, r.stop)
		defer t.Stop()

		defer r.stop()

		for role := range r.procs {
			if role == "sut" || role == "driver" {
				continue
			}

			if ready := r.runRole(role, out); !<-ready {
				return
			}
		}

		if ready := r.runRole("sut", out); !<-ready {
			return
		}

		if ready := r.runRole("driver", out); !<-ready {
			return
		}

		for _, stopped := range r.stopped {
			<-stopped
		}
	}()

	return r.controller(out)
}

func (r *runner) stop() {
	for role := range r.kills {
		r.close(role)
	}

	for _, stopped := range r.stopped {
		<-stopped
	}
}

func (r *runner) runRole(role string, gout chan<- Event) (ready_ <-chan bool) {
	ready := make(chan bool, 1)

	proc := r.procs[role]
	if proc == nil {
		gout <- r.traceError(fmt.Errorf("missing role: %q", role))
		ready <- false
		return ready
	}

	out, kill, done := r.execer.Exec(r.wdir, role, proc)

	r.mtx.Lock()
	r.kills[role] = kill
	r.stopped[role] = done
	r.mtx.Unlock()

	go func() {
		defer r.close(role)

		isReady := false

		for event := range out {
			if event.Role == "" {
				event.Role = role
			}
			if event.Time == 0 {
				event.RecordTime()
			}

			gout <- event
			if event.Type == "ready" {
				isReady = true
				ready <- true
			}
		}

		if !isReady {
			ready <- false
		}
	}()

	return ready
}

func (r *runner) traceError(err error) Event {
	e := Event{Type: "log", Info: Info{"line": "error: " + err.Error()}}
	e.RecordTime()
	return e
}

func (r *runner) traceStatus(success bool) Event {
	e := Event{Type: "status", Info: Info{"success": success}}
	e.RecordTime()
	return e
}

func (r *runner) traceExec() Event {
	e := Event{Type: "exec", Info: Info{"sut": r.sut, "driver": r.driver}}
	e.RecordTime()
	return e
}
