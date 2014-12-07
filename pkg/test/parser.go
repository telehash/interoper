package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

func Parse(name string, data []byte) (*Description, error) {
	var (
		t   *Description
		r   = bytes.NewReader(data)
		dec = json.NewDecoder(r)
		err = dec.Decode(&t)
	)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse: %s", name, err)
	}

	data, err = ioutil.ReadAll(io.MultiReader(dec.Buffered(), r))
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
