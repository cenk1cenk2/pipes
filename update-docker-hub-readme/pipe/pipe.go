package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v3"
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
	return TL.New(p).SetTasks(
		TL.JobSequence(
			Setup(&TL).Job(),
			LoginToDockerHubRegistry(&TL).Job(),
			ReadReadmeFile(&TL).Job(),
			UpdateDockerReadme(&TL).Job(),
		),
	)
}
