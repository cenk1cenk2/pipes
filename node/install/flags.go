package install

import (
	"github.com/urfave/cli/v2"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "node.install_cwd",
		Usage:       "Install CWD for nodejs.",
		Required:    false,
		EnvVars:     []string{"NODE_INSTALL_CWD"},
		Value:       ".",
		Destination: &TL.Pipe.NodeInstall.Cwd,
	},

	&cli.BoolFlag{
		Name:        "node.use_lock_file",
		Usage:       "Whether to use lock file or not.",
		Required:    false,
		EnvVars:     []string{"NODE_INSTALL_USE_LOCK_FILE"},
		Value:       true,
		Destination: &TL.Pipe.NodeInstall.UseLockFile,
	},

	&cli.StringFlag{
		Name:        "node.install_args",
		Usage:       "Arguments for appending to installation.",
		Required:    false,
		EnvVars:     []string{"NODE_INSTALL_ARGS"},
		Value:       "",
		Destination: &TL.Pipe.NodeInstall.Args,
	},
}
