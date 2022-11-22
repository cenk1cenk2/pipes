package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

type (
	NodeAdd struct {
		Packages   []string
		Global     bool
		ScriptArgs string
		Cwd        string
	}

	Pipe struct {
		NodeAdd
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		ShouldDisable(func(tl *TaskList[Pipe]) bool {
			return len(tl.Pipe.NodeAdd.Packages) == 0
		}).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				AddNodeModules(tl).Job(),
			)
		})
}
