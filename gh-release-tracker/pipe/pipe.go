package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

type (
	DockerImage struct {
		TagsFile string
	}

	Github struct {
		Token      string
		Repository string
	}

	Pipe struct {
		Ctx

		Github
		DockerImage
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).SetTasks(
		TL.JobSequence(
			Setup(&TL).Job(),
			GithubLogin(&TL).Job(),
			FetchLatestTag(&TL).Job(),
			WriteTagsFile(&TL).Job(),
		),
	)
}
