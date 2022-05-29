package login

import (
	"github.com/urfave/cli/v2"

	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type (
	Npm struct {
		Login     string
		NpmRcFile cli.StringSlice
		NpmRc     string
	}

	Pipe struct {
		Npm
		Ctx
	}
)

var P = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(a *App) *TaskList[Pipe] {
	return P.New(a).SetTasks(
		P.JobSequence(
			Decode(&P).Job(),
			GenerateNpmRc(&P).Job(),
			VerifyNpmLogin(&P).Job(),
		),
	)
}
