package install

import (
	"github.com/urfave/cli/v3"
)

const CATEGORY_INSTALL = "Install"

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.BoolFlag{
		Category: CATEGORY_INSTALL,
		Name:     "terraform.install.reconfigure",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_INSTALL_RECONFIGURE"),
			cli.EnvVar("TF_INSTALL_RECONFIGURE"),
		),
		Usage:       "Reconfigure flag for terraform init.",
		Required:    false,
		Value:       false,
		Destination: &P.Install.Reconfigure,
	},

	&cli.BoolFlag{
		Category: CATEGORY_INSTALL,
		Name:     "terraform.install.use-lockfile",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_INSTALL_USE_LOCKFILE"),
			cli.EnvVar("TF_INSTALL_USE_LOCKFILE"),
		),
		Usage:       "Use lockfile for terraform init.",
		Required:    false,
		Value:       false,
		Destination: &P.Install.UseLockfile,
	},

	&cli.StringFlag{
		Category: CATEGORY_INSTALL,
		Name:     "terraform.install.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_INSTALL_ARGS"),
			cli.EnvVar("TF_INSTALL_ARGS"),
		),
		Usage:       "Additional arguments for terraform init.",
		Required:    false,
		Value:       "",
		Destination: &P.Install.Args,
	},
}
