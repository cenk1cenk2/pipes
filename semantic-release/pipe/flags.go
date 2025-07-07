package pipe

import (
	"github.com/urfave/cli/v3"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

const (
	CATEGORY_SEMANTIC_RELEASE = "Semantic Release"
)

var Flags = []cli.Flag{

	// CATEGORY_SEMANTIC_RELEASE

	&cli.BoolFlag{
		Category:    CATEGORY_SEMANTIC_RELEASE,
		Name:        "semantic_release.dry_run",
		Usage:       "Run semantic-release in dry mode without making changes.",
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
	tl.Plumber.SetDeprecationNotices([]DeprecationNotice{
		{
			Flag:        []string{"packages.apk"},
			Environment: []string{"ADD_APKS"},
			Level:       LOG_LEVEL_ERROR,
			Message:     "Installing os packages directly through this pipe is deprecated whereas argument has been found: %s",
		},
		{
			Flag:        []string{"packages.node", "packages.node.global"},
			Environment: []string{"PACKAGES_NODE", "PACKAGES_NODE_GLOBAL"},
			Level:       LOG_LEVEL_ERROR,
			Message:     "Installing node packages directly through this pipe is deprecated whereas argument has been found: %s",
		},
	})

	return nil
}
