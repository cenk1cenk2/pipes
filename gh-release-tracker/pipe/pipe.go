package pipe

import (
	. "github.com/cenk1cenk2/plumber/v6"
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
	return TL.New(p).
		SetRuntimeDepth(3).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				Setup(tl).Job(),
				GithubLogin(tl).Job(),
				FetchLatestTag(tl).Job(),
				WriteTagsFile(tl).Job(),
			)
		})
}
