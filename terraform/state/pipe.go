package state

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type (
	State struct {
		Type string `validate:"oneof=gitlab-http"`
		Name string
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
