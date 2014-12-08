package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Command interface {
	Do() error
}

type implementations map[string]implementation

type implementation struct {
	Repo string `json:"repo"`
}

func (i *implementations) read(f string) error {
	if f == "" {
		f = "interop/tests/_implementations.json"
	}

	data, err := ioutil.ReadFile(f)
	if err != nil {
		return fmt.Errorf("unable to determine implementations sources: %s", err)
	}

	err = json.Unmarshal(data, i)
	if err != nil {
		return fmt.Errorf("unable to determine implementations sources: %s", err)
	}

	return nil
}
