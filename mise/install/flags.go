package install

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategoryInstall = "Install"
)

var Flags = []cli.Flag{
	&cli.BoolFlag{
		Category: CategoryInstall,
		Name:     "mise.install.prune",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("MISE_INSTALL_PRUNE"),
		),
		Usage:       "Remove the tool versions the configuration no longer asks for, so the data directory does not grow with every bump.",
		Required:    false,
		Value:       true,
		Destination: &P.Prune,
	},

	&cli.StringFlag{
		Category: CategoryInstall,
		Name:     "mise.install.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("MISE_INSTALL_ARGS"),
		),
		Usage:       "Arguments to append to the install command.",
		Required:    false,
		Value:       "",
		Destination: &P.Args,
	},
}
