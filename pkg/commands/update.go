package commands

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const (
	repoURL = "https://github.com/telehash/interoper/archive/master.zip"
)

type Update struct {
	buf   bytes.Buffer
	files map[string][]byte
}

func (u *Update) Do() error {
	if err := u.download(); err != nil {
		return err
	}

	if err := u.extract(); err != nil {
		return err
	}

	if err := u.write(); err != nil {
		return err
	}

	return nil
}

func (u *Update) download() error {
	resp, err := http.Get(repoURL)
	if err != nil {
		return fmt.Errorf("failed to download tests: %s", err)
	}

	defer resp.Body.Close()

	_, err = io.Copy(&u.buf, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to download tests: %s", err)
	}

	return nil
}

func (u *Update) extract() error {
	u.files = make(map[string][]byte)

	r, err := zip.NewReader(bytes.NewReader(u.buf.Bytes()), int64(u.buf.Len()))
	if err != nil {
		return fmt.Errorf("failed to extract tests: %s", err)
	}

	for _, file := range r.File {
		if !file.Mode().IsRegular() {
			continue
		}

		matched, err := filepath.Match("*/tests/*.md", file.Name)
		if err != nil {
			return fmt.Errorf("failed to extract tests: %s", err)
		}

		if !matched {
			matched, err = filepath.Match("*/tests/_implementations.json", file.Name)
			if err != nil {
				return fmt.Errorf("failed to extract tests: %s", err)
			}
		}

		if !matched {
			continue
		}

		r, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to extract tests: %s", err)
		}

		data, err := ioutil.ReadAll(r)
		if err != nil {
			return fmt.Errorf("failed to extract tests: %s", err)
		}

		err = r.Close()
		if err != nil {
			return fmt.Errorf("failed to extract tests: %s", err)
		}

		u.files[filepath.Base(file.Name)] = data
	}

	return nil
}

func (u *Update) write() error {
	os.RemoveAll("interop/tests")

	err := os.MkdirAll("interop/tests", 0755)
	if err != nil {
		return fmt.Errorf("failed to write tests: %s", err)
	}

	for name, data := range u.files {
		err := ioutil.WriteFile(filepath.Join("interop/tests", name), data, 0644)
		if err != nil {
			return fmt.Errorf("failed to write tests: %s", err)
		}
	}

	return nil
}
