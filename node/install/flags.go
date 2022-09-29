package install

import (
	"github.com/urfave/cli/v2"
)

//revive:disable:line-length-limit

const (
	category_node_install = "Install"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category:    category_node_install,
		Name:        "node.install_cwd",
		Usage:       "Install CWD for nodejs.",
		Required:    false,
		EnvVars:     []string{"NODE_INSTALL_CWD"},
		Value:       ".",
		Destination: &TL.Pipe.NodeInstall.Cwd,
	},

	&cli.BoolFlag{
		Category:    category_node_install,
		Name:        "node.use_lock_file",
		Usage:       "Whether to use lock file or not.",
		Required:    false,
		EnvVars:     []string{"NODE_INSTALL_USE_LOCK_FILE"},
		Value:       true,
		Destination: &TL.Pipe.NodeInstall.UseLockFile,
	},

	&cli.StringFlag{
		Category:    category_node_install,
		Name:        "node.install_args",
		Usage:       "Arguments for appending to installation.",
		Required:    false,
		EnvVars:     []string{"NODE_INSTALL_ARGS"},
		Value:       "",
		Destination: &TL.Pipe.NodeInstall.Args,
	},
}
