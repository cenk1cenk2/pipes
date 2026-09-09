package manifest

import (
	"context"
	"fmt"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/buildah/login"
)

type (
	Manifest struct {
		Target string
		Images []string
		Files  []string
		Matrix []ManifestMatrix
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
		Matches          []string
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, tl *TaskList) error {
			if login.P.Uri != "" {
				P.Manifest.Target = fmt.Sprintf("%s/%s", login.P.Uri, P.Manifest.Target)

				tl.Log.Info(fmt.Sprintf("Using default manifest target: %s", P.Manifest.Target))
			}

			if err := p.Validate(P); err != nil {
				return err
			}

			C.ManifestedImages = make(map[string][]string)

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				JobParallel(
					JobSequence(
						discoverFile(tl).Job(),
						fetchFile(tl).Job(),
					),
					fetchUser(tl).Job(),
				),
				manifest(tl).Job(),
			)
		})
}
