package setup

import (
	"regexp"

	ucli "github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
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
	CwdEnvAliases:  []string{"TF_ROOT"},
	VersionArgs:    []string{"version"},
	VersionPattern: regexp.MustCompile(`Terraform (v\d+\.\d+\.\d+)`),
}

var Flags = CombineFlags(
	tool.Flags(Spec, &P.Config),
	[]ucli.Flag{
		// CATEGORY_CONFIG
		&ucli.StringFlag{
			Category:    CATEGORY_CONFIG,
			Name:        "terraform-config.log-level",
			Sources:     cli.EnvVars("TF_LOG_LEVEL", "TF_LOG"),
			Usage:       `Terraform log level. enum("trace", "debug", "info", "warn", "error")`,
			Required:    false,
			Value:       "",
			Destination: &P.LogLevel,
		},

		// CATEGORY_CI_VARIABLES

		&ucli.StringFlag{
			Category:    CATEGORY_CI_VARIABLES,
			Name:        "terraform-var.api-url",
			Sources:     cli.EnvVars("TF_VAR_CI_API_V4_URL", "CI_API_V4_URL"),
			Usage:       "Injected CI api-url variable to the deployment.",
			Required:    false,
			Value:       "",
			Destination: &P.CiVariables.ApiUrl,
		},

		&ucli.StringFlag{
			Category:    CATEGORY_CI_VARIABLES,
			Name:        "terraform-var.project-id",
			Sources:     cli.EnvVars("TF_VAR_CI_PROJECT_ID", "CI_PROJECT_ID"),
			Usage:       "Injected CI project-id variable to the deployment.",
			Required:    false,
			Value:       "",
			Destination: &P.CiVariables.ProjectId,
		},
	})
