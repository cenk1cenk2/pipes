package write

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/joho/godotenv"
	"gitlab.kilic.dev/devops/pipes/select-env/setup"
)

func environmentFile(tl *TaskList) *Task {
	return tl.CreateTask("environment", "file").
		Set(func(_ context.Context, t *Task) error {
			return godotenv.Write(setup.EnvironmentCtx.EnvVars, P.Environment.File)
		})
}
