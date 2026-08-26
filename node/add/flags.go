package pipe

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
)

//revive:disable:line-length-limit

const (
	CATEGORY_PACKAGES = "Packages"
)

var Flags = []cli.Flag{
	&cli.StringSliceFlag{
		Category: CATEGORY_PACKAGES,
		Name:     "node.add.packages",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PACKAGES_NODE"),
			cli.EnvVar("NODE_ADD_PACKAGES"),
		),
		Usage:       "Install node packages before performing operations.",
		Required:    true,
		Value:       []string{},
		Destination: &P.Add.Packages,
	},

	&cli.BoolFlag{
		Category: CATEGORY_PACKAGES,
		Name:     "node.add.global",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PACKAGES_NODE_GLOBAL"),
			cli.EnvVar("NODE_ADD_GLOBAL"),
		),
		Usage:       "Install node packages globally.",
		Required:    false,
		Value:       true,
		Destination: &P.Add.Global,
	},

	&cli.StringFlag{
		Category: CATEGORY_PACKAGES,
		Name:     "node.add.script-args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PACKAGES_NODE_SCRIPT_ARGS"),
			cli.EnvVar("NODE_ADD_SCRIPT_ARGS"),
		),
		Usage:       fmt.Sprintf("package.json script arguments for building operation. %s", environment.HELP_FORMAT_TEMPLATE),
		Required:    false,
		Value:       "",
		Destination: &P.Add.ScriptArgs,
	},

	&cli.StringFlag{
		Category: CATEGORY_PACKAGES,
		Name:     "node.add.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("PACKAGES_NODE_CWD"),
			cli.EnvVar("NODE_ADD_CWD"),
		),
		Usage:       "Working directory for build operation.",
		Required:    false,
		Value:       ".",
		Destination: &P.Add.Cwd,
	},
}
