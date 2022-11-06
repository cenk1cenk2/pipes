package build

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/common/flags"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
)

//revive:disable:line-length-limit

const (
	CATEGORY_NODE_BUILD = "Build"
)

var Flags = TL.Plumber.AppendFlags(flags.NewGitFlags(
	flags.GitFlagsSetup{
		GitBranchDestination: &TL.Pipe.Git.Branch,
		GitTagDestination:    &TL.Pipe.Git.Tag,
	},
), []cli.Flag{
	// CATEGORY_BUILD

	&cli.StringFlag{
		Category:    CATEGORY_NODE_BUILD,
		Name:        "node.build_script",
		Usage:       fmt.Sprintf("package.json script for building operation. format(%s)", environment.HELP_FORMAT_ENVIRONMENT_TEMPLATE),
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_SCRIPT"},
		Value:       "build",
		Destination: &TL.Pipe.NodeBuild.Script,
	},

	&cli.StringFlag{
		Category:    CATEGORY_NODE_BUILD,
		Name:        "node.build_script_args",
		Usage:       fmt.Sprintf("package.json script arguments for building operation. format(%s)", environment.HELP_FORMAT_ENVIRONMENT_TEMPLATE),
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_SCRIPT_ARGS"},
		Value:       "",
		Destination: &TL.Pipe.NodeBuild.ScriptArgs,
	},

	&cli.StringFlag{
		Category:    CATEGORY_NODE_BUILD,
		Name:        "node.build_cwd",
		Usage:       "Working directory for build operation.",
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_CWD"},
		Value:       ".",
		Destination: &TL.Pipe.NodeBuild.Cwd,
	},
})
