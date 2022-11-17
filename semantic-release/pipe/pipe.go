package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Packages struct {
		Apk  []string
		Node []string
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
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
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
