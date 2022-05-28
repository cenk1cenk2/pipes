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
		Npm Npm
	}
)

var P = TaskList[Pipe, Ctx]{
	Pipe:    Pipe{},
	Context: Ctx{},
}

func New(a *App) *TaskList[Pipe, Ctx] {
	return P.New(a).SetTasks(
		P.JobSequence(
			Unmarshal(&P).Job(),
			GenerateNpmRc(&P).Job(),
			VerifyNpmLogin(&P).Job(),
		),
	)
}
