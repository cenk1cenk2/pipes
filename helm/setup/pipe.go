package setup

import (
	"github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
	helmv2 "helm.sh/helm/v4/pkg/chart/v2"
)

// Ctx carries the chart alongside what every tool pipe resolves, so the pipes
// that publish it read one context instead of two.
type Ctx struct {
	*tool.Ctx
	Chart *helmv2.Chart
}

var P = &tool.Config{}
var C = &Ctx{Ctx: tool.NewCtx()}

// Step is what every helm command starts with, since the tool has to be resolved
// and the chart read before anything can be asked of either.
var Step = cli.Step{Flags: Flags, New: New}

func New(p *plumber.Plumber) *plumber.TaskList {
	return tool.Setup(p, Spec, P, C.Ctx, HelmLoadChart)
}
