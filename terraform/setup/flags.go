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

	&cli.StringSliceFlag{
		Category: CATEGORY_PROJECT,
		Name:     "terraform-project.workspaces",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_WORKSPACES"),
		),
		Usage:       "Workspaces that this command will be executed on.",
		Required:    false,
		Destination: &P.Project.Workspaces,
	},

	// CATEGORY_CI_VARIABLES

	&cli.StringFlag{
		Category: CATEGORY_CI_VARIABLES,
		Name:     "terraform-var.job-id",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_VAR_CI_JOB_ID"),
			cli.EnvVar("CI_JOB_ID"),
		),
		Usage:       "Injected CI job-id variable to the deployment.",
		Required:    false,
		Value:       "",
		Destination: &P.CiVariables.JobId,
	},
	&cli.StringFlag{
		Category: CATEGORY_CI_VARIABLES,
		Name:     "terraform-var.commit-sha",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_VAR_CI_COMMIT_SHA"),
			cli.EnvVar("CI_COMMIT_SHA"),
		),
		Usage:       "Injected CI commit-sha variable to the deployment.",
		Required:    false,
		Value:       "",
		Destination: &P.CiVariables.CommitSha,
	},
	&cli.StringFlag{
		Category: CATEGORY_CI_VARIABLES,
		Name:     "terraform-var.job-stage",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_VAR_CI_JOB_STAGE"),
			cli.EnvVar("CI_JOB_STAGE"),
		),
		Usage:       "Injected CI job-stage variable to the deployment.",
		Required:    false,
		Value:       "",
		Destination: &P.CiVariables.JobStage,
	},
	&cli.StringFlag{
		Category: CATEGORY_CI_VARIABLES,
		Name:     "terraform-var.project-id",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_VAR_CI_PROJECT_ID"),
			cli.EnvVar("CI_PROJECT_ID"),
		),
		Usage:       "Injected CI project-id variable to the deployment.",
		Required:    false,
		Value:       "",
		Destination: &P.CiVariables.ProjectId,
	},
	&cli.StringFlag{
		Category: CATEGORY_CI_VARIABLES,
		Name:     "terraform-var.project-name",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_VAR_CI_PROJECT_NAME"),
			cli.EnvVar("CI_PROJECT_NAME"),
		),
		Usage:       "Injected CI project-name variable to the deployment.",
		Required:    false,
		Value:       "",
		Destination: &P.CiVariables.ProjectName,
	},
	&cli.StringFlag{
		Category: CATEGORY_CI_VARIABLES,
		Name:     "terraform-var.project-namespace",
		Usage:    "Injected CI project-namespace variable to the deployment.",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_VAR_CI_PROJECT_NAMESPACE"),
			cli.EnvVar("CI_PROJECT_NAMESPACE"),
		),
		Required:    false,
		Value:       "",
		Destination: &P.CiVariables.ProjectNamespace,
	},
	&cli.StringFlag{
		Category: CATEGORY_CI_VARIABLES,
		Name:     "terraform-var.project-path",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_VAR_CI_PROJECT_PATH"),
			cli.EnvVar("CI_PROJECT_PATH"),
		),
		Usage:       "Injected CI project-path variable to the deployment.",
		Required:    false,
		Value:       "",
		Destination: &P.CiVariables.ProjectPath,
	},
	&cli.StringFlag{
		Category: CATEGORY_CI_VARIABLES,
		Name:     "terraform-var.project-url",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_VAR_CI_PROJECT_URL"),
			cli.EnvVar("CI_PROJECT_URL"),
		),
		Usage:       "Injected CI project-url variable to the deployment.",
		Required:    false,
		Value:       "",
		Destination: &P.CiVariables.ProjectUrl,
	},
	&cli.StringFlag{
		Category: CATEGORY_CI_VARIABLES,
		Name:     "terraform-var.api-url",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_VAR_CI_API_V4_URL"),
			cli.EnvVar("CI_API_V4_URL"),
		),
		Usage:       "Injected CI api-url variable to the deployment.",
		Required:    false,
		Value:       "",
		Destination: &P.CiVariables.ApiUrl,
	},
}
