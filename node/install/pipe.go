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
		NodeInstall NodeInstall
	}
)

var P = TaskList[Pipe, Ctx]{
	Pipe:    Pipe{},
	Context: Ctx{},
}

func New(a *App) *TaskList[Pipe, Ctx] {
	return P.New(a).SetTasks(
		P.JobSequence(
			InstallNodeDependencies(&P).Job(),
		),
	)
}
