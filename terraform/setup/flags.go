package setup

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v5"
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
		Category:    CATEGORY_CONFIG,
		Name:        "terraform-config.log-level",
		Usage:       `Terraform log level. enum("trace", "debug", "info", "warn", "error")`,
		Required:    false,
		EnvVars:     []string{"TF_LOG_LEVEL", "TF_LOG"},
		Value:       "info",
		Destination: &TL.Pipe.Config.LogLevel,
	},

	// CATEGORY_PROJECT

	&cli.StringFlag{
		Category:    CATEGORY_PROJECT,
		Name:        "terraform-project.cwd",
		Usage:       "Terraform project working directory",
		Required:    false,
		EnvVars:     []string{"TF_ROOT"},
		Value:       ".",
		Destination: &TL.Pipe.Project.Cwd,
	},

	// CATEGORY_CI_VARIABLES

	&cli.StringFlag{
		Category:    CATEGORY_CI_VARIABLES,
		Name:        "terraform-var.job-id",
		Usage:       "Injected CI job-id variable to the deployment.",
		Required:    false,
		EnvVars:     []string{"TF_VAR_CI_JOB_ID", "CI_JOB_ID"},
		Value:       "",
		Destination: &TL.Pipe.CiVariables.JobId,
	},
	&cli.StringFlag{
		Category:    CATEGORY_CI_VARIABLES,
		Name:        "terraform-var.commit-sha",
		Usage:       "Injected CI commit-sha variable to the deployment.",
		Required:    false,
		EnvVars:     []string{"TF_VAR_CI_COMMIT_SHA", "CI_COMMIT_SHA"},
		Value:       "",
		Destination: &TL.Pipe.CiVariables.CommitSha,
	},
	&cli.StringFlag{
		Category:    CATEGORY_CI_VARIABLES,
		Name:        "terraform-var.job-stage",
		Usage:       "Injected CI job-stage variable to the deployment.",
		Required:    false,
		EnvVars:     []string{"TF_VAR_CI_JOB_STAGE", "CI_JOB_STAGE"},
		Value:       "",
		Destination: &TL.Pipe.CiVariables.JobStage,
	},
	&cli.StringFlag{
		Category:    CATEGORY_CI_VARIABLES,
		Name:        "terraform-var.project-id",
		Usage:       "Injected CI project-id variable to the deployment.",
		Required:    false,
		EnvVars:     []string{"TF_VAR_CI_PROJECT_ID", "CI_PROJECT_ID"},
		Value:       "",
		Destination: &TL.Pipe.CiVariables.ProjectId,
	},
	&cli.StringFlag{
		Category:    CATEGORY_CI_VARIABLES,
		Name:        "terraform-var.project-name",
		Usage:       "Injected CI project-name variable to the deployment.",
		Required:    false,
		EnvVars:     []string{"TF_VAR_CI_PROJECT_NAME", "CI_PROJECT_NAME"},
		Value:       "",
		Destination: &TL.Pipe.CiVariables.ProjectName,
	},
	&cli.StringFlag{
		Category:    CATEGORY_CI_VARIABLES,
		Name:        "terraform-var.project-namespace",
		Usage:       "Injected CI project-namespace variable to the deployment.",
		Required:    false,
		EnvVars:     []string{"TF_VAR_CI_PROJECT_NAMESPACE", "CI_PROJECT_NAMESPACE"},
		Value:       "",
		Destination: &TL.Pipe.CiVariables.ProjectNamespace,
	},
	&cli.StringFlag{
		Category:    CATEGORY_CI_VARIABLES,
		Name:        "terraform-var.project-path",
		Usage:       "Injected CI project-path variable to the deployment.",
		Required:    false,
		EnvVars:     []string{"TF_VAR_CI_PROJECT_PATH", "CI_PROJECT_PATH"},
		Value:       "",
		Destination: &TL.Pipe.CiVariables.ProjectPath,
	},
	&cli.StringFlag{
		Category:    CATEGORY_CI_VARIABLES,
		Name:        "terraform-var.project-url",
		Usage:       "Injected CI project-url variable to the deployment.",
		Required:    false,
		EnvVars:     []string{"TF_VAR_CI_PROJECT_URL", "CI_PROJECT_URL"},
		Value:       "",
		Destination: &TL.Pipe.CiVariables.ProjectUrl,
	},
	&cli.StringFlag{
		Category:    CATEGORY_CI_VARIABLES,
		Name:        "terraform-var.api-url",
		Usage:       "Injected CI api-url variable to the deployment.",
		Required:    false,
		EnvVars:     []string{"TF_VAR_CI_API_V4_URL", "CI_API_V4_URL"},
		Value:       "",
		Destination: &TL.Pipe.CiVariables.ApiUrl,
	},
}

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList[Pipe]) error {
	return nil
}
