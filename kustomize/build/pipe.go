package build

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
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
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			return p.Validate(P)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				build(tl).Job(),
			)
		})
}
