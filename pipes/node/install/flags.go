package install

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_NODE_INSTALL = "Install"
)

var Flags = []cli.Flag{

	// CATEGORY_NODE_INSTALL

	&cli.StringFlag{
		Category: CATEGORY_NODE_INSTALL,
		Name:     "node.install_cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_INSTALL_CWD"),
		),
		Usage:       "Install CWD for the package manager.",
		Required:    false,
		Value:       ".",
		Destination: &P.NodeInstall.Cwd,
	},

	&cli.BoolFlag{
		Category: CATEGORY_NODE_INSTALL,
		Name:     "node.use_lock_file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_INSTALL_USE_LOCK_FILE"),
		),
		Usage:       "Use the lockfile while installing the packages.",
		Required:    false,
		Value:       true,
		Destination: &P.NodeInstall.UseLockFile,
	},

	&cli.StringFlag{
		Category: CATEGORY_NODE_INSTALL,
		Name:     "node.install_args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_INSTALL_ARGS"),
		),
		Usage:       "Arguments to append to install command.",
		Required:    false,
		Value:       "",
		Destination: &P.NodeInstall.Args,
	},

	&cli.BoolFlag{
		Category: CATEGORY_NODE_INSTALL,
		Name:     "node.install_cache_enable",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_INSTALL_CACHE_ENABLE"),
		),
		Usage:       "Enable caching for the package manager.",
		Required:    false,
		Value:       true,
		Destination: &P.NodeInstall.Cache,
	},
}
