package state

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	State struct {
		Type   string `validate:"omitempty,oneof=gitlab-http"`
		Name   string
		Strict bool
	}

	Credentials struct {
		Username string
		Password string
	}

	GitlabHttpState struct {
		HttpAddress       string
		HttpLockAddress   string
		HttpLockMethod    string
		HttpUnlockAddress string
		HttpUnlockMethod  string
		HttpUsername      string
		HttpPassword      string
		HttpRetryWaitMin  string
	}

	Pipe struct {
		State
		Credentials
		GitlabHttpState
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				GenerateTerraformEnvVarsState(tl).Job(),
			)
		})
}
