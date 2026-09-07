package setup

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_SETUP = "Setup"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category:    CATEGORY_SETUP,
		Name:        "go.cwd",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("GO_CWD")),
		Usage:       "Working directory for go commands.",
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},

	&cli.StringFlag{
		Category:    CATEGORY_SETUP,
		Name:        "go.cache",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("GO_CACHE")),
		Usage:       "Enable go cache.",
		Required:    false,
		Value:       "./.go/",
		Destination: &P.Cache,
	},

	&cli.BoolFlag{
		Category: CATEGORY_SETUP,
		Name:     "go.workspace",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_WORKSPACE"),
			cli.EnvVar("GO_LINT_WORKSPACE"),
		),
		Usage:       "Drive the modules as a Go workspace instead of the single module in the working directory.",
		Required:    false,
		Value:       false,
		Destination: &P.Workspace,
	},
}
