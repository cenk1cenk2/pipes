package setup

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type (
	Project struct {
		Cwd        string   `validate:"omitempty,dir"`
		Workspaces []string `validate:"omitempty,dir,dive"`
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
		Ctx

		Project
		Config
		CiVariables
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
			return tl.JobParallel(
				Version(tl).Job(),
				DiscoverWorkspaces(tl).Job(),
				GenerateTerraformEnvVars(tl).Job(),
			)
		})
}
