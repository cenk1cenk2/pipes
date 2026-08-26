package build

import (
	. "github.com/cenk1cenk2/plumber/v6"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/kustomize/setup"
)

type (
	Pipe struct {
		EnableHelm     bool
		HelmCommand    string
		LoadRestrictor string `validate:"omitempty,oneof=rootOnly none"`
		KubeVersion    string
		Output         string `validate:"omitempty,dirpath"`
	}

	OverlayResult struct {
		Overlay  string
		Yaml     []byte
		DocCount int
		Err      error
	}

	Ctx struct {
		Results []OverlayResult
	}

	// Deps is what the setup step resolved: the overlays to render, which is the
	// whole of the work this step has to do.
	Deps struct {
		Tool *setup.Ctx
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			return icli.Validated(p, P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				RenderOverlays(tl, deps).Job(),
			)
		})
}

func Step(deps Deps) icli.Step {
	return icli.Step{
		Flags: Flags,
		New:   func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
