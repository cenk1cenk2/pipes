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
		Ctx
	}
)

var P = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(a *App) *TaskList[Pipe] {
	return P.New(a).SetTasks(
		P.JobSequence(
			P.JobParallel(
				InstallApkPackages(&P).Job(),
				InstallNodePackages(&P).Job(),
			),
			RunSemanticRelease(&P).Job(),
		),
	)
}
