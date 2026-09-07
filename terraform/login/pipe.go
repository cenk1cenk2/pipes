package login

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

type (
	Registry struct {
		Credentials []TerraformRegistryCredentialsJson
	}

	Pipe struct {
		Registry
	}

	// Deps is where the registry tokens land: the environment the terraform
	// commands of the surrounding pipe are run with.
	Deps struct {
		Tool *tool.Ctx
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				GenerateTerraformRegistryCredentialsEnvVars(tl, deps).Job(),
			)
		})
}
