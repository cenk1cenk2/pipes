package build

import (
	"gitlab.kilic.dev/devops/pipes/common/flags"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Git flags.GitFlags

	Environment struct {
		Enable bool
	}

	NodeBuild struct {
		Script     string
		ScriptArgs string
		Cwd        string `validate:"dir"`
	}

	Pipe struct {
		Environment
		Git
		NodeBuild
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				BuildNodeApplication(tl).Job(),
			)
		})
}
