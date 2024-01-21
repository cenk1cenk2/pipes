package state

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

const (
	CATEGORY_STATE = "State"
)

var Flags = []cli.Flag{
	// CATEGORY_STATE

	&cli.StringFlag{
		Category:    CATEGORY_STATE,
		Name:        "terraform-state.type",
		Usage:       `Terraform state type. enum("gitlab-http")`,
		Required:    false,
		EnvVars:     []string{"TF_STATE_TYPE"},
		Value:       "",
		Destination: &TL.Pipe.State.Type,
	},

	&cli.StringFlag{
		Category:    CATEGORY_STATE,
		Name:        "terraform-state.name",
		Usage:       "Terraform state name.",
		Required:    false,
		EnvVars:     []string{"TF_STATE_NAME"},
		Value:       "default",
		Destination: &TL.Pipe.State.Name,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_STATE,
		Name:        "terraform-state.strict",
		Usage:       "Terraform state strict.",
		Required:    true,
		EnvVars:     []string{"TF_STATE_STRICT"},
		Value:       false,
		Destination: &TL.Pipe.State.Strict,
	},

	// gitlab http state

	&cli.StringFlag{
		Category:    CATEGORY_STATE,
		Name:        "terraform-state.gitlab-http.http-address",
		Usage:       "State configuration for terraform: http-address",
		Required:    false,
		EnvVars:     []string{"TF_HTTP_ADDRESS", "TF_ADDRESS"},
		Value:       "",
		Destination: &TL.Pipe.GitlabHttpState.HttpAddress,
	},
	&cli.StringFlag{
		Category:    CATEGORY_STATE,
		Name:        "terraform-state.gitlab-http.http-lock-address",
		Usage:       "State configuration for terraform: http-lock-address",
		Required:    false,
		EnvVars:     []string{"TF_HTTP_LOCK_ADDRESS"},
		Value:       "",
		Destination: &TL.Pipe.GitlabHttpState.HttpLockAddress,
	},
	&cli.StringFlag{
		Category:    CATEGORY_STATE,
		Name:        "terraform-state.gitlab-http.http-lock-method",
		Usage:       "State configuration for terraform: http-lock-method",
		Required:    false,
		EnvVars:     []string{"TF_HTTP_LOCK_METHOD"},
		Value:       "POST",
		Destination: &TL.Pipe.GitlabHttpState.HttpLockMethod,
	},
	&cli.StringFlag{
		Category:    CATEGORY_STATE,
		Name:        "terraform-state.gitlab-http.http-unlock-address",
		Usage:       "State configuration for terraform: http-unlock-address",
		Required:    false,
		EnvVars:     []string{"TF_HTTP_UNLOCK_ADDRESS"},
		Value:       "",
		Destination: &TL.Pipe.GitlabHttpState.HttpUnlockAddress,
	},
	&cli.StringFlag{
		Category:    CATEGORY_STATE,
		Name:        "terraform-state.gitlab-http.http-unlock-method",
		Usage:       "State configuration for terraform: http-unlock-method",
		Required:    false,
		EnvVars:     []string{"TF_HTTP_UNLOCK_METHOD"},
		Value:       "DELETE",
		Destination: &TL.Pipe.GitlabHttpState.HttpUnlockMethod,
	},
	&cli.StringFlag{
		Category:    CATEGORY_STATE,
		Name:        "terraform-state.gitlab-http.http-username",
		Usage:       "State configuration for terraform: http-username",
		Required:    false,
		EnvVars:     []string{"TF_HTTP_USERNAME", "TF_USERNAME"},
		Value:       "gitlab-ci-token",
		Destination: &TL.Pipe.GitlabHttpState.HttpUsername,
	},
	&cli.StringFlag{
		Category:    CATEGORY_STATE,
		Name:        "terraform-state.gitlab-http.http-password",
		Usage:       "State configuration for terraform: http-password",
		Required:    false,
		EnvVars:     []string{"TF_HTTP_PASSWORD", "TF_PASSWORD", "CI_JOB_TOKEN"},
		Value:       "",
		Destination: &TL.Pipe.GitlabHttpState.HttpPassword,
	},
	&cli.StringFlag{
		Category:    CATEGORY_STATE,
		Name:        "terraform-state.gitlab-http.http-retry-wait-min",
		Usage:       "State configuration for terraform: http-retry-wait-min",
		Required:    false,
		EnvVars:     []string{"TF_HTTP_RETRY_WAIT_MIN"},
		Value:       "5",
		Destination: &TL.Pipe.GitlabHttpState.HttpRetryWaitMin,
	},
}

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList[Pipe]) error {
	return nil
}
