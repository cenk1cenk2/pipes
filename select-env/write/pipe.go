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

	// Deps is the environment selection whose variables land in the file.
	Deps struct {
		Environment *environment.Ctx
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				WriteEnvironmentFile(tl, deps).Job(),
			)
		})
}
