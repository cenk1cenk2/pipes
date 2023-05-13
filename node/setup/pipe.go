package setup

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	Node struct {
		PackageManager string `validate:"oneof=npm yarn pnpm"`
	}

	Pipe struct {
		Ctx

		Node
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				SetupPackageManager(tl).Job(),
				NodeVersion(tl).Job(),
			)
		})
}
