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
		Name:        "mise.cwd",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("MISE_CWD")),
		Usage:       "Working directory for mise commands.",
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},

	&cli.StringFlag{
		Category:    CategorySetup,
		Name:        "mise.data-dir",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("MISE_DATA_DIR")),
		Usage:       "Data directory for mise, which holds the binary, the installs and the shims. Every job that restores it has to use the same absolute path, since the shims link into it.",
		Required:    false,
		Value:       "./.mise/",
		Destination: &P.DataDir,
	},
}
