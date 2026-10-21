package setup

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategorySetup = "Setup"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category:    CategorySetup,
		Name:        "python.cwd",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("PYTHON_CWD")),
		Usage:       "Working directory of the uv project.",
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},

	&cli.StringFlag{
		Category:    CategorySetup,
		Name:        "python.cache",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("PYTHON_CACHE")),
		Usage:       "Cache directory for uv. Leave empty to use the environment defaults.",
		Required:    false,
		Value:       "./.uv/",
		Destination: &P.Cache,
	},
}
