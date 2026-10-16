package login

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

type Pipe struct {
	Container []registry.LoginEntry
}

var P = &Pipe{}

// New is the login stage of every helm command that pulls a dependency or pushes a chart.
func New(p *Plumber) *TaskList {
	tl := &TaskList{}

	container := registry.NewLogin(registry.LoginOptions{
		Command: []string{"helm", "registry", "login"},
		Entries: P.Container,
	})

	return tl.New(p).
		SetRuntimeDepth(3).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				container.Job(tl),
			)
		})
}
