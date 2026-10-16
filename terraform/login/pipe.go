package login

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
)

type (
	// Login is one registry the providers and modules are fetched from.
	Login struct {
		Uri   string `yaml:"uri"   validate:"required"`
		Token string `yaml:"token"`
	}

	Pipe struct {
		Logins []Login `validate:"dive"`
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
