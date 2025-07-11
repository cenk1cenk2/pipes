package login

import (
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
)

func GenerateTerraformRegistryCredentialsEnvVars(tl *TaskList) *Task {
	return tl.CreateTask("environment", "credentials").
		ShouldDisable(func(t *Task) bool {
			return len(P.Registry.Credentials) == 0
		}).
		Set(func(t *Task) error {
			for _, c := range P.Registry.Credentials {
				t.Log.Infof("Generating registry token: %s", c.Registry)

				sanitized := strings.ReplaceAll(c.Registry, ".", "_")
				sanitized = strings.ReplaceAll(sanitized, "-", "__")

				setup.C.EnvVars["TF_TOKEN_"+sanitized] = c.Token
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
