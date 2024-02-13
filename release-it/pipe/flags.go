package pipe

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

const (
	CATEGORY_RELEASE_IT = "release-it"
)

var Flags = []cli.Flag{

	// CATEGORY_RELEASE_IT

	&cli.BoolFlag{
		Category:    CATEGORY_RELEASE_IT,
		Name:        "release-it.dry_run",
		Usage:       "Run release-it in dry mode without making changes.",
		Required:    false,
		EnvVars:     []string{"RELEASE_IT_DRY_RUN"},
		Value:       false,
		Destination: &TL.Pipe.ReleaseIt.IsDryRun,
	},

	&cli.StringFlag{
		Category:    CATEGORY_RELEASE_IT,
		Name:        "release-it.config-file",
		Usage:       "release-it configuration file.",
		Required:    false,
		EnvVars:     []string{"RELEASE_IT_CONFIG_FILE"},
		Value:       "",
		Destination: &TL.Pipe.ReleaseIt.ConfigFile,
	},
}

func ProcessFlags(tl *TaskList[Pipe]) error {
	return nil
}
