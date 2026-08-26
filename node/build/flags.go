package build

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
)

//revive:disable:line-length-limit

const (
	CATEGORY_NODE_BUILD = "Build"
)

var Flags = []cli.Flag{
	// CATEGORY_BUILD

	&cli.StringFlag{
		Category: CATEGORY_NODE_BUILD,
		Name:     "node.build.script",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_BUILD_SCRIPT"),
		),
		Usage:       fmt.Sprintf("package.json script for building operation. %s", environment.HELP_FORMAT_TEMPLATE),
		Required:    false,
		Value:       "build",
		Destination: &P.Build.Script,
	},

	&cli.StringFlag{
		Category: CATEGORY_NODE_BUILD,
		Name:     "node.build.script-args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_BUILD_SCRIPT_ARGS"),
		),
		Usage:       fmt.Sprintf("package.json script arguments for building operation. %s", environment.HELP_FORMAT_TEMPLATE),
		Required:    false,
		Value:       "",
		Destination: &P.Build.ScriptArgs,
	},

	&cli.StringFlag{
		Category: CATEGORY_NODE_BUILD,
		Name:     "node.build.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_BUILD_CWD"),
		),
		Usage:       "Working directory for build operation.",
		Required:    false,
		Value:       ".",
		Destination: &P.Build.Cwd,
	},
}
