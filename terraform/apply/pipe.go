package apply

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type (
	Apply struct {
		Args   string
		Output string
	}

	Pipe struct {
		Apply
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				TerraformApply(tl).Job(),
			)
		})
}
