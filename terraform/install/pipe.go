package install

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type (
	Install struct {
		UseLockfile bool
		Reconfigure bool
		Args        string
	}

	Pipe struct {
		Install
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
				TerraformInstall(tl).Job(),
			)
		})
}
