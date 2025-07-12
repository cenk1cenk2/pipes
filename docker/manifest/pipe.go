package manifest

import (
	"encoding/json"
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
		ShouldRunBefore(func(tl *TaskList) error {
			if v := p.Cli.String("docker-manifest.matrix"); v != "" {
				if err := json.Unmarshal([]byte(v), &P.DockerManifest.Matrix); err != nil {
					return fmt.Errorf("Can not unmarshal Docker manifest matrix: %w", err)
				}
			}

			if err := p.Validate(P); err != nil {
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
