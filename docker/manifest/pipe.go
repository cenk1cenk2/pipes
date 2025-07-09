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
		Ctx

		DockerManifest
	}
)

var TL = TaskList[Pipe]{
	Pipe: Pipe{},
}

func New(p *Plumber) *TaskList[Pipe] {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList[Pipe]) error {
			if login.TL.Pipe.DockerRegistry.Registry != "" {
				tl.Pipe.DockerManifest.Target = fmt.Sprintf("%s/%s", login.TL.Pipe.DockerRegistry.Registry, tl.Pipe.DockerManifest.Target)
			}

			tl.Pipe.ManifestedImages = make(map[string][]string)

			return nil
		}).
		Set(func(tl *TaskList[Pipe]) Job {
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
