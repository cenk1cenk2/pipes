package login

import (
	"github.com/urfave/cli/v2"

	. "gitlab.kilic.dev/libraries/plumber/v3"
)

type (
	Npm struct {
		Login     string
		NpmRcFile cli.StringSlice
		NpmRc     string
	}

	Pipe struct {
		Ctx

		Npm
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).SetTasks(
		TL.JobSequence(
			Setup(&TL).Job(),
			GenerateNpmRc(&TL).Job(),
			VerifyNpmLogin(&TL).Job(),
		),
	)
}
