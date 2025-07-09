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

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				GenerateTerraformEnvVarsState(tl).Job(),
			)
		})
}
