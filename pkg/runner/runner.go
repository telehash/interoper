package runner

import (
	"fmt"
	"time"

	"github.com/telehash/interoper/pkg/test"
)

type Execer interface {
	Exec(name, implementation, command string) (Process, error)
}

type Process interface {
	Kill() error
	Wait() error
	Output() []*OutputLine
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
		)

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

		{ // start the SUT
			proc, err := ctx.Execer.Exec("sut", sut, t.SystemUnderTest.Command)
			if err != nil {
				// TODO report
				return nil
			}
			procs = append(procs, proc)
		}

		for name, desc := range t.Services { // start the SUT
			proc, err := ctx.Execer.Exec("srv-"+name, driver, desc.Command)
			if err != nil {
				// TODO report
				return nil
			}
			procs = append(procs, proc)
		}

		{ // start the Driver
			proc, err := ctx.Execer.Exec("driver", driver, t.TestDriver.Command)
			if err != nil {
				// TODO report
				return nil
			}
			procs = append(procs, proc)
		}

		for _, proc := range procs {
			err := proc.Wait()
			if err != nil {
				// TODO report
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
