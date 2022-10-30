package pipe

import (
	"github.com/urfave/cli/v2"
)

//revive:disable:line-length-limit

const (
	category_packages         = "Packages"
	category_semantic_release = "Semantic Release"
)

var Flags = []cli.Flag{

	// category_packages

	&cli.StringSliceFlag{
		Category:    category_packages,
		Name:        "packages.apk",
		Usage:       "APK applications to install before running semantic-release.",
		Required:    false,
		EnvVars:     []string{"ADD_APKS"},
		Value:       &cli.StringSlice{},
		Destination: &TL.Pipe.Packages.Apk,
	},

	&cli.StringSliceFlag{
		Category:    category_packages,
		Name:        "packages.node",
		Usage:       "Node packages to install before running semantic-release.",
		Required:    false,
		EnvVars:     []string{"ADD_MODULES"},
		Value:       &cli.StringSlice{},
		Destination: &TL.Pipe.Packages.Node,
	},

	// category_semantic_release

	&cli.BoolFlag{
		Category:    category_semantic_release,
		Name:        "semantic_release.dry_run",
		Usage:       "Node packages to install before running semantic-release.",
		Required:    false,
		EnvVars:     []string{"DRY_RUN"},
		Value:       false,
		Destination: &TL.Pipe.SemanticRelease.IsDryRun,
	},

	&cli.BoolFlag{
		Category:    category_semantic_release,
		Name:        "semantic_release.run_multi",
		Usage:       "Uses @qiwi/multi-semantic-release package to do a workspace release.",
		Required:    false,
		EnvVars:     []string{"RUN_MULTI"},
		Value:       false,
		Destination: &TL.Pipe.SemanticRelease.UseMulti,
	},
}
