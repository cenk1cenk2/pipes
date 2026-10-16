package manifest

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
)

type (
	Manifest struct {
		Target          string
		Images          []string
		Matrix          []ManifestMatrix
		PublishedPrefix string
	}

	ManifestMatrix struct {
		Target string   `json:"target,omitempty" yaml:"target,omitempty"`
		Images []string `json:"images"           yaml:"images"`
	}

	Pipe struct {
		Manifest
	}

	Ctx struct {
		ManifestedImages map[string][]string
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			if err := p.Validate(P); err != nil {
				return err
			}

			C.ManifestedImages = make(map[string][]string)

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				JobParallel(
					fetchEnv(tl).Job(),
					fetchUser(tl).Job(),
				),
				manifest(tl).Job(),
			)
		})
}
