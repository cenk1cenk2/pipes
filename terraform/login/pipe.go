package login

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
)

type (
	Registry struct {
		Credentials []TerraformRegistryCredentialsJson
	}

	Pipe struct {
		Registry
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				environmentCredentials(tl).Job(),
			)
		})
}
