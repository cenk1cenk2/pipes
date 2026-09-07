package write

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/select-env/setup"
)

func WriteEnvironmentFile(tl *TaskList) *Task {
	return tl.CreateTask("environment", "file").
		Set(func(t *Task) error {
			return environment.WriteFile(P.Environment.File, setup.EnvironmentCtx.EnvVars)
		})
}
