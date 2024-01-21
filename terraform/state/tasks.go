package state

import (
	"fmt"
	"net/url"
	"strings"

	"gitlab.kilic.dev/devops/pipes/terraform/setup"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func GenerateTerraformEnvVarsState(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("state").
		Set(func(t *Task[Pipe]) error {
			if t.Pipe.State.Strict && t.Pipe.State.Type == "" {
				return fmt.Errorf("State has to be setup when in strict mode.")
			}

			return nil
		}).
		SetJobWrapper(func(job Job, t *Task[Pipe]) Job {
			return t.TL.JobParallel(
				job,
				GenerateTerraformEnvVarsGitlabState(t.TL).Job(),
			)
		})
}

func GenerateTerraformEnvVarsGitlabState(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("state", "gitlab-http").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.State.Type != TF_STATE_TYPE_GITLAB_HTTP
		}).
		Set(func(t *Task[Pipe]) error {
			if t.Pipe.GitlabHttpState.HttpAddress == "" {
				t.Pipe.GitlabHttpState.HttpAddress = strings.Join([]string{
					setup.TL.Pipe.ApiUrl,
					"projects",
					url.PathEscape(setup.TL.Pipe.CiVariables.ProjectId),
					"terraform/state",
					url.PathEscape(t.Pipe.State.Name),
				},
					"/",
				)

				t.Log.Debugf("State HTTP address has not been set, using default for state type: %s", t.Pipe.GitlabHttpState.HttpAddress)
			}

			setup.TL.Pipe.Ctx.EnvVars["TF_HTTP_ADDRESS"] = t.Pipe.GitlabHttpState.HttpAddress

			if t.Pipe.GitlabHttpState.HttpLockAddress == "" {
				t.Pipe.GitlabHttpState.HttpLockAddress = strings.Join([]string{
					t.Pipe.GitlabHttpState.HttpAddress,
					"lock",
				}, "/")

				t.Log.Debugf("State HTTP lock address has not been set, using default for state type: %s", t.Pipe.GitlabHttpState.HttpLockAddress)
			}

			setup.TL.Pipe.Ctx.EnvVars["TF_HTTP_LOCK_ADDRESS"] = t.Pipe.GitlabHttpState.HttpLockAddress
			setup.TL.Pipe.Ctx.EnvVars["TF_HTTP_LOCK_METHOD"] = t.Pipe.GitlabHttpState.HttpLockMethod

			if t.Pipe.GitlabHttpState.HttpUnlockAddress == "" {
				t.Pipe.GitlabHttpState.HttpUnlockAddress = strings.Join([]string{
					t.Pipe.GitlabHttpState.HttpAddress,
					"lock",
				}, "/")

				t.Log.Debugf("State HTTP unlock address has not been set, using default for state type: %s", t.Pipe.GitlabHttpState.HttpUnlockAddress)
			}

			setup.TL.Pipe.Ctx.EnvVars["TF_HTTP_UNLOCK_ADDRESS"] = t.Pipe.GitlabHttpState.HttpUnlockAddress
			setup.TL.Pipe.Ctx.EnvVars["TF_HTTP_UNLOCK_METHOD"] = t.Pipe.GitlabHttpState.HttpUnlockMethod

			if t.Pipe.GitlabHttpState.HttpUsername == "" {
				return fmt.Errorf("TF_HTTP_USERNAME is required for state type: %s", t.Pipe.State.Type)
			}
			setup.TL.Pipe.Ctx.EnvVars["TF_HTTP_USERNAME"] = t.Pipe.GitlabHttpState.HttpUsername

			if t.Pipe.GitlabHttpState.HttpPassword == "" {
				return fmt.Errorf("TF_HTTP_PASSWORD is required for state type: %s", t.Pipe.State.Type)
			}
			setup.TL.Pipe.Ctx.EnvVars["TF_HTTP_PASSWORD"] = t.Pipe.GitlabHttpState.HttpPassword
			t.Plumber.AppendSecrets(t.Pipe.GitlabHttpState.HttpPassword)

			setup.TL.Pipe.Ctx.EnvVars["TF_HTTP_RETRY_WAIT_MIN"] = t.Pipe.GitlabHttpState.HttpRetryWaitMin

			t.Log.Debugf("Generated following environment variables for terraform to consume: %+v", setup.TL.Pipe.Ctx.EnvVars)

			return nil
		})
}
