package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type BuildImages struct {
	Which           string
	implementations implementations
}

func (cmd *BuildImages) Do() error {
	if err := cmd.implementations.read(""); err != nil {
		return err
	}

	if cmd.Which == "" {
		if err := cmd.buildAll(); err != nil {
			return err
		}
	} else {
		if err := cmd.buildOne(); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *BuildImages) buildAll() error {
	for name, imp := range cmd.implementations {
		err := cmd.buildImage(name, imp.Repo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cmd *BuildImages) buildOne() error {
	imp, found := cmd.implementations[cmd.Which]
	if !found {
		return fmt.Errorf("unknown implementation: %q", cmd.Which)
	}

	return cmd.buildImage(cmd.Which, imp.Repo)
}

func (cmd *BuildImages) buildImage(name, url string) error {
	args := []string{
		"docker",
		"build",
		"--force-rm=true",
		"--tag", "telehashinterop/" + name,
		url,
	}

	fmt.Printf("$ %s\n", strings.Join(args, " "))

	c := exec.Command(args[0], args[1:]...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}
