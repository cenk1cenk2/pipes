package setup

import (
	"github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

var P = &tool.Config{}
var C = tool.NewCtx()

// Step is what every pulumi command starts with, since the tool has to be
// resolved before anything can be asked of it.
var Step = cli.Step{Flags: Flags, New: New}

func New(p *plumber.Plumber) *plumber.TaskList {
	return tool.Setup(p, Spec, P, C)
}
