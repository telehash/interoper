package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/telehash/interoper/pkg/web"
)

type Inspect struct {
	Report string
}

func (i *Inspect) Do() error {
	data, err := ioutil.ReadFile(i.Report)
	if err != nil {
		return fmt.Errorf("failed to read report: %s", err)
	}

	web.Run(data)
	return nil
}
