package publish

import (
	"regexp"

	. "github.com/cenk1cenk2/plumber/v6"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
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

	// Deps dials the registry only once the flags have been parsed, so the pipe
	// carries the way to reach one rather than a connection to it.
	Deps struct {
		Registry func(apiUrl, projectId, token string) gitlab.ModuleRegistry
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if P.Module.Name != "" {
				P.Module.Name = regexp.MustCompile(`[_ ]`).ReplaceAllString(P.Module.Name, "-")
			}

			if err := icli.Validated(p, P); err != nil {
				return err
			}

			C.Registry = deps.Registry(
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

func Step(deps Deps) icli.Step {
	return icli.Step{
		Flags: Flags,
		New:   func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
