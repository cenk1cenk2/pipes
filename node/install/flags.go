package install

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "node.install_cwd",
		Usage:       "Install CWD for nodejs.",
		Required:    false,
		EnvVars:     []string{"NODE_INSTALL_CWD"},
		Value:       ".",
		Destination: &P.Pipe.NodeInstall.Cwd,
	},

	&cli.BoolFlag{
		Name:        "node.use_lock_file",
		Usage:       "Whether to use lock file or not.",
		Required:    false,
		EnvVars:     []string{"NODE_INSTALL_USE_LOCK_FILE"},
		Value:       true,
		Destination: &P.Pipe.NodeInstall.UseLockFile,
	},
}
