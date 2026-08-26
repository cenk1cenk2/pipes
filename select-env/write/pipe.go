package write

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
)

type (
	Environment struct {
		File string
	}

	Pipe struct {
		Environment
	}
)

var TL = TaskList{}

var P = &Pipe{}

// The environment selection lives here rather than in main so that the flags,
// the environment task list and the file this pipe writes all read the same
// instance.
var (
	EnvironmentConfig = &environment.Config{}
	EnvironmentCtx    = &environment.Ctx{}
)

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				WriteEnvironmentFile(tl).Job(),
			)
		})
}
