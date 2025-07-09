package pipe

import (
	"github.com/joho/godotenv"
	"gitlab.kilic.dev/devops/pipes/select-env/setup"
	. "github.com/cenk1cenk2/plumber/v6"
)

func WriteEnvironmentFile(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("environment", "file").
		Set(func(t *Task[Pipe]) error {
			return godotenv.Write(setup.TL.Pipe.Ctx.EnvVars, t.Pipe.Environment.File)
		})
}
