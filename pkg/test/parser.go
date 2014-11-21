package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

func Load(file string) (*Description, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return Parse(file, data)
}

func MustLoad(file string) *Description {
	t, err := Load(file)
	if err != nil {
		panic(err)
	}
	return t
}

func Parse(name string, data []byte) (*Description, error) {
	var (
		t   *Description
		dec = json.NewDecoder(bytes.NewReader(data))
		err = dec.Decode(&t)
	)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse: %s", name, err)
	}

	data, err = ioutil.ReadAll(dec.Buffered())
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse: %s", name, err)
	}

	t.Name = name
	t.Description = strings.TrimSpace(string(data))

	err = t.Normalize()
	if err != nil {
		return nil, err
	}

	return t, nil
}
