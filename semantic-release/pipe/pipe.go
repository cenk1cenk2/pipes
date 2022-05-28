package pipe

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type (
	Packages struct {
		Apk  cli.StringSlice
		Node cli.StringSlice
	}

	SemanticRelease struct {
		IsDryRun bool
		UseMulti bool
	}

	Pipe struct {
		Packages
		SemanticRelease
	}
)

var P = TaskList[Pipe, Ctx]{
	Pipe:    Pipe{},
	Context: Ctx{},
}

func New(a *App) *TaskList[Pipe, Ctx] {
	return P.New(a).SetTasks(
		P.JobSequence(
			InstallPackages(&P).Job(),
			RunSemanticRelease(&P).Job(),
		),
	)
}
