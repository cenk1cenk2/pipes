package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
	helmv2 "helm.sh/helm/v4/pkg/chart/v2"
)

type (
	Pipe struct {
		Cwd string `validate:"omitempty,dirpath"`
	}

	Ctx struct {
		Chart *helmv2.Chart
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				HelmVersion(tl).Job(),
				HelmLoadChart(tl).Job(),
			)
		})
}
