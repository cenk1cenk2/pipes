package publish

import (
	"regexp"

	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Module struct {
		TagsFile string
		Name     string
		Cwd      string
		System   string
	}

	Registry struct {
		Name   string `validate:"oneof=gitlab"`
		Gitlab struct {
			ApiUrl    string
			ProjectId string
			Token     string
		}
	}

	Pipe struct {
		Module
		Registry
	}

	Ctx struct {
		Tags     []string
		Packages []PublishablePackage
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if P.Module.Name != "" {
				P.Module.Name = regexp.MustCompile(`[_ ]`).ReplaceAllString(P.Module.Name, "-")
			}

			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				TerraformTagsFile(tl).Job(),
				TerraformPackage(tl).Job(),
				TerraformPublish(tl).Job(),
			)
		})
}
