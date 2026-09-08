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
		Usage:       "Cache directory for go commands. Leave empty to use the environment defaults.",
		Required:    false,
		Value:       "./.go/",
		Destination: &P.Cache,
	},
}
