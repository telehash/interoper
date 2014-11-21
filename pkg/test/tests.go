package test

import (
	"fmt"
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
		return fmt.Errorf("%s: `sut` must be present.", t.Name)
	}

	if t.TestDriver == nil {
		return fmt.Errorf("%s: `driver` must be present.", t.Name)
	}

	return nil
}

func (t *Description) Timeout() time.Duration {
	return t.timeout
}
