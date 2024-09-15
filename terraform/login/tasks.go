package login

import (
	"strings"

	"gitlab.kilic.dev/devops/pipes/terraform/setup"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func GenerateTerraformRegistryCredentilsEnvVars(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("environment", "credentials").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(t.Pipe.Registry.Credentials) == 0
		}).
		Set(func(t *Task[Pipe]) error {
			for _, c := range t.Pipe.Registry.Credentials {
				t.Log.Infof("Generating registry token: %s", c.Registry)

				sanitized := strings.ReplaceAll(c.Registry, ".", "_")
				sanitized = strings.ReplaceAll(sanitized, "-", "__")

				setup.TL.Pipe.EnvVars["TF_TOKEN_"+sanitized] = c.Token
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
