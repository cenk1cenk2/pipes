package setup

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Environment struct {
		Conditions        []EnvironmentConditionJson
		FailOnNoReference bool
	}

	Git struct {
		Branch string
		Tag    string
	}

	Pipe struct {
		Ctx

		Environment
		Git
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).Set(func(tl *TaskList[Pipe]) Job {
		return tl.JobSequence(
			Setup(tl).Job(),

			SelectEnvironment(tl).Job(),
			FetchEnvironment(tl).Job(),
		)
	})
}
