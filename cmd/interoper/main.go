package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt-go"

	"github.com/telehash/interoper/pkg/runner"
	"github.com/telehash/interoper/pkg/suite"
	"github.com/telehash/interoper/pkg/test"
)

var usage = `interoper -- Interoperation testing for telehash

Usage:
  interoper test
  interoper suite pull
  interoper suite push [--suite=<dir>]
  interoper -h | --help
  interoper --version

Options:
  --suite=<dir>  Path to the test suite. [default: tests]
  -h --help      Show this screen.
  --version      Show version.

Pushing The Test Suite:
  Run interoper suite push from the root of the interoper directory.
  Make sure you have the INTEROPER_AWS_BUCKET, INTEROPER_AWS_ACCOUNT and
  INTEROPER_AWS_SECRET environment variables set.
`

func main() {
	args, _ := docopt.Parse(usage, nil, true, "0.1", false)

	var (
		isTest, _  = args["test"].(bool)
		isSuite, _ = args["suite"].(bool)
		isPull, _  = args["pull"].(bool)
		isPush, _  = args["push"].(bool)
	)

	switch {

	case isTest:
		mainTest(args)

	case isSuite && isPull:

	case isSuite && isPush:
		mainSuitePush(args)

	}
}

func mainTest(args map[string]interface{}) {
	var (
		ctx runner.Context
	)

	ctx.Drivers = []string{"go", "c"}
	ctx.SUT = "go"
	ctx.Execer = &runner.Docker{}

	t, err := test.LoadAll("tests")
	assert(err)

	for _, t := range t {
		runner.Run(t, &ctx)
	}

	err = ctx.Dump()
	assert(err)
}

func mainSuitePush(args map[string]interface{}) {
	var (
		INTEROPER_AWS_BUCKET  = os.Getenv("INTEROPER_AWS_BUCKET")
		INTEROPER_AWS_ACCOUNT = os.Getenv("INTEROPER_AWS_ACCOUNT")
		INTEROPER_AWS_SECRET  = os.Getenv("INTEROPER_AWS_SECRET")
		dir, _                = args["--suite"].(string)
	)

	cmd := suite.Push{
		Root:       dir,
		AwsBucket:  INTEROPER_AWS_BUCKET,
		AwsAccount: INTEROPER_AWS_ACCOUNT,
		AwsSecret:  INTEROPER_AWS_SECRET,
	}

	assert(cmd.Do())
}

func assert(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
