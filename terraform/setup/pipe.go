package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	CiVariables struct {
		ProjectId string
		ApiUrl    string
	}

	Pipe struct {
		CiVariables
		Cwd      string `validate:"omitempty,dir"`
		LogLevel string `validate:"omitempty,oneof=trace debug info warn error"`
	}

	// Ctx is what New resolves; the values only land once the setup task list has run,
	// so the pipe holds on to the same instance the steps read back.
	Ctx struct {
		Cwd     string
		Version string
		Env     map[string]string
	}
)

var P = &Pipe{}
var C = &Ctx{Env: map[string]string{}}

func New(p *Plumber) *TaskList {
	tl := &TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				initialize(tl).Job(),
				version(tl).Job(),
				GenerateTerraformEnvVars(tl).Job(),
			)
		})
}
