package test

import (
	"io/ioutil"
	"path"
	"strings"
)

func ListAll(dir string) ([]string, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var l []string
	for _, entry := range entries {
		if !entry.Mode().IsRegular() {
			continue
		}

		if strings.HasPrefix(entry.Name(), ".") || strings.HasPrefix(entry.Name(), "_") {
			continue
		}

		if path.Ext(entry.Name()) != ".md" {
			continue
		}

		l = append(l, entry.Name())
	}

	return l, nil
}

func LoadAll(dir string) ([]*Description, error) {
	var (
		l []*Description
	)

	entries, err := ListAll(dir)
	if err != nil {
		return nil, err
	}

	for _, name := range entries {
		d, err := Load(path.Join(dir, name))
		if err != nil {
			return nil, err
		}

		l = append(l, d)
	}

	return l, nil
}

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
