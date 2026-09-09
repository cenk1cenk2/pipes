package release

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategorySemanticRelease = "Semantic Release"
	CategoryCIVariables     = "CI Variables"
)

var Flags = []cli.Flag{

	// CategorySemanticRelease

	&cli.BoolFlag{
		Category: CategorySemanticRelease,
		Name:     "semantic-release.dry-run",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("SEMANTIC_RELEASE_DRY_RUN"),
		),
		Usage:       "Run semantic-release in dry mode without making changes.",
		Required:    false,
		Value:       false,
		Destination: &P.SemanticRelease.DryRun,
	},

	&cli.BoolFlag{
		Category: CategorySemanticRelease,
		Name:     "semantic-release.workspace",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("SEMANTIC_RELEASE_WORKSPACE"),
		),
		Usage:       "Use @qiwi/multi-semantic-release package to do a workspace release.",
		Required:    false,
		Value:       false,
		Destination: &P.SemanticRelease.Workspace,
	},

	// CategoryCIVariables

	&cli.StringFlag{
		Category: CategoryCIVariables,
		Name:     "semantic-release.ci.commit-reference",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("SEMANTIC_RELEASE_CI_COMMIT_REFERENCE"),
			cli.EnvVar("CI_COMMIT_REF_NAME"),
		),
		Usage:       "Current commit reference, either the branch or the tag name of the project.",
		Required:    false,
		Value:       "",
		Destination: &P.CI.CommitReference,
	},
}
