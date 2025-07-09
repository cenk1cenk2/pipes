package manifest

import (
	"fmt"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/docker/login"
)

type (
	DockerManifest struct {
		Target string
		Images []string
		Files  []string
		Matrix []DockerManifestMatrixJson
	}

	Pipe struct {
		DockerManifest
	}
)

var TL = TaskList{}

var P = &Pipe{}
var C = &Ctx{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if err := p.Validator.Struct(P); err != nil {
				return err
			}

			if login.P.DockerRegistry.Registry != "" {
				P.DockerManifest.Target = fmt.Sprintf("%s/%s", login.P.DockerRegistry.Registry, P.DockerManifest.Target)
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
