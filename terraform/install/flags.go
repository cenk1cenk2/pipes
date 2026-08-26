package install

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name: "terraform-install.reconfigure",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_INSTALL_RECONFIGURE"),
		),
		Usage:       "Reconfigure flag for terraform init.",
		Required:    false,
		Value:       false,
		Destination: &P.Install.Reconfigure,
	},

	&cli.BoolFlag{
		Name: "terraform-install.use-lockfile",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_INSTALL_USE_LOCKFILE"),
		),
		Usage:       "Use lockfile for terraform init.",
		Required:    false,
		Value:       false,
		Destination: &P.Install.UseLockfile,
	},

	&cli.StringFlag{
		Name: "terraform-install.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_INSTALL_ARGS"),
		),
		Usage:       "Additional arguments for terraform init.",
		Required:    false,
		Value:       "",
		Destination: &P.Install.Args,
	},
}
