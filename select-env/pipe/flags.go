package pipe

import (
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/common/flags"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category:    flags.CATEGORY_ENVIRONMENT,
		Name:        "environment.file",
		Usage:       "File for writing the environment variables for selected environment.",
		Required:    true,
		EnvVars:     []string{"ENVIRONMENT_FILE"},
		Value:       "env.environment",
		Destination: &TL.Pipe.Environment.File,
	},
}
