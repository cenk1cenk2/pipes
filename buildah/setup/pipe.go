package setup

import (
	"github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

var P = &tool.Config{}
var C = tool.NewCtx()

func New(p *plumber.Plumber) *plumber.TaskList {
	return tool.Setup(p, Spec, P, C)
}
