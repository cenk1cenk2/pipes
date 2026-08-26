package pipe

import (
	"fmt"

	"github.com/urfave/cli/v3"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
)

//revive:disable:line-length-limit

const (
	CATEGORY_PACKAGES = "Packages"
)

var Flags = []cli.Flag{
	&cli.StringSliceFlag{
		Category: CATEGORY_PACKAGES,
		Name:     "packages.node",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PACKAGES_NODE"),
		),
		Usage:       "Install node packages before performing operations.",
		Required:    true,
		Value:       []string{},
		Destination: &P.NodeAdd.Packages,
	},

	&cli.BoolFlag{
		Category: CATEGORY_PACKAGES,
		Name:     "packages.node.global",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PACKAGES_NODE_GLOBAL"),
		),
		Usage:       "Install node packages globally.",
		Required:    false,
		Value:       true,
		Destination: &P.NodeAdd.Global,
	},

	&cli.StringFlag{
		Category: CATEGORY_PACKAGES,
		Name:     "packages.node.script_args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PACKAGES_NODE_SCRIPT_ARGS"),
		),
		Usage:       fmt.Sprintf("package.json script arguments for building operation. %s", environment.HELP_FORMAT_ENVIRONMENT_TEMPLATE),
		Required:    false,
		Value:       "",
		Destination: &P.NodeAdd.ScriptArgs,
	},

	&cli.StringFlag{
		Category: CATEGORY_PACKAGES,
		Name:     "packages.node.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PACKAGES_NODE_CWD"),
		),
		Usage:       "Working directory for build operation.",
		Required:    false,
		Value:       ".",
		Destination: &P.NodeAdd.Cwd,
	},
}
