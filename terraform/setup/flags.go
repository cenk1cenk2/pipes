package setup

import (
	"regexp"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
	"gitlab.kilic.dev/devops/pipes/internal/tool"

	. "github.com/cenk1cenk2/plumber/v6"
)

//revive:disable:line-length-limit

const (
	CATEGORY_CONFIG       = "Config"
	CATEGORY_PROJECT      = "Project"
	CATEGORY_CI_VARIABLES = "Injected Variables"
)

var Spec = tool.Spec{
	Name:           "terraform",
	Category:       CATEGORY_PROJECT,
	FlagPrefix:     "terraform",
	EnvPrefix:      "TERRAFORM",
	VersionArgs:    []string{"version"},
	VersionPattern: regexp.MustCompile(`Terraform (v\d+\.\d+\.\d+)`),
}

var Flags = CombineFlags(
	tool.NewFlags(Spec, &P.Config),
	[]cli.Flag{
		// CATEGORY_CONFIG
		&cli.StringFlag{
			Category:    CATEGORY_CONFIG,
			Name:        "terraform.log-level",
			Sources:     flags.EnvVars("TERRAFORM_LOG_LEVEL", "TF_LOG"),
			Usage:       `Terraform log level. enum("trace", "debug", "info", "warn", "error")`,
			Required:    false,
			Value:       "",
			Destination: &P.LogLevel,
		},

		// CATEGORY_CI_VARIABLES

		&cli.StringFlag{
			Category:    CATEGORY_CI_VARIABLES,
			Name:        "terraform.ci.api-url",
			Sources:     flags.EnvVars("TERRAFORM_CI_API_URL", "CI_API_V4_URL"),
			Usage:       "Injected CI api-url variable to the deployment.",
			Required:    false,
			Value:       "",
			Destination: &P.CiVariables.ApiUrl,
		},

		&cli.StringFlag{
			Category:    CATEGORY_CI_VARIABLES,
			Name:        "terraform.ci.project-id",
			Sources:     flags.EnvVars("TERRAFORM_CI_PROJECT_ID", "CI_PROJECT_ID"),
			Usage:       "Injected CI project-id variable to the deployment.",
			Required:    false,
			Value:       "",
			Destination: &P.CiVariables.ProjectId,
		},
	})
