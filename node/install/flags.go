package install

import (
	"github.com/urfave/cli/v2"
)

//revive:disable:line-length-limit

const (
	CATEGORY_NODE_INSTALL = "Install"
)

var Flags = []cli.Flag{

	// CATEGORY_NODE_INSTALL

	&cli.StringFlag{
		Category:    CATEGORY_NODE_INSTALL,
		Name:        "node.install_cwd",
		Usage:       "Install CWD for the package manager.",
		Required:    false,
		EnvVars:     []string{"NODE_INSTALL_CWD"},
		Value:       ".",
		Destination: &TL.Pipe.NodeInstall.Cwd,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_NODE_INSTALL,
		Name:        "node.use_lock_file",
		Usage:       "Use the lockfile while installing the packages.",
		Required:    false,
		EnvVars:     []string{"NODE_INSTALL_USE_LOCK_FILE"},
		Value:       true,
		Destination: &TL.Pipe.NodeInstall.UseLockFile,
	},

	&cli.StringFlag{
		Category:    CATEGORY_NODE_INSTALL,
		Name:        "node.install_args",
		Usage:       "Arguments to append to install command.",
		Required:    false,
		EnvVars:     []string{"NODE_INSTALL_ARGS"},
		Value:       "",
		Destination: &TL.Pipe.NodeInstall.Args,
	},
}
