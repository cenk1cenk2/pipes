package pipe

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v4"
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
		Ctx

		Packages
		SemanticRelease
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		SetName("semantic-release").
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				tl.JobParallel(
					InstallApkPackages(tl).Job(),
					InstallNodePackages(tl).Job(),
				),
				RunSemanticRelease(tl).Job(),
			)
		})
}
