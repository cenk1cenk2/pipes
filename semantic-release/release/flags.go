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
		Usage:       "Use @anolilab/multi-semantic-release package to do a workspace release.",
		Required:    false,
		Value:       false,
		Destination: &P.SemanticRelease.Workspace,
	},

	&cli.StringFlag{
		Category: CategorySemanticRelease,
		Name:     "semantic-release.tags-file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TAGS_FILE"),
		),
		Usage:       "File the release writes its tag into, which is handed over to the later jobs as a dotenv report.",
		Required:    false,
		Value:       ".tags",
		Destination: &P.SemanticRelease.TagsFile,
	},

	&cli.StringFlag{
		Category: CategorySemanticRelease,
		Name:     "semantic-release.env-file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("SEMANTIC_RELEASE_ENV_FILE"),
		),
		Usage:       "Dotenv report the released tags are written into for the later jobs.",
		Required:    false,
		Value:       "semantic-release.env",
		Destination: &P.SemanticRelease.EnvFile,
	},

	&cli.StringFlag{
		Category: CategorySemanticRelease,
		Name:     "semantic-release.env-variable",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("SEMANTIC_RELEASE_ENV_VARIABLE"),
		),
		Usage:       "Variable the released tag is handed over in, suffixed with the package name of every released package in a workspace.",
		Required:    false,
		Value:       "SEMANTIC_RELEASE_TAG",
		Destination: &P.SemanticRelease.EnvVariable,
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
