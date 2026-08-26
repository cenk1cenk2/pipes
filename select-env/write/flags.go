package write

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: environment.CATEGORY_ENVIRONMENT,
		Name:     "environment.file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("ENVIRONMENT_FILE"),
		),
		Usage:       "File for writing the environment variables for selected environment.",
		Required:    true,
		Value:       "env.environment",
		Destination: &P.Environment.File,
	},
}
