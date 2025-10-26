package lint

import (
	"time"

	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_LINT = "Lint"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_LINT,
		Name:     "go.lint.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_LINT_ARGS"),
		),
		Usage:       "Arguments to append to lint command.",
		Required:    false,
		Value:       "",
		Destination: &P.Args,
	},

	&cli.DurationFlag{
		Category: CATEGORY_LINT,
		Name:     "go.lint.timeout",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_LINT_TIMEOUT"),
		),
		Usage:       "Timeout for the lint command.",
		Required:    false,
		Value:       time.Duration(5 * time.Minute),
		Destination: &P.Timeout,
	},
}
