package pipe

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_SEMANTIC_RELEASE = "Semantic Release"
)

var Flags = []cli.Flag{

	// CATEGORY_SEMANTIC_RELEASE

	&cli.BoolFlag{
		Category: CATEGORY_SEMANTIC_RELEASE,
		Name:     "semantic_release.dry_run",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("SEMANTIC_RELEASE_DRY_RUN"),
		),
		Usage:       "Run semantic-release in dry mode without making changes.",
		Required:    false,
		Value:       false,
		Destination: &P.SemanticRelease.IsDryRun,
	},

	&cli.BoolFlag{
		Category: CATEGORY_SEMANTIC_RELEASE,
		Name:     "semantic_release.workspace",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("SEMANTIC_RELEASE_WORKSPACE"),
		),
		Usage:       "Use @qiwi/multi-semantic-release package to do a workspace release.",
		Required:    false,
		Value:       false,
		Destination: &P.SemanticRelease.Workspace,
	},
}
