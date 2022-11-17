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
		Npm
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		SetName("node", "login").
		ShouldRunBefore(
			func(tl *TaskList[Pipe]) error {
				return ProcessFlags(tl)
			}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				GenerateNpmRc(tl).Job(),
				VerifyNpmLogin(tl).Job(),
			)
		})
}
