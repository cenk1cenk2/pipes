package manifest

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/docker/login"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{

	// CATEGORY_DOCKER_MANIFEST

	&cli.StringSliceFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker_manifest.files",
		Usage:    "Read published tags from a file. format(glob)",
		Required: false,
		EnvVars:  []string{"DOCKER_MANIFEST_FILES"},
		Value:    cli.NewStringSlice("**/.published-docker-images*"),
	},

	&cli.StringFlag{
		Category:    setup.CATEGORY_DOCKER_MANIFEST,
		Name:        "docker_manifest.target",
		Usage:       "Target image names for patching the manifest.",
		Required:    false,
		EnvVars:     []string{"DOCKER_MANIFEST_TARGET"},
		Destination: &TL.Pipe.DockerManifest.Target,
	},

	&cli.StringSliceFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker_manifest.images",
		Usage:    "Image names for patching the manifest with the given target.",
		Required: false,
		EnvVars:  []string{"DOCKER_MANIFEST_IMAGES"},
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker_manifest.matrix",
		Usage:    "Matrix of all the images that should be manifested. json([]struct{ target: string, images: []string })",
		Required: false,
		EnvVars:  []string{"DOCKER_MANIFEST_MATRIX"},
	},
}

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList[Pipe]) error {
	tl.Pipe.DockerManifest.Files = tl.CliContext.StringSlice("docker_manifest.files")
	tl.Pipe.DockerManifest.Images = tl.CliContext.StringSlice("docker_manifest.images")

	// extend image names with registry

	if login.TL.Pipe.DockerRegistry.Registry != "" {
		tl.Pipe.DockerManifest.Target = fmt.Sprintf("%s/%s", login.TL.Pipe.DockerRegistry.Registry, tl.Pipe.DockerManifest.Target)
	}

	if v := tl.CliContext.String("docker_manifest.matrix"); v != "" {
		if err := json.Unmarshal([]byte(v), &tl.Pipe.DockerManifest.Matrix); err != nil {
			return fmt.Errorf("Can not unmarshal Docker manifest matrix: %w", err)
		}
	}

	tl.Pipe.ManifestedImages = make(map[string][]string)

	return nil
}
