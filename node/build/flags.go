package build

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "git.branch",
		Usage:       "Source control management branch.",
		Required:    false,
		EnvVars:     []string{"CI_COMMIT_REF_NAME"},
		Value:       "",
		Destination: &P.Pipe.Git.Branch,
	},

	&cli.StringFlag{
		Name:        "git.tag",
		Usage:       "Source control management tag.",
		Required:    false,
		EnvVars:     []string{"CI_COMMIT_TAG"},
		Value:       "",
		Destination: &P.Pipe.Git.Tag,
	},

	&cli.StringFlag{
		Name:        "node.build_script",
		Usage:       "package.json script for building operation.",
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_SCRIPT"},
		Value:       "build",
		Destination: &P.Pipe.NodeBuild.Script,
	},

	&cli.StringFlag{
		Name:        "node.build_script_args",
		Usage:       "package.json script arguments for building operation.",
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_SCRIPT_ARGS"},
		Value:       "",
		Destination: &P.Pipe.NodeBuild.ScriptArgs,
	},

	&cli.StringFlag{
		Name:        "node.build_cwd",
		Usage:       "Working directory for build operation.",
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_CWD"},
		Value:       ".",
		Destination: &P.Pipe.NodeBuild.Cwd,
	},

	&cli.StringSliceFlag{
		Name:        "node.build_environment_files",
		Usage:       "Yaml files to inject to build.",
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_ENVIRONMENT_FILES"},
		Value:       &cli.StringSlice{},
		Destination: &P.Pipe.NodeBuild.EnvironmentFiles,
	},

	&cli.StringFlag{
		Name:        "node.build_environment_conditions",
		Usage:       "Tagging regex patterns to match. json({ [name: string]: RegExp })",
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_ENVIRONMENT_CONDITIONS"},
		Value:       `{ "production": "^v\\d*\\.\\d*\\.\\d*$", "stage": "^v\\d*\\.\\d*\\.\\d*-.*$" }`,
		Destination: &P.Pipe.NodeBuild.EnvironmentConditions,
	},

	&cli.StringFlag{
		Name:        "node.build_environment_fallback",
		Usage:       "Fallback, if it does not match any conditions. Defaults to current branch name.",
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_ENVIRONMENT_FALLBACK"},
		Value:       "develop",
		Destination: &P.Pipe.NodeBuild.EnvironmentFallback,
	},
}
