package manifest

import (
	"fmt"

	. "github.com/cenk1cenk2/plumber/v6"
	icli "gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

type (
	ContainerManifest struct {
		Target string
		Images []string
		Files  []string
		Matrix []ContainerManifestMatrix
	}

	ContainerManifestMatrix struct {
		Target string   `json:"target,omitempty" yaml:"target,omitempty"`
		Images []string `json:"images"           yaml:"images"`
	}

	Pipe struct {
		ContainerManifest
	}

	Ctx struct {
		ManifestedImages map[string][]string
		Matches          []string
	}

	// Deps is the registry the login step authenticated against, whose uri the
	// target manifest is prefixed with so it is created under the name the images
	// beneath it were pushed to.
	Deps struct {
		Registry *registry.Credentials
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber, deps Deps) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if deps.Registry.Uri != "" {
				P.ContainerManifest.Target = fmt.Sprintf("%s/%s", deps.Registry.Uri, P.ContainerManifest.Target)

				tl.Log.Infof("Using default manifest target: %s", P.ContainerManifest.Target)
			}

			if err := icli.Validated(p, P); err != nil {
				return err
			}

			C.ManifestedImages = make(map[string][]string)

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				JobParallel(
					JobSequence(
						DiscoverPublishedImageFiles(tl).Job(),
						FetchPublishedImagesFromFiles(tl).Job(),
					),
					FetchUserPublishedImages(tl).Job(),
				),
				UpdateManifests(tl).Job(),
			)
		})
}

func Step(deps Deps) icli.Step {
	return icli.Step{
		Flags: Flags,
		New:   func(p *Plumber) *TaskList { return New(p, deps) },
	}
}
