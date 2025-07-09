package build

import (
	"gitlab.kilic.dev/devops/pipes/common/flags"
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Git flags.GitFlags

	NodeBuild struct {
		Script     string
		ScriptArgs string
		Cwd        string `validate:"dir"`
	}

	Pipe struct {
		Git
		NodeBuild
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		SetRuntimeDepth(3).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				BuildNodeApplication(tl).Job(),
			)
		})
}
