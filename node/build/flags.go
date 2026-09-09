package build

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
)

//revive:disable:line-length-limit

const (
	CategoryNodeBuild = "Build"
)

var Flags = []cli.Flag{
	// CategoryBuild

	&cli.StringFlag{
		Category: CategoryNodeBuild,
		Name:     "node.build.script",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_BUILD_SCRIPT"),
		),
		Usage:       fmt.Sprintf("package.json script for the build operation. %s", environment.HelpFormatTemplate),
		Required:    false,
		Value:       "build",
		Destination: &P.Build.Script,
	},

	&cli.StringFlag{
		Category: CategoryNodeBuild,
		Name:     "node.build.script-args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_BUILD_SCRIPT_ARGS"),
		),
		Usage:       fmt.Sprintf("package.json script arguments for the build operation. %s", environment.HelpFormatTemplate),
		Required:    false,
		Value:       "",
		Destination: &P.Build.ScriptArgs,
	},

	&cli.StringFlag{
		Category: CategoryNodeBuild,
		Name:     "node.build.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("NODE_BUILD_CWD"),
		),
		Usage:       "Working directory for the build operation.",
		Required:    false,
		Value:       ".",
		Destination: &P.Build.Cwd,
	},
}
