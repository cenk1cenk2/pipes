package setup

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	helmv2 "helm.sh/helm/v4/pkg/chart/v2"
)

type (
	Pipe struct {
		Cwd string `validate:"omitempty,dir"`
	}

	// Ctx carries the chart alongside what the setup resolves, so the pipes that publish it read one context.
	Ctx struct {
		Cwd     string
		Version string
		Env     map[string]string
		Chart   *helmv2.Chart
	}
)

var P = &Pipe{}
var C = &Ctx{Env: map[string]string{}}

func New(p *Plumber) *TaskList {
	tl := &TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, _ *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				initialize(tl).Job(),
				version(tl).Job(),
				read(tl).Job(),
			)
		})
}
