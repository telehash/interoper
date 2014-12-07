package suite

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/s3"

	"github.com/telehash/interoper/pkg/test"
)

type Push struct {
	Root       string
	AwsBucket  string
	AwsAccount string
	AwsSecret  string

	buf bytes.Buffer
}

func (p *Push) Do() error {
	if p.Root == "" {
		p.Root = "."
	}

	if err := p.pack(); err != nil {
		return fmt.Errorf("pack failed: %s", err)
	}

	if err := p.upload(); err != nil {
		return fmt.Errorf("push failed: %s", err)
	}

	return nil
}

func (p *Push) pack() error {
	p.buf.Reset()

	l, err := test.ListAll(p.Root)
	if err != nil {
		return err
	}

	z := zip.NewWriter(&p.buf)

	{ // write the images file
		name := "_images.json"

		w, err := z.Create(name)
		if err != nil {
			return err
		}

		data, err := ioutil.ReadFile(path.Join(p.Root, name))
		if err != nil {
			return err
		}

		_, err = w.Write(data)
		if err != nil {
			return err
		}
	}

	for _, name := range l {

		w, err := z.Create(name)
		if err != nil {
			return err
		}

		data, err := ioutil.ReadFile(path.Join(p.Root, name))
		if err != nil {
			return err
		}

		_, err = w.Write(data)
		if err != nil {
			return err
		}

	}

	return z.Close()
}

func (p *Push) upload() error {
	b := s3.New(aws.Auth{AccessKey: p.AwsAccount, SecretKey: p.AwsSecret}, aws.USEast).Bucket(p.AwsBucket)
	return b.Put("suite.zip", p.buf.Bytes(), "application/zip", s3.PublicRead, s3.Options{})
}
