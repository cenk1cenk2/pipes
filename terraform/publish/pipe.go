package publish

import (
	"regexp"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
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
		Registry gitlab.ModuleRegistry
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

			C.Registry = gitlab.NewModuleRegistry(
				P.Registry.Gitlab.ApiUrl,
				P.Registry.Gitlab.ProjectId,
				P.Registry.Gitlab.Token,
			)

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
