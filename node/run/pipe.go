package run

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Environment struct {
		Enable bool
	}

	NodeCommand struct {
		Script string
		Cwd    string `validate:"dir"`
	}

	Pipe struct {
		Ctx

		Environment
		NodeCommand
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		SetName("node", "run").
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				RunNodeScript(tl).Job(),
			)
		})
}
