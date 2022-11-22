package flags

import (
	"github.com/urfave/cli/v2"
)

type SelectEnvEnableFlagSetup struct {
	SelectEnvEnableDestination *bool
	SelectEnvEnableRequired    bool
	SelectEnvEnableValue       bool
}

func NewSelectEnvEnableFlag(setup SelectEnvEnableFlagSetup) []cli.Flag {
	return []cli.Flag{

		// CATEGORY_ENVIRONMENT

		&cli.BoolFlag{
			Category:    CATEGORY_ENVIRONMENT,
			Name:        "environment.enable",
			Usage:       "Enable environment injection.",
			Required:    setup.SelectEnvEnableRequired,
			EnvVars:     []string{"ENVIRONMENT_ENABLE"},
			Value:       setup.SelectEnvEnableValue,
			Destination: setup.SelectEnvEnableDestination,
		},
	}
}
