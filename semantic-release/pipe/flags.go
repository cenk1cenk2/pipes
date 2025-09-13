package pipe

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_SEMANTIC_RELEASE = "Semantic Release"
	CATEGORY_CI_VARIABLES     = "CI Variables"
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

	// CATEGORY_CI_VARIABLES

	&cli.StringFlag{
		Category: CATEGORY_CI_VARIABLES,
		Name:     "ci.commit-ref-name",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CI_COMMIT_REF_NAME"),
		),
		Usage:       "Current commit reference that can be branch or tag name of the project..",
		Required:    false,
		Value:       "",
		Destination: &P.CI.CommitReference,
	},
}
