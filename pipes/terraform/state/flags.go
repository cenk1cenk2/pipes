package state

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_STATE = "State"
)

var Flags = []cli.Flag{
	// CATEGORY_STATE

	&cli.StringFlag{
		Category: CATEGORY_STATE,
		Name:     "terraform-state.type",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_STATE_TYPE"),
		),
		Usage:       `Terraform state type. enum("gitlab-http")`,
		Required:    false,
		Value:       "",
		Destination: &P.State.Type,
	},

	&cli.StringFlag{
		Category: CATEGORY_STATE,
		Name:     "terraform-state.name",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_STATE_NAME"),
		),
		Usage:       "Terraform state name.",
		Required:    false,
		Value:       "default",
		Destination: &P.State.Name,
	},

	&cli.BoolFlag{
		Category: CATEGORY_STATE,
		Name:     "terraform-state.strict",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_STATE_STRICT"),
		),
		Usage:       "Terraform state strict.",
		Required:    false,
		Value:       false,
		Destination: &P.State.Strict,
	},

	// gitlab http state

	&cli.StringFlag{
		Category: CATEGORY_STATE,
		Name:     "terraform-state.gitlab-http.http-address",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_HTTP_ADDRESS"),
			cli.EnvVar("TF_ADDRESS"),
		),
		Usage:       "State configuration for terraform: http-address",
		Required:    false,
		Value:       "",
		Destination: &P.GitlabHttpState.HttpAddress,
	},
	&cli.StringFlag{
		Category: CATEGORY_STATE,
		Name:     "terraform-state.gitlab-http.http-lock-address",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_HTTP_LOCK_ADDRESS"),
		),
		Usage:       "State configuration for terraform: http-lock-address",
		Required:    false,
		Value:       "",
		Destination: &P.GitlabHttpState.HttpLockAddress,
	},
	&cli.StringFlag{
		Category: CATEGORY_STATE,
		Name:     "terraform-state.gitlab-http.http-lock-method",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_HTTP_LOCK_METHOD"),
		),
		Usage:       "State configuration for terraform: http-lock-method",
		Required:    false,
		Value:       "POST",
		Destination: &P.GitlabHttpState.HttpLockMethod,
	},
	&cli.StringFlag{
		Category: CATEGORY_STATE,
		Name:     "terraform-state.gitlab-http.http-unlock-address",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_HTTP_UNLOCK_ADDRESS"),
		),
		Usage:       "State configuration for terraform: http-unlock-address",
		Required:    false,
		Value:       "",
		Destination: &P.GitlabHttpState.HttpUnlockAddress,
	},
	&cli.StringFlag{
		Category: CATEGORY_STATE,
		Name:     "terraform-state.gitlab-http.http-unlock-method",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_HTTP_UNLOCK_METHOD"),
		),
		Usage:       "State configuration for terraform: http-unlock-method",
		Required:    false,
		Value:       "DELETE",
		Destination: &P.GitlabHttpState.HttpUnlockMethod,
	},
	&cli.StringFlag{
		Category: CATEGORY_STATE,
		Name:     "terraform-state.gitlab-http.http-username",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_HTTP_USERNAME"),
			cli.EnvVar("TF_USERNAME"),
		),
		Usage:       "State configuration for terraform: http-username",
		Required:    false,
		Value:       "gitlab-ci-token",
		Destination: &P.GitlabHttpState.HttpUsername,
	},
	&cli.StringFlag{
		Category: CATEGORY_STATE,
		Name:     "terraform-state.gitlab-http.http-password",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_HTTP_PASSWORD"),
			cli.EnvVar("TF_PASSWORD"),
			cli.EnvVar("CI_JOB_TOKEN"),
		),
		Usage:       "State configuration for terraform: http-password",
		Required:    false,
		Value:       "",
		Destination: &P.GitlabHttpState.HttpPassword,
	},
	&cli.StringFlag{
		Category: CATEGORY_STATE,
		Name:     "terraform-state.gitlab-http.http-retry-wait-min",
		Usage:    "State configuration for terraform: http-retry-wait-min",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_HTTP_RETRY_WAIT_MIN"),
		),
		Required:    false,
		Value:       "5",
		Destination: &P.GitlabHttpState.HttpRetryWaitMin,
	},
}
