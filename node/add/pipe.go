package pipe

import (
	. "github.com/cenk1cenk2/plumber/v6"
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

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldDisable(func(tl *TaskList) bool {
			return len(P.NodeAdd.Packages) == 0
		}).
		ShouldRunBefore(func(tl *TaskList) error {
			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				AddNodeModules(tl).Job(),
			)
		})
}
