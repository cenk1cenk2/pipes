package add

import (
	"github.com/urfave/cli/v3"
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
			cli.EnvVar("NODE_ADD_PACKAGES"),
			cli.EnvVar("PACKAGES_NODE"),
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
			cli.EnvVar("NODE_ADD_GLOBAL"),
			cli.EnvVar("PACKAGES_NODE_GLOBAL"),
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
			cli.EnvVar("NODE_ADD_SCRIPT_ARGS"),
			cli.EnvVar("PACKAGES_NODE_SCRIPT_ARGS"),
		),
		Usage:       "Script arguments to append to the install command.",
		Required:    false,
		Value:       "",
		Destination: &P.Add.ScriptArgs,
	},

	&cli.StringFlag{
		Category: CATEGORY_PACKAGES,
		Name:     "node.add.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_ADD_CWD"),
			cli.EnvVar("PACKAGES_NODE_CWD"),
		),
		Usage:       "Working directory for the add operation.",
		Required:    false,
		Value:       ".",
		Destination: &P.Add.Cwd,
	},
}
