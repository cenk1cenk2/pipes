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
		Name:     "python.install.locked",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PYTHON_INSTALL_LOCKED"),
		),
		Usage:       "Fail when the lockfile is out of date with the project. Otherwise the lockfile is installed as it is.",
		Required:    false,
		Value:       true,
		Destination: &P.Locked,
	},

	&cli.StringFlag{
		Category: CategoryInstall,
		Name:     "python.install.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PYTHON_INSTALL_ARGS"),
		),
		Usage:       "Arguments to append to the install command.",
		Required:    false,
		Value:       "",
		Destination: &P.Args,
	},
}
