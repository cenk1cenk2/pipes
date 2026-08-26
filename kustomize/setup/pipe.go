package setup

import (
	"github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

type (
	Pipe struct {
		tool.Config
		Paths []string
	}

	// Ctx carries the resolved overlays alongside what every tool pipe resolves, so
	// the build reads one context instead of two.
	Ctx struct {
		*tool.Ctx
		Overlays []string
	}
)

var P = &Pipe{}
var C = &Ctx{Ctx: tool.NewCtx()}

// Step is what every kustomize command starts with, since the tool has to be
// resolved and the overlays enumerated before anything can be built.
var Step = cli.Step{Flags: Flags, New: New}

func New(p *plumber.Plumber) *plumber.TaskList {
	return tool.Setup(p, Spec, &P.Config, C.Ctx, ResolveOverlays)
}
