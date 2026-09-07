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

	// CATEGORY_PROJECT

	&cli.StringFlag{
		Category:    CATEGORY_PROJECT,
		Name:        "terraform.cwd",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("TERRAFORM_CWD")),
		Usage:       "Working directory for terraform commands.",
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},

	// CATEGORY_CONFIG

	&cli.StringFlag{
		Category:    CATEGORY_CONFIG,
		Name:        "terraform.log-level",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("TERRAFORM_LOG_LEVEL"), cli.EnvVar("TF_LOG")),
		Usage:       `Terraform log level. enum("trace", "debug", "info", "warn", "error")`,
		Required:    false,
		Value:       "",
		Destination: &P.LogLevel,
	},

	// CATEGORY_CI_VARIABLES

	&cli.StringFlag{
		Category:    CATEGORY_CI_VARIABLES,
		Name:        "terraform.ci.api-url",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("TERRAFORM_CI_API_URL"), cli.EnvVar("CI_API_V4_URL")),
		Usage:       "Injected CI api-url variable to the deployment.",
		Required:    false,
		Value:       "",
		Destination: &P.CiVariables.ApiUrl,
	},

	&cli.StringFlag{
		Category:    CATEGORY_CI_VARIABLES,
		Name:        "terraform.ci.project-id",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("TERRAFORM_CI_PROJECT_ID"), cli.EnvVar("CI_PROJECT_ID")),
		Usage:       "Injected CI project-id variable to the deployment.",
		Required:    false,
		Value:       "",
		Destination: &P.CiVariables.ProjectId,
	},
}
