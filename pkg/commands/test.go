package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/telehash/interoper/pkg/runner"
	"github.com/telehash/interoper/pkg/test"
)

type Test struct {
	ctx             runner.Context
	implementations implementations
	tests           []*test.Description
}

func (cmd *Test) Do() error {
	if err := cmd.implementations.read(""); err != nil {
		return err
	}

	if err := cmd.readTests(); err != nil {
		return err
	}

	if err := cmd.buildImage(); err != nil {
		return err
	}

	if err := cmd.runTests(); err != nil {
		return err
	}

	return nil
}

func (cmd *Test) readTests() error {
	t, err := test.LoadAll("interop/tests")
	if err != nil {
		return fmt.Errorf("failed to load tests: %s", err)
	}

	cmd.tests = t
	return nil
}

func (cmd *Test) buildImage() error {
	args := []string{
		"docker",
		"build",
		"--force-rm=true",
		"--tag", "telehashinterop/local",
		".",
	}

	fmt.Printf("$ %s", strings.Join(args, " "))

	c := exec.Command(args[0], args[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func (cmd *Test) runTests() error {
	cmd.ctx.SUT = "local"
	cmd.ctx.Execer = &runner.Docker{}

	cmd.ctx.Drivers = append(cmd.ctx.Drivers, "local")
	for name := range cmd.implementations {
		cmd.ctx.Drivers = append(cmd.ctx.Drivers, name)
	}

	for _, t := range cmd.tests {
		runner.Run(t, &cmd.ctx)
	}

	return cmd.ctx.Report.Dump()
}
