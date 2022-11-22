package pipe

import (
	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/common/flags"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

//revive:disable:line-length-limit

const (
	CATEGORY_SEMANTIC_RELEASE = "Semantic Release"
)

var Flags = []cli.Flag{

	// CATEGORY_PACKAGES

	&cli.StringSliceFlag{
		Category: flags.CATEGORY_PACKAGES,
		Name:     "packages.apk",
		Usage:    "APK applications to install before running semantic-release.",
		Required: false,
		EnvVars:  []string{"ADD_APKS"},
		Value:    &cli.StringSlice{},
	},

	// CATEGORY_SEMANTIC_RELEASE

	&cli.BoolFlag{
		Category:    CATEGORY_SEMANTIC_RELEASE,
		Name:        "semantic_release.dry_run",
		Usage:       "Node packages to install before running semantic-release.",
		Required:    false,
		EnvVars:     []string{"SEMANTIC_RELEASE_DRY_RUN"},
		Value:       false,
		Destination: &TL.Pipe.SemanticRelease.IsDryRun,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_SEMANTIC_RELEASE,
		Name:        "semantic_release.run_multi",
		Usage:       "Use @qiwi/multi-semantic-release package to do a workspace release.",
		Required:    false,
		EnvVars:     []string{"SEMANTIC_RELEASE_RUN_MULTI"},
		Value:       false,
		Destination: &TL.Pipe.SemanticRelease.UseMulti,
	},
}

func ProcessFlags(tl *TaskList[Pipe]) error {
	tl.Pipe.Apk = tl.CliContext.StringSlice("packages.apk")

	return nil
}

var DeprecationNotices = []DeprecationNotice{
	{
		Environment: []string{"DRY_RUN", "RUN_MULTI"},
		Level:       LOG_LEVEL_ERROR,
		Message:     `"%s" is deprecated, please use the environment variable with the "SEMANTIC_RELEASE_" prefix instead.`,
	},
	{
		Environment: []string{"ADD_MODULES"},
		Level:       LOG_LEVEL_ERROR,
		Message:     `"%s" is deprecated, please use the environment variable "ADD_NODE_MODULES" through node pipe instead.`,
	},
}
