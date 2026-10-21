package test

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

//revive:disable:line-length-limit

const (
	CategoryTest = "Test"
)

var Flags = []cli.Flag{
	flags.YAMLFlag(&P.Commands, &cli.StringFlag{
		Category: CategoryTest,
		Name:     "python.test.commands",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PYTHON_TEST_COMMANDS"),
		),
		Usage:    "Commands to test with, each run in order inside the installed environment. format(yaml([]string))",
		Required: false,
		Value:    `["pytest --junitxml=junit.xml"]`,
	}),
}
