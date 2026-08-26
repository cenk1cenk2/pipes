package state

import (
	. "github.com/cenk1cenk2/plumber/v6"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
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

	// Deps is the environment the state configuration is written into, plus the CI
	// coordinates the default GitLab state address is built out of when none was
	// configured.
	Deps struct {
		Tool *tool.Ctx
		CI   *setup.CiVariables
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			return icli.Validated(p, P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				GenerateTerraformEnvVarsState(tl, deps).Job(),
			)
		})
}

func Step(deps Deps) icli.Step {
	return icli.Step{
		Flags: Flags,
		New:   func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
