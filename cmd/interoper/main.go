package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"

	"github.com/telehash/interoper/pkg/runner"
	"github.com/telehash/interoper/pkg/test"
)

var usage = `interoper -- Interoperation testing for telehash

Usage:
  interoper run <test>
  interoper -h | --help
  interoper --version

Options:
  -h --help     Show this screen.
  --version     Show version.
`

func main() {
	args, _ := docopt.Parse(usage, nil, true, "0.1", false)

	var (
		testFile = args["<test>"].(string)
		ctx      runner.Context
	)

	ctx.Drivers = []string{"go"}
	ctx.SUTs = []string{"go"}
	ctx.Execer = &runner.Docker{}

	t, err := test.Load(testFile)
	assert(err)

	err = runner.Run(t, &ctx)
	fmt.Printf("%v\n", ctx.Output)
	assert(err)
}

func assert(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
