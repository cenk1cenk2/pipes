package pipe

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:        "packages.apk",
		Usage:       "APK applications to install before running semantic-release.",
		Required:    false,
		EnvVars:     []string{"ADD_APKS"},
		Value:       &cli.StringSlice{},
		Destination: &TL.Pipe.Packages.Apk,
	},

	&cli.StringSliceFlag{
		Name:        "packages.node",
		Usage:       "Node packages to install before running semantic-release.",
		Required:    false,
		EnvVars:     []string{"ADD_MODULES"},
		Value:       &cli.StringSlice{},
		Destination: &TL.Pipe.Packages.Node,
	},

	&cli.BoolFlag{
		Name:        "semantic_release.dry_run",
		Usage:       "Node packages to install before running semantic-release.",
		Required:    false,
		EnvVars:     []string{"DRY_RUN"},
		Value:       false,
		Destination: &TL.Pipe.SemanticRelease.IsDryRun,
	},

	&cli.BoolFlag{
		Name:        "semantic_release.run_multi",
		Usage:       "Uses @qiwi/multi-semantic-release package to do a workspace release.",
		Required:    false,
		EnvVars:     []string{"RUN_MULTI"},
		Value:       false,
		Destination: &TL.Pipe.SemanticRelease.UseMulti,
	},
}
