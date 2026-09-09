package setup

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategoryConfig      = "Config"
	CategoryProject     = "Project"
	CategoryCIVariables = "Injected Variables"
)

var Flags = []cli.Flag{

	// CategoryProject

	&cli.StringFlag{
		Category: CategoryProject,
		Name:     "terraform.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_CWD"),
			cli.EnvVar("TF_ROOT"),
		),
		Usage:       "Working directory for terraform commands.",
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},

	// CategoryConfig

	&cli.StringFlag{
		Category: CategoryConfig,
		Name:     "terraform.log-level",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_LOG_LEVEL"),
			cli.EnvVar("TF_LOG_LEVEL"),
			cli.EnvVar("TF_LOG"),
		),
		Usage:       `Terraform log level. format(enum("trace", "debug", "info", "warn", "error"))`,
		Required:    false,
		Value:       "",
		Destination: &P.LogLevel,
	},

	// CategoryCIVariables

	&cli.StringFlag{
		Category: CategoryCIVariables,
		Name:     "terraform.ci.api-url",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_CI_API_URL"),
			cli.EnvVar("TF_VAR_CI_API_V4_URL"),
			cli.EnvVar("CI_API_V4_URL"),
		),
		Usage:       "Injected CI api-url variable to the deployment.",
		Required:    false,
		Value:       "",
		Destination: &P.CiVariables.ApiUrl,
	},

	&cli.StringFlag{
		Category: CategoryCIVariables,
		Name:     "terraform.ci.project-id",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_CI_PROJECT_ID"),
			cli.EnvVar("TF_VAR_CI_PROJECT_ID"),
			cli.EnvVar("CI_PROJECT_ID"),
		),
		Usage:       "Injected CI project-id variable to the deployment.",
		Required:    false,
		Value:       "",
		Destination: &P.CiVariables.ProjectId,
	},
}
