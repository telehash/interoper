package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/telehash/interoper/pkg/test"
)

type Execer interface {
	Exec(wdir, name, implementation, command string) (Process, error)
}

type Process interface {
	Kill() error
	Wait() error
	Output() []*OutputLine
	OutputReader() io.Reader
}

type OutputLine struct {
	Process string
	Stream  string // stdout stderr
	Time    time.Time
	Line    string
}

type Context struct {
	Drivers []string
	SUTs    []string
	Execer  Execer
	Output  []*OutputLine
}

func Run(t *test.Description, ctx *Context) error {
	return permute(ctx.SUTs, ctx.Drivers, func(sut, driver string) error {
		fmt.Printf("SUT: %s Driver: %s\n", sut, driver)

		var (
			procs = make([]Process, 0, 2+len(t.Services))
			wdir  = path.Join(os.TempDir(), fmt.Sprintf("interop-%d", time.Now().UnixNano()))
		)

		os.MkdirAll(wdir, 0755)
		defer os.RemoveAll(wdir)

		var kill = func() {
			for _, proc := range procs {
				if proc != nil {
					proc.Kill()
				}
			}
		}

		var capture = func() {
			for _, proc := range procs {
				ctx.Output = append(ctx.Output, proc.Output()...)
			}
		}

		defer capture()
		defer kill()
		timeout := time.AfterFunc(t.Timeout(), kill)
		defer timeout.Stop()

		for name, desc := range t.Services { // start the Services
			proc, err := ctx.Execer.Exec(wdir, "srv-"+name, driver, desc.Command)
			if err != nil {
				// TODO report
				fmt.Printf("error: %s\n", err)
				return nil
			}
			err = awaitReady(proc.OutputReader())
			if err != nil {
				// TODO report
				fmt.Printf("error: %s\n", err)
				return nil
			}
			procs = append(procs, proc)
		}

		{ // start the SUT
			proc, err := ctx.Execer.Exec(wdir, "sut", sut, t.SystemUnderTest.Command)
			if err != nil {
				// TODO report
				fmt.Printf("error: %s\n", err)
				return nil
			}
			err = awaitReady(proc.OutputReader())
			if err != nil {
				// TODO report
				fmt.Printf("error: %s\n", err)
				return nil
			}
			procs = append(procs, proc)
		}

		{ // start the Driver
			proc, err := ctx.Execer.Exec(wdir, "driver", driver, t.TestDriver.Command)
			if err != nil {
				// TODO report
				fmt.Printf("error: %s\n", err)
				return nil
			}
			err = awaitReady(proc.OutputReader())
			if err != nil {
				// TODO report
				fmt.Printf("error: %s\n", err)
				return nil
			}
			procs = append(procs, proc)
		}

		for _, proc := range procs {
			err := proc.Wait()
			if err != nil {
				// TODO report
				fmt.Printf("error: %s\n", err)
				return nil
			}
		}

		return nil
	})
}

func (o *OutputLine) String() string {
	return fmt.Sprintf("ps=%s stream=%s at=%s %q", o.Process, o.Stream, o.Time, o.Line)
}

func permute(suts, drivers []string, f func(sut, driver string) error) error {
	for _, sut := range suts {
		for _, driver := range drivers {
			err := f(sut, driver)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type Command struct {
	Cmd string
}

func readCommand(r io.Reader) (*Command, error) {
	var cmd *Command
	err := json.NewDecoder(r).Decode(&cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func awaitReady(r io.Reader) error {
	cmd, err := readCommand(r)
	if err != nil {
		return err
	}
	if cmd.Cmd != "ready" {
		return fmt.Errorf("expected %q instead of %q", "ready", cmd.Cmd)
	}
	return nil
}
