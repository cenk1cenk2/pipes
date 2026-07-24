package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Pipe struct {
		Cwd               string `validate:"omitempty,dirpath"`
		Paths             []string
		DiscoveryPattern  []string
		DiscoveryStrategy string `validate:"omitempty,oneof=roots all"`
	}

	Ctx struct {
		Overlays []string
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
				KustomizeVersion(tl).Job(),
				DiscoverOverlays(tl).Job(),
			)
		})
}
