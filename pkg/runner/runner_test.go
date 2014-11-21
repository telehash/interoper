package runner

import (
	"fmt"
	"time"

	"github.com/telehash/interoper/pkg/test"
)

func Example_permute() {
	permute([]string{"go", "c", "js"}, []string{"c"}, func(sut, driver string) error {
		fmt.Printf("test: sut=%q driver=%q\n", sut, driver)
		return nil
	})

	// Output:
	// test: sut="go" driver="c"
	// test: sut="c" driver="c"
	// test: sut="js" driver="c"
}

func Example_Run() {
	var (
		ctx = &Context{execer: &execer{}, suts: []string{"go", "c", "js"}, drivers: []string{"go", "c"}}
		t   = test.MustLoad("testdata/example.md")
	)

	err := Run(t, ctx)
	if err != nil {
		panic(err)
	}

	// Output:
	// run: name="sut" impl="go" impl="test-net-link await"
	// run: name="driver" impl="go" impl="test-net-link establish"
	// run: name="sut" impl="go" impl="test-net-link await"
	// run: name="driver" impl="c" impl="test-net-link establish"
	// run: name="sut" impl="c" impl="test-net-link await"
	// run: name="driver" impl="go" impl="test-net-link establish"
	// run: name="sut" impl="c" impl="test-net-link await"
	// run: name="driver" impl="c" impl="test-net-link establish"
	// run: name="sut" impl="js" impl="test-net-link await"
	// run: name="driver" impl="go" impl="test-net-link establish"
	// run: name="sut" impl="js" impl="test-net-link await"
	// run: name="driver" impl="c" impl="test-net-link establish"
}

type execer struct {
}

type process struct {
	c <-chan time.Time
}

func (e *execer) Exec(name, impl, command string) (Process, error) {
	fmt.Printf("run: name=%q impl=%q impl=%q\n", name, impl, command)
	return &process{time.After(10 * time.Millisecond)}, nil
}

func (p *process) Kill() error {
	return nil
}

func (p *process) Wait() error {
	<-p.c
	return nil
}
