package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Default struct {
		Flag string
	}

	Pipe struct {
		Ctx

		Default
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				DefaultTask(tl).Job(),
			)
		})
}
