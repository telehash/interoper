package test

import (
	"fmt"
	"path"
	"strings"
	"time"
)

type Description struct {
	Name          string
	TimeoutString string `json:"timeout"`

	SystemUnderTest *Process            `json:"sut"`
	TestDriver      *Process            `json:"driver"`
	Services        map[string]*Process `json:"services"`

	Description string

	timeout time.Duration
}

type Process struct {
	Command string
}

func (t *Description) Normalize() error {
	if t.TimeoutString == "" {
		t.TimeoutString = "0s"
	}

	if d, err := time.ParseDuration(t.TimeoutString); err != nil {
		return fmt.Errorf("%s: `timeout` is invalid: %s", t.Name, err)
	} else {
		t.timeout = d
	}

	if t.timeout <= 0 {
		t.timeout = 10 * time.Minute
	}

	if t.SystemUnderTest == nil {
		t.SystemUnderTest = &Process{}
	}

	if t.SystemUnderTest.Command == "" {
		t.SystemUnderTest.Command = fmt.Sprintf("th-test %s sut", strings.TrimSuffix(path.Base(t.Name), ".md"))
	}

	if t.TestDriver == nil {
		t.TestDriver = &Process{}
	}

	if t.TestDriver.Command == "" {
		t.TestDriver.Command = fmt.Sprintf("th-test %s driver", strings.TrimSuffix(path.Base(t.Name), ".md"))
	}

	return nil
}

func (t *Description) Timeout() time.Duration {
	return t.timeout
}
