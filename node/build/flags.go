package build

import (
	"github.com/urfave/cli/v2"
)

//revive:disable:line-length-limit

const (
	category_git        = "GIT"
	category_node_build = "Build"
)

var Flags = []cli.Flag{

	// category_git

	&cli.StringFlag{
		Category:    category_git,
		Name:        "git.tag",
		Usage:       "Source control management tag.",
		Required:    false,
		EnvVars:     []string{"CI_COMMIT_TAG", "BITBUCKET_TAG"},
		Value:       "",
		Destination: &TL.Pipe.Git.Tag,
	},

	&cli.StringFlag{
		Category:    category_git,
		Name:        "git.branch",
		Usage:       "Source control management branch.",
		Required:    false,
		EnvVars:     []string{"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
		Value:       "",
		Destination: &TL.Pipe.Git.Branch,
	},

	// category_build

	&cli.StringFlag{
		Category:    category_node_build,
		Name:        "node.build_script",
		Usage:       "package.json script for building operation.",
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_SCRIPT"},
		Value:       "build",
		Destination: &TL.Pipe.NodeBuild.Script,
	},

	&cli.StringFlag{
		Category:    category_node_build,
		Name:        "node.build_script_args",
		Usage:       "package.json script arguments for building operation.",
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_SCRIPT_ARGS"},
		Value:       "",
		Destination: &TL.Pipe.NodeBuild.ScriptArgs,
	},

	&cli.StringFlag{
		Category:    category_node_build,
		Name:        "node.build_cwd",
		Usage:       "Working directory for build operation.",
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_CWD"},
		Value:       ".",
		Destination: &TL.Pipe.NodeBuild.Cwd,
	},

	&cli.StringSliceFlag{
		Category:    category_node_build,
		Name:        "node.build_environment_files",
		Usage:       "Yaml files to inject to build.",
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_ENVIRONMENT_FILES"},
		Value:       &cli.StringSlice{},
		Destination: &TL.Pipe.NodeBuild.EnvironmentFiles,
	},

	&cli.StringFlag{
		Category:    category_node_build,
		Name:        "node.build_environment_conditions",
		Usage:       "Tagging regex patterns to match. json(map[string]RegExp)",
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_ENVIRONMENT_CONDITIONS"},
		Value:       `{ "production": "^v\\d*\\.\\d*\\.\\d*$", "stage": "^v\\d*\\.\\d*\\.\\d*-.*$" }`,
		Destination: &TL.Pipe.NodeBuild.EnvironmentConditions,
	},

	&cli.StringFlag{
		Category:    category_node_build,
		Name:        "node.build_environment_fallback",
		Usage:       "Fallback, if it does not match any conditions. Defaults to current branch name.",
		Required:    false,
		EnvVars:     []string{"NODE_BUILD_ENVIRONMENT_FALLBACK"},
		Value:       "develop",
		Destination: &TL.Pipe.NodeBuild.EnvironmentFallback,
	},
}
