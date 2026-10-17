package login

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

type Pipe struct {
	Container []registry.LoginEntry
}

var P = &Pipe{}

// New is the login stage of every helm command that pulls a dependency or pushes a chart,
// requiring the credentials of every entry for a command that pushes.
func New(p *Plumber, required bool) *TaskList {
	tl := &TaskList{}

	container := registry.NewLogin(registry.LoginOptions{
		Command:  []string{"helm", "registry", "login"},
		Entries:  P.Container,
		Required: required,
	})

	return tl.New(p).
		SetRuntimeDepth(3).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				container.Job(tl),
			)
		})
}
