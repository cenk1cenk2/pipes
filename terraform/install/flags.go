package install

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.BoolFlag{
		Name:        "terraform-install.reconfigure",
		Usage:       "Reconfigure flag for terraform init.",
		Required:    false,
		EnvVars:     []string{"TF_INSTALL_RECONFIGURE"},
		Value:       false,
		Destination: &TL.Pipe.Install.Reconfigure,
	},

	&cli.BoolFlag{
		Name:        "terraform-install.use-lockfile",
		Usage:       "Use lockfile for terraform init.",
		Required:    false,
		EnvVars:     []string{"TF_INSTALL_USE_LOCKFILE"},
		Value:       false,
		Destination: &TL.Pipe.Install.UseLockfile,
	},

	&cli.StringFlag{
		Name:        "terraform-install.args",
		Usage:       "Additional arguments for terraform init.",
		Required:    false,
		EnvVars:     []string{"TF_INSTALL_ARGS"},
		Value:       "",
		Destination: &TL.Pipe.Install.Args,
	},
}

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList[Pipe]) error {
	return nil
}
