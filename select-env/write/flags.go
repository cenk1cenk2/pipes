package write

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_ENVIRONMENT = "Environment"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_ENVIRONMENT,
		Name:     "environment.file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("ENVIRONMENT_FILE"),
		),
		Usage:       "File for writing the environment variables of the selected environment.",
		Required:    true,
		Value:       "env.environment",
		Destination: &P.Environment.File,
	},
}
