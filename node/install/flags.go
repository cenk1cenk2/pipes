package install

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategoryNodeInstall = "Install"
)

var Flags = []cli.Flag{

	// CategoryNodeInstall

	&cli.StringFlag{
		Category: CategoryNodeInstall,
		Name:     "node.install.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_INSTALL_CWD"),
		),
		Usage:       "Working directory for the install operation.",
		Required:    false,
		Value:       ".",
		Destination: &P.Install.Cwd,
	},

	&cli.BoolFlag{
		Category: CategoryNodeInstall,
		Name:     "node.install.use-lock-file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_INSTALL_USE_LOCK_FILE"),
		),
		Usage:       "Use the lockfile while installing the packages.",
		Required:    false,
		Value:       true,
		Destination: &P.Install.UseLockFile,
	},

	&cli.StringFlag{
		Category: CategoryNodeInstall,
		Name:     "node.install.args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_INSTALL_ARGS"),
		),
		Usage:       "Arguments to append to the install command.",
		Required:    false,
		Value:       "",
		Destination: &P.Install.Args,
	},

	&cli.BoolFlag{
		Category: CategoryNodeInstall,
		Name:     "node.install.cache",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_INSTALL_CACHE"),
		),
		Usage:       "Enable caching for the package manager.",
		Required:    false,
		Value:       true,
		Destination: &P.Install.Cache,
	},
}
