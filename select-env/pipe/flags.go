package pipe

import (
	"github.com/urfave/cli/v2"

	"gitlab.kilic.dev/devops/pipes/select-env/setup"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category:    setup.CATEGORY_ENVIRONMENT,
		Name:        "environment.file",
		Usage:       "File for writing the environment variables for selected environment.",
		Required:    true,
		EnvVars:     []string{"ENVIRONMENT_FILE"},
		Value:       "env.environment",
		Destination: &TL.Pipe.Environment.File,
	},
}
