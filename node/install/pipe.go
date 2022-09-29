package install

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	NodeInstall struct {
		Cwd         string `validate:"dir"`
		UseLockFile bool
		Args        string
	}

	Pipe struct {
		NodeInstall
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).SetTasks(
		TL.JobSequence(
			InstallNodeDependencies(&TL).Job(),
		),
	)
}
