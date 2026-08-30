package setup

import (
	"github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

type (
	Pipe struct {
		tool.Config
		Cache     string `validate:"omitempty,dirpath"`
		Workspace bool
	}

	// Ctx carries whether the modules are driven as a workspace alongside what
	// every tool pipe resolves, so the commands that have to tell the two apart
	// read one context instead of probing the toolchain again.
	Ctx struct {
		*tool.Ctx
		Workspace bool
	}
)

var P = &Pipe{}
var C = &Ctx{Ctx: tool.NewCtx()}

// Step is what every go command starts with, since the tool has to be resolved
// and the cache pointed at before anything can be asked of it.
var Step = cli.Step{Flags: Flags, New: New}

func New(p *plumber.Plumber) *plumber.TaskList {
	return tool.Setup(p, Spec, &P.Config, C.Ctx, GoEnv, GoWorkspace).
		ShouldRunBefore(func(_ *plumber.TaskList) error {
			return cli.Validated(p, P)
		})
}
