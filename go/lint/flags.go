package lint

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_LINT = "Lint"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_LINT,
		Name:     "go.lint.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_LINT_CWD"),
		),
		Usage:       "Lint CWD for the package manager.",
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},

	&cli.StringFlag{
		Category: CATEGORY_LINT,
		Name:     "go.lint.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_LINT_ARGS"),
		),
		Usage:       "Arguments to append to lint command.",
		Required:    false,
		Value:       "run --timeout 3m",
		Destination: &P.Args,
	},

	&cli.StringFlag{
		Category: CATEGORY_LINT,
		Name:     "go.lint.source",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_LINT_SOURCE"),
		),
		Usage:       `Source for the linter. enum("tools")`,
		Required:    false,
		Value:       "tools",
		Destination: &P.Source,
	},

	&cli.StringFlag{
		Category: CATEGORY_LINT,
		Name:     "go.lint.tool",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_LINT_TOOL"),
		),
		Usage:       "Binary that provides the linting.",
		Required:    false,
		Value:       "golangci-lint",
		Destination: &P.Tool,
	},
}
