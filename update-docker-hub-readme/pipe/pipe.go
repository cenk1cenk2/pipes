package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
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
	return TL.New(p).Set(func(tl *TaskList[Pipe]) Job {
		return tl.JobSequence(
			Setup(tl).Job(),
			LoginToDockerHubRegistry(tl).Job(),
			ReadReadmeFile(tl).Job(),
			UpdateDockerReadme(tl).Job(),
		)
	})
}
