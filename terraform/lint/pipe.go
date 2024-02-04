package lint

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type (
	Lint struct {
		FormatCheckEnable bool
		FormatCheckArgs   string
		ValidateEnable    bool
		ValidateArgs      string
	}

	Pipe struct {
		Lint
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
				TerraformLint(tl).Job(),
			)
		})
}
