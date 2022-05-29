package install

import (
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type (
	NodeInstall struct {
		Cwd         string `validate:"dir"`
		UseLockFile bool
	}

	Pipe struct {
		NodeInstall
		Ctx
	}
)

var P = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(a *App) *TaskList[Pipe] {
	return P.New(a).SetTasks(
		P.JobSequence(
			InstallNodeDependencies(&P).Job(),
		),
	)
}
