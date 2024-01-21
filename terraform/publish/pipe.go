package publish

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

type (
	Module struct {
		TagsFile string
		Name     string
		Cwd      string
		System   string
	}

	Registry struct {
		Name   string
		Gitlab struct {
			ApiUrl    string
			ProjectId string
			Token     string
		}
	}

	Pipe struct {
		Ctx

		Module
		Registry
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
				TerraformTagsFile(tl).Job(),
				TerraformPackage(tl).Job(),
				TerraformPublish(tl).Job(),
			)
		})
}
