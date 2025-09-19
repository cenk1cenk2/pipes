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
		Category: CATEGORY_SETUP,
		Name:     "go.build.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_CWD"),
		),
		Usage:       "Build CWD for the package manager.",
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},

	&cli.StringFlag{
		Category: CATEGORY_SETUP,
		Name:     "go.cache",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("GO_CACHE"),
		),
		Usage:       "Enable go cache.",
		Required:    false,
		Value:       "",
		Destination: &P.Cache,
	},
}
