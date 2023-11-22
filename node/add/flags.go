package pipe

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/common/flags"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringSliceFlag{
		Category: flags.CATEGORY_PACKAGES,
		Name:     "packages.node",
		Usage:    "Install node packages before performing operations.",
		Required: true,
		EnvVars:  []string{"PACKAGES_NODE"},
		Value:    &cli.StringSlice{},
	},

	&cli.BoolFlag{
		Category:    flags.CATEGORY_PACKAGES,
		Name:        "packages.node.global",
		Usage:       "Install node packages globally.",
		Required:    false,
		EnvVars:     []string{"PACKAGES_NODE_GLOBAL"},
		Value:       true,
		Destination: &TL.Pipe.NodeAdd.Global,
	},

	&cli.StringFlag{
		Category:    flags.CATEGORY_PACKAGES,
		Name:        "packages.node.script_args",
		Usage:       fmt.Sprintf("package.json script arguments for building operation. format(%s)", environment.HELP_FORMAT_ENVIRONMENT_TEMPLATE),
		Required:    false,
		EnvVars:     []string{"PACKAGES_NODE_SCRIPT_ARGS"},
		Value:       "",
		Destination: &TL.Pipe.NodeAdd.ScriptArgs,
	},

	&cli.StringFlag{
		Category:    flags.CATEGORY_PACKAGES,
		Name:        "packages.node.cwd",
		Usage:       "Working directory for build operation.",
		Required:    false,
		EnvVars:     []string{"PACKAGES_NODE_CWD"},
		Value:       ".",
		Destination: &TL.Pipe.NodeAdd.Cwd,
	},
}

func ProcessFlags(tl *TaskList[Pipe]) error {
	tl.Pipe.NodeAdd.Packages = tl.CliContext.StringSlice("packages.node")

	return nil
}
