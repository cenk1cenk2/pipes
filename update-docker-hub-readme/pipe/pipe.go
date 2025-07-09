package pipe

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	DockerHub struct {
		Username string
		Password string
		Address  string
	}

	Readme struct {
		Repository  string
		File        string
		Description string
		Matrix      []ReadmeMatrixJson
	}

	Pipe struct {
		Ctx

		DockerHub
		Readme
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				LoginToDockerHubRegistry(tl).Job(),
				DiscoverJobs(tl).Job(),
				UpdateDockerReadme(tl).Job(),
			)
		})
}
