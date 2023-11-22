package pipe

import (
	"github.com/joho/godotenv"
	"gitlab.kilic.dev/devops/pipes/select-env/setup"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func WriteEnvironmentFile(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("environment", "file").
		Set(func(t *Task[Pipe]) error {
			return godotenv.Write(setup.TL.Pipe.Ctx.EnvVars, t.Pipe.Environment.File)
		})
}
