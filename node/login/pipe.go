package login

import (
	"github.com/urfave/cli/v2"

	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Npm struct {
		Login     []NpmLoginJson
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
	return TL.New(p).Set(func(tl *TaskList[Pipe]) Job {
		return tl.JobSequence(
			ProcessFlags(tl),

			GenerateNpmRc(tl).Job(),
			VerifyNpmLogin(tl).Job(),
		)
	})
}
