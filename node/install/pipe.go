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

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		SetRuntimeDepth(3).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				InstallNodeDependencies(tl).Job(),
			)
		})
}
