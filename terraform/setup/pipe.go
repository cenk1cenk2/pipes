package setup

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type (
	Project struct {
		Cwd string
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

	Pipe struct {
		Ctx

		Project
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
			return tl.JobSequence(
				Setup(tl).Job(),
				Version(tl).Job(),
				GenerateTerraformEnvVars(tl).Job(),
			)
		})
}
