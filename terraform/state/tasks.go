package state

import (
	"fmt"
	"net/url"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
)

func GenerateTerraformEnvVarsState(tl *TaskList) *Task {
	return tl.CreateTask("state").
		Set(func(t *Task) error {
			if P.State.Strict && P.State.Type == "" {
				return fmt.Errorf("State has to be setup when in strict mode.")
			}

			return nil
		}).
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobParallel(
				job,
				GenerateTerraformEnvVarsGitlabState(t.TL).Job(),
			)
		})
}

func GenerateTerraformEnvVarsGitlabState(tl *TaskList) *Task {
	return tl.CreateTask("state", "gitlab-http").
		ShouldDisable(func(t *Task) bool {
			return P.State.Type != TF_STATE_TYPE_GITLAB_HTTP
		}).
		Set(func(t *Task) error {
			if P.GitlabHttpState.HttpAddress == "" {
				P.GitlabHttpState.HttpAddress = strings.Join([]string{
					setup.P.ApiUrl,
					"projects",
					url.PathEscape(setup.P.CiVariables.ProjectId),
					"terraform/state",
					url.PathEscape(P.State.Name),
				},
					"/",
				)

				t.Log.Debugf("State HTTP address has not been set, using default for state type: %s", P.GitlabHttpState.HttpAddress)
			}

			setup.C.EnvVars["TF_HTTP_ADDRESS"] = P.GitlabHttpState.HttpAddress

			if P.GitlabHttpState.HttpLockAddress == "" {
				P.GitlabHttpState.HttpLockAddress = strings.Join([]string{
					P.GitlabHttpState.HttpAddress,
					"lock",
				}, "/")

				t.Log.Debugf("State HTTP lock address has not been set, using default for state type: %s", P.GitlabHttpState.HttpLockAddress)
			}

			setup.C.EnvVars["TF_HTTP_LOCK_ADDRESS"] = P.GitlabHttpState.HttpLockAddress
			setup.C.EnvVars["TF_HTTP_LOCK_METHOD"] = P.GitlabHttpState.HttpLockMethod

			if P.GitlabHttpState.HttpUnlockAddress == "" {
				P.GitlabHttpState.HttpUnlockAddress = strings.Join([]string{
					P.GitlabHttpState.HttpAddress,
					"lock",
				}, "/")

				t.Log.Debugf("State HTTP unlock address has not been set, using default for state type: %s", P.GitlabHttpState.HttpUnlockAddress)
			}

			setup.C.EnvVars["TF_HTTP_UNLOCK_ADDRESS"] = P.GitlabHttpState.HttpUnlockAddress
			setup.C.EnvVars["TF_HTTP_UNLOCK_METHOD"] = P.GitlabHttpState.HttpUnlockMethod

			if P.GitlabHttpState.HttpUsername == "" {
				return fmt.Errorf("TF_HTTP_USERNAME is required for state type: %s", P.State.Type)
			}
			setup.C.EnvVars["TF_HTTP_USERNAME"] = P.GitlabHttpState.HttpUsername

			if P.GitlabHttpState.HttpPassword == "" {
				return fmt.Errorf("TF_HTTP_PASSWORD is required for state type: %s", P.State.Type)
			}
			setup.C.EnvVars["TF_HTTP_PASSWORD"] = P.GitlabHttpState.HttpPassword
			t.Plumber.AppendSecrets(P.GitlabHttpState.HttpPassword)

			setup.C.EnvVars["TF_HTTP_RETRY_WAIT_MIN"] = P.GitlabHttpState.HttpRetryWaitMin

			t.Log.Debugf("Generated following environment variables for terraform to consume: %+v", setup.C.EnvVars)

			return nil
		})
}
