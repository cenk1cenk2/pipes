package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Tags struct {
		File string
	}

	Github struct {
		Token      string
		Repository string
	}

	Pipe struct {
		Ctx

		Github
		Tags
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).Set(func(tl *TaskList[Pipe]) Job {
		return tl.JobSequence(
			Setup(tl).Job(),
			GithubLogin(tl).Job(),
			FetchLatestTag(tl).Job(),
			WriteTagsFile(tl).Job(),
		)
	})
}
