package build

import (
	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/common/flags"
)

//revive:disable:line-length-limit

const (
	CATEGORY_NODE_BUILD = "Build"
)

var Flags = TL.Plumber.AppendFlags(flags.NewGitFlags(flags.GitFlagsDestination{
	GitBranch: &TL.Pipe.Git.Branch,
	GitTag:    &TL.Pipe.Git.Tag,
}), []cli.Flag{
	// CATEGORY_BUILD

	&cli.StringFlag{
		Category:    CATEGORY_NODE_BUILD,
		Name:        "node.build_script",
		Usage:       "package.json script for building operation. format(Template(struct{ Environment: string }))",
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_SCRIPT"},
		Value:       "build",
		Destination: &TL.Pipe.NodeBuild.Script,
	},

	&cli.StringFlag{
		Category:    CATEGORY_NODE_BUILD,
		Name:        "node.build_script_args",
		Usage:       "package.json script arguments for building operation. format(Template(struct{ Environment: string }))",
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
