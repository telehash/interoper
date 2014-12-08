package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"

	"github.com/telehash/interoper/pkg/commands"
)

var usage = `interoper -- Inter-operation testing for telehash

Usage:
  interoper test
  interoper tournament
  interoper tests update
  interoper tests build [IMPLEMENTATION]
  interoper -h | --help
  interoper --version

Options:
  -h --help      Show this screen.
  --version      Show version.

Commands:

  $ interoper tournament
  Runs all implementations against each other.

  $ interoper test
  Builds a test image from the Dockerfile in the current directory.
  Then runs this image against all implementations.

  $ interoper tests update
  This command will download the latests test descriptions from github.com.

  $ interoper tests build
  This command builds all the test docker images defined by the repositories
  in the "interop/tests/_implementations.json" file.

  $ interoper tests build IMPLEMENTATION
  This command builds just the specified docker image.
`

func main() {
	args, _ := docopt.Parse(usage, nil, true, "0.1", false)

	var (
		cmd               commands.Command
		isTest, _         = args["test"].(bool)
		isTests, _        = args["tests"].(bool)
		isUpdate, _       = args["update"].(bool)
		isBuild, _        = args["build"].(bool)
		implementation, _ = args["IMPLEMENTATION"].(string)
	)

	switch {

	case isTest:
		cmd = &commands.Test{}

	case isTests && isUpdate:
		cmd = &commands.Update{}

	case isTests && isBuild:
		cmd = &commands.BuildImages{Which: implementation}

	default:
		assert(fmt.Errorf("unsupported command"))

	}

	assert(cmd.Do())
}

func assert(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
