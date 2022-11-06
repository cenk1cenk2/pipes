package pipe

import (
	"github.com/urfave/cli/v2"
)

//revive:disable:line-length-limit

const (
	CATEGORY_PACKAGES         = "Packages"
	CATEGORY_SEMANTIC_RELEASE = "Semantic Release"
)

var Flags = []cli.Flag{

	// CATEGORY_PACKAGES

	&cli.StringSliceFlag{
		Category:    CATEGORY_PACKAGES,
		Name:        "packages.apk",
		Usage:       "APK applications to install before running semantic-release.",
		Required:    false,
		EnvVars:     []string{"ADD_APKS"},
		Value:       &cli.StringSlice{},
		Destination: &TL.Pipe.Packages.Apk,
	},

	&cli.StringSliceFlag{
		Category:    CATEGORY_PACKAGES,
		Name:        "packages.node",
		Usage:       "Node packages to install before running semantic-release.",
		Required:    false,
		EnvVars:     []string{"ADD_MODULES"},
		Value:       &cli.StringSlice{},
		Destination: &TL.Pipe.Packages.Node,
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
		Usage:       "Uses @qiwi/multi-semantic-release package to do a workspace release.",
		Required:    false,
		EnvVars:     []string{"RUN_MULTI"},
		Value:       false,
		Destination: &TL.Pipe.SemanticRelease.UseMulti,
	},
}
