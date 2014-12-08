package commands

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/telehash/interoper/pkg/web"
)

type Inspect struct {
	Report string
}

func (i *Inspect) Do() error {
	if i.Report == "latest" {
		r, err := i.getLatest()
		if err != nil {
			return err
		}

		i.Report = r
	}

	data, err := ioutil.ReadFile(i.Report)
	if err != nil {
		return fmt.Errorf("failed to read report: %s", err)
	}

	web.Run(data)
	return nil
}

func (i *Inspect) getLatest() (string, error) {
	entries, err := ioutil.ReadDir("interop/reports")
	if err != nil {
		return "", fmt.Errorf("failed to read report: %s", err)
	}

	last := ""
	for _, entry := range entries {
		if path.Ext(entry.Name()) == ".dump" && entry.Mode().IsRegular() {
			name := path.Join("interop/reports", entry.Name())

			if name > last {
				last = name
			}
		}
	}

	if last == "" {
		return "", fmt.Errorf("failed to read report: %s", "no reports we found")
	}

	return last, nil
}
