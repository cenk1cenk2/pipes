package lint

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

//revive:disable:line-length-limit

const (
	CategoryLint = "Lint"
)

var Flags = []cli.Flag{
	flags.YAMLFlag(&P.Commands, &cli.StringFlag{
		Category: CategoryLint,
		Name:     "python.lint.commands",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PYTHON_LINT_COMMANDS"),
		),
		Usage:    "Commands to lint with, each run in order inside the installed environment. format(yaml([]string))",
		Required: false,
		Value:    `["ruff format --check .", "ruff check .", "ty check ."]`,
	}),
}
