package setup

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Docker struct {
		UseBuildKit    bool
		UseBuildx      bool
		BuildxInstance string
	}

	Pipe struct {
		Ctx

		Docker
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
			return tl.JobParallel(
				DockerVersion(tl).Job(),
				DockerBuildXVersion(tl).Job(),
			)
		})
}
