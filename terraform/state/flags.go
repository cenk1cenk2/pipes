package state

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategoryState = "State"
)

var Flags = []cli.Flag{
	// CategoryState

	&cli.StringFlag{
		Category: CategoryState,
		Name:     "terraform.state.type",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_STATE_TYPE"),
			cli.EnvVar("TF_STATE_TYPE"),
		),
		Usage:       `Terraform state type. format(enum("gitlab-http"))`,
		Required:    false,
		Value:       "",
		Destination: &P.State.Type,
	},

	&cli.StringFlag{
		Category: CategoryState,
		Name:     "terraform.state.name",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_STATE_NAME"),
			cli.EnvVar("TF_STATE_NAME"),
		),
		Usage:       "Terraform state name.",
		Required:    false,
		Value:       "default",
		Destination: &P.State.Name,
	},

	&cli.BoolFlag{
		Category: CategoryState,
		Name:     "terraform.state.strict",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_STATE_STRICT"),
			cli.EnvVar("TF_STATE_STRICT"),
		),
		Usage:       "Fail when no Terraform state type is configured.",
		Required:    false,
		Value:       false,
		Destination: &P.State.Strict,
	},

	// gitlab http state

	&cli.StringFlag{
		Category: CategoryState,
		Name:     "terraform.state.gitlab-http.http-address",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_STATE_GITLAB_HTTP_HTTP_ADDRESS"),
			cli.EnvVar("TF_HTTP_ADDRESS"),
			cli.EnvVar("TF_ADDRESS"),
		),
		Usage:       "HTTP address for the GitLab HTTP state backend.",
		Required:    false,
		Value:       "",
		Destination: &P.GitlabHttpState.HttpAddress,
	},
	&cli.StringFlag{
		Category: CategoryState,
		Name:     "terraform.state.gitlab-http.http-lock-address",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_STATE_GITLAB_HTTP_HTTP_LOCK_ADDRESS"),
			cli.EnvVar("TF_HTTP_LOCK_ADDRESS"),
		),
		Usage:       "HTTP lock address for the GitLab HTTP state backend.",
		Required:    false,
		Value:       "",
		Destination: &P.GitlabHttpState.HttpLockAddress,
	},
	&cli.StringFlag{
		Category: CategoryState,
		Name:     "terraform.state.gitlab-http.http-lock-method",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_STATE_GITLAB_HTTP_HTTP_LOCK_METHOD"),
			cli.EnvVar("TF_HTTP_LOCK_METHOD"),
		),
		Usage:       "HTTP lock method for the GitLab HTTP state backend.",
		Required:    false,
		Value:       "POST",
		Destination: &P.GitlabHttpState.HttpLockMethod,
	},
	&cli.StringFlag{
		Category: CategoryState,
		Name:     "terraform.state.gitlab-http.http-unlock-address",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_STATE_GITLAB_HTTP_HTTP_UNLOCK_ADDRESS"),
			cli.EnvVar("TF_HTTP_UNLOCK_ADDRESS"),
		),
		Usage:       "HTTP unlock address for the GitLab HTTP state backend.",
		Required:    false,
		Value:       "",
		Destination: &P.GitlabHttpState.HttpUnlockAddress,
	},
	&cli.StringFlag{
		Category: CategoryState,
		Name:     "terraform.state.gitlab-http.http-unlock-method",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_STATE_GITLAB_HTTP_HTTP_UNLOCK_METHOD"),
			cli.EnvVar("TF_HTTP_UNLOCK_METHOD"),
		),
		Usage:       "HTTP unlock method for the GitLab HTTP state backend.",
		Required:    false,
		Value:       "DELETE",
		Destination: &P.GitlabHttpState.HttpUnlockMethod,
	},
	&cli.StringFlag{
		Category: CategoryState,
		Name:     "terraform.state.gitlab-http.http-username",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_STATE_GITLAB_HTTP_HTTP_USERNAME"),
			cli.EnvVar("TF_HTTP_USERNAME"),
			cli.EnvVar("TF_USERNAME"),
		),
		Usage:       "HTTP username for the GitLab HTTP state backend.",
		Required:    false,
		Value:       "gitlab-ci-token",
		Destination: &P.GitlabHttpState.HttpUsername,
	},
	&cli.StringFlag{
		Category: CategoryState,
		Name:     "terraform.state.gitlab-http.http-password",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_STATE_GITLAB_HTTP_HTTP_PASSWORD"),
			cli.EnvVar("TF_HTTP_PASSWORD"),
			cli.EnvVar("TF_PASSWORD"),
			cli.EnvVar("CI_JOB_TOKEN"),
		),
		Usage:       "HTTP password for the GitLab HTTP state backend.",
		Required:    false,
		Value:       "",
		Destination: &P.GitlabHttpState.HttpPassword,
	},
	&cli.StringFlag{
		Category: CategoryState,
		Name:     "terraform.state.gitlab-http.http-retry-wait-min",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_STATE_GITLAB_HTTP_HTTP_RETRY_WAIT_MIN"),
			cli.EnvVar("TF_HTTP_RETRY_WAIT_MIN"),
		),
		Usage:       "Minimum time to wait between HTTP retries for the GitLab HTTP state backend, in seconds.",
		Required:    false,
		Value:       "5",
		Destination: &P.GitlabHttpState.HttpRetryWaitMin,
	},
}
