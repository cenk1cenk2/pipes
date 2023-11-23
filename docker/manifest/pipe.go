package manifest

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
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
			return ProcessFlags(tl)
		}).
		Set(func(tl *TaskList[Pipe]) Job {
			return tl.JobSequence(
				tl.JobParallel(
					tl.JobSequence(
						DiscoverPublishedImageFiles(tl).Job(),
						FetchPublishedImagesFromFiles(tl).Job(),
					),
					FetchUserPublishedImages(tl).Job(),
				),
				UpdateManifests(tl).Job(),
			)
		})
}
