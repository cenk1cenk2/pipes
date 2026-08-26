package write

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
)

func WriteEnvironmentFile(tl *TaskList) *Task {
	return tl.CreateTask("environment", "file").
		Set(func(t *Task) error {
			return environment.WriteFile(P.Environment.File, EnvironmentCtx.EnvVars)
		})
}
