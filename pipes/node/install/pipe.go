package install

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	NodeInstall struct {
		Cwd         string `validate:"dir"`
		UseLockFile bool
		Args        string
		Cache       bool
	}

	Pipe struct {
		NodeInstall
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
				InstallNodeDependencies(tl).Job(),
			)
		})
}
