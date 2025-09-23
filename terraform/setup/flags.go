package setup

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_CONFIG       = "Config"
	CATEGORY_PROJECT      = "Project"
	CATEGORY_CI_VARIABLES = "Injected Variables"
)

var Flags = []cli.Flag{
	// CATEGORY_CONFIG
	&cli.StringFlag{
		Category: CATEGORY_CONFIG,
		Name:     "terraform-config.log-level",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_LOG_LEVEL"),
			cli.EnvVar("TF_LOG"),
		),
		Usage:       `Terraform log level. enum("trace", "debug", "info", "warn", "error")`,
		Required:    false,
		Value:       "",
		Destination: &P.Config.LogLevel,
	},

	// CATEGORY_PROJECT

	&cli.StringFlag{
		Category: CATEGORY_PROJECT,
		Name:     "terraform-project.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_ROOT"),
		),
		Usage:       "Terraform project working directory",
		Required:    false,
		Value:       ".",
		Destination: &P.Project.Cwd,
	},
}
