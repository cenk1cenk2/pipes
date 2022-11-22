package pipe

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

//revive:disable:line-length-limit

const (
	CATEGORY_PACKAGES         = "Packages"
	CATEGORY_SEMANTIC_RELEASE = "Semantic Release"
)

var Flags = []cli.Flag{

	// CATEGORY_PACKAGES

	&cli.StringSliceFlag{
		Category: CATEGORY_PACKAGES,
		Name:     "packages.apk",
		Usage:    "APK applications to install before running semantic-release.",
		Required: false,
		EnvVars:  []string{"ADD_APKS"},
		Value:    &cli.StringSlice{},
	},

	&cli.StringSliceFlag{
		Category: CATEGORY_PACKAGES,
		Name:     "packages.node",
		Usage:    "Node packages to install before running semantic-release.",
		Required: false,
		EnvVars:  []string{"ADD_MODULES"},
		Value:    &cli.StringSlice{},
	},

	// CATEGORY_SEMANTIC_RELEASE

	&cli.BoolFlag{
		Category:    CATEGORY_SEMANTIC_RELEASE,
		Name:        "semantic_release.dry_run",
		Usage:       "Node packages to install before running semantic-release.",
		Required:    false,
		EnvVars:     []string{"DRY_RUN"},
		Value:       false,
		Destination: &TL.Pipe.SemanticRelease.IsDryRun,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_SEMANTIC_RELEASE,
		Name:        "semantic_release.run_multi",
		Usage:       "Use @qiwi/multi-semantic-release package to do a workspace release.",
		Required:    false,
		EnvVars:     []string{"RUN_MULTI"},
		Value:       false,
		Destination: &TL.Pipe.SemanticRelease.UseMulti,
	},
}

func ProcessFlags(tl *TaskList[Pipe]) error {
	tl.Pipe.Apk = tl.CliContext.StringSlice("packages.apk")
	tl.Pipe.Node = tl.CliContext.StringSlice("packages.node")

	return nil
}
