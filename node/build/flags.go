package build

import (
	"fmt"

	"github.com/urfave/cli/v2"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

//revive:disable:line-length-limit

const (
	CATEGORY_NODE_BUILD = "Build"
)

var Flags = []cli.Flag{
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
}

var DeprecationNotices = []DeprecationNotice{
	{
		Flag:        []string{"--node.build_environment_files", "--node.build_environment_fallback", "--node.build_environment_conditions"},
		Environment: []string{"NODE_BUILD_ENVIRONMENT_FILES", "NODE_BUILD_ENVIRONMENT_CONDITIONS", "NODE_BUILD_ENVIRONMENT_FALLBACK"},
		Level:       LOG_LEVEL_ERROR,
		Message:     `"%s" is deprecated, please utilize the new select-env flags instead.`,
	},
}
