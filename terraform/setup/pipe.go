package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Project struct {
		Cwd        string   `validate:"omitempty,dir"`
		Workspaces []string `validate:"omitempty"`
	}

	CiVariables struct {
		JobId            string
		CommitSha        string
		JobStage         string
		ProjectId        string
		ProjectName      string
		ProjectNamespace string
		ProjectPath      string
		ProjectUrl       string
		ApiUrl           string
	}

	Config struct {
		LogLevel string `validate:"omitempty,oneof=trace debug info warn error"`
	}

	Pipe struct {
		Project
		Config
		CiVariables
	}

	Ctx struct {
		EnvVars map[string]string
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if err := p.Validate(P); err != nil {
				return err
			}

			C.EnvVars = make(map[string]string)

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobParallel(
				Version(tl).Job(),
				DiscoverWorkspaces(tl).Job(),
				GenerateTerraformEnvVars(tl).Job(),
			)
		})
}
