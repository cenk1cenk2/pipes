package build

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/flags"
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

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				BuildNodeApplication(tl).Job(),
			)
		})
}
