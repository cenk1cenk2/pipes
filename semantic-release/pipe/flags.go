package pipe

import (
	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/common/flags"
	. "gitlab.kilic.dev/libraries/plumber/v5"
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
		Name:        "semantic_release.workspace",
		Usage:       "Use @qiwi/multi-semantic-release package to do a workspace release.",
		Required:    false,
		EnvVars:     []string{"SEMANTIC_RELEASE_WORKSPACE"},
		Value:       false,
		Destination: &TL.Pipe.SemanticRelease.Workspace,
	},
}

func ProcessFlags(tl *TaskList[Pipe]) error {
	tl.Pipe.Apk = tl.CliContext.StringSlice("packages.apk")

	return nil
}
