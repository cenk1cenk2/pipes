package install

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

var TL = TaskList{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				HelmInstall(tl).Job(),
			)
		})
}
