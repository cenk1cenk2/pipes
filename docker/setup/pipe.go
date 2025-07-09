package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
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
		SetRuntimeDepth(3).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobParallel(
				DockerVersion(tl).Job(),
				DockerBuildXVersion(tl).Job(),
			)
		})
}
