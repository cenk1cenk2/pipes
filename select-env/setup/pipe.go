package setup

import (
	"gitlab.kilic.dev/devops/pipes/common/flags"
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Environment struct {
		Enable            bool
		Conditions        []EnvironmentConditionJson
		FailOnNoReference bool
		Strict            bool
	}

	Git flags.GitFlags

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
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldDisable(func(tl *TaskList[Pipe]) bool {
			return !tl.Pipe.Environment.Enable
		}).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				Setup(tl).Job(),

				SelectEnvironment(tl).Job(),
				FetchEnvironment(tl).Job(),
			)
		})
}
