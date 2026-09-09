package test

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategoryTest = "Test"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CategoryTest,
		Name:     "go.test.runner",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_TEST_RUNNER"),
		),
		Usage:       `Test runner for the test command. format(enum("ginkgo", "go"))`,
		Required:    false,
		Value:       string(RunnerGinkgo),
		Destination: &P.Runner,
	},

	&cli.StringFlag{
		Category: CategoryTest,
		Name:     "go.test.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_TEST_ARGS"),
		),
		Usage:       "Arguments to append to the test command.",
		Required:    false,
		Value:       "",
		Destination: &P.Args,
	},

	&cli.StringFlag{
		Category: CategoryTest,
		Name:     "go.test.coverage",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_TEST_COVERAGE"),
		),
		Usage:       "Coverage profile to write. Leave empty to skip coverage.",
		Required:    false,
		Value:       "coverage.out",
		Destination: &P.Coverage,
	},
}
