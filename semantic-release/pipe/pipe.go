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

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).SetTasks(
		TL.JobSequence(
			TL.JobParallel(
				InstallApkPackages(&TL).Job(),
				InstallNodePackages(&TL).Job(),
			),
			RunSemanticRelease(&TL).Job(),
		),
	)
}
