package test

import (
	"fmt"
	"path"
	"strings"
	"time"
)

type Description struct {
	Name          string              `json:"name"`
	TimeoutString string              `json:"timeout"`
	Containers    map[string]*Process `json:"containers"`
	Description   string              `json:"description"`

	timeout time.Duration
}

type Process struct {
	Image   string `json:"image"`
	Command string `json:"command"`
}

func (p *Process) ForImplementation(impl string) *Process {
	o := new(Process)
	*o = *p
	if o.Image == "" {
		o.Image = "telehashinterop/" + impl
	}
	return o
}

func (t *Description) Normalize() error {
	t.Name = strings.TrimSuffix(path.Base(t.Name), ".md")

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

	if t.Containers == nil {
		t.Containers = make(map[string]*Process)
	}

	if t.Containers["worker"] == nil {
		t.Containers["worker"] = &Process{}
	}

	if t.Containers["worker"].Command == "" {
		t.Containers["worker"].Command = fmt.Sprintf("th-test %s worker", strings.TrimSuffix(path.Base(t.Name), ".md"))
	}

	if t.Containers["driver"] == nil {
		t.Containers["driver"] = &Process{}
	}

	if t.Containers["driver"].Command == "" {
		t.Containers["driver"].Command = fmt.Sprintf("th-test %s driver", strings.TrimSuffix(path.Base(t.Name), ".md"))
	}

	return nil
}

func (t *Description) Timeout() time.Duration {
	return t.timeout
}
