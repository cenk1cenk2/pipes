package write

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/joho/godotenv"
	"gitlab.kilic.dev/devops/pipes/select-env/setup"
)

func WriteEnvironmentFile(tl *TaskList) *Task {
	return tl.CreateTask("environment", "file").
		Set(func(t *Task) error {
			return godotenv.Write(setup.EnvironmentCtx.EnvVars, P.Environment.File)
		})
}
