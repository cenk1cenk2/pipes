package login

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

type Pipe struct {
	Container []registry.LoginEntry
}

var P = &Pipe{}

// New is the login stage of every buildah command, verifying the ambient login
// of the image only for a registry the container login does not cover.
func New(p *Plumber, image string) *TaskList {
	tl := &TaskList{}

	container := registry.NewLogin(registry.LoginOptions{
		Command: []string{"buildah", "login"},
		Entries: P.Container,
	})

	return tl.New(p).
		SetRuntimeDepth(3).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				container.Job(tl),
				loginVerify(tl, container, image).Job(),
			)
		})
}
