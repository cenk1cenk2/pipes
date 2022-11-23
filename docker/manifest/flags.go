package manifest

import (
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

	&cli.StringSliceFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker_manifest.targets",
		Usage:    "Target image names for patching the manifest.",
		Required: true,
		EnvVars:  []string{"DOCKER_MANIFEST_TARGETS"},
	},

	&cli.StringSliceFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker_manifest.images",
		Usage:    "Image names for patching the manifest with the given target.",
		Required: false,
		EnvVars:  []string{"DOCKER_MANIFEST_IMAGES"},
	},
}

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList[Pipe]) error {
	tl.Pipe.DockerManifest.Files = tl.CliContext.StringSlice("docker_manifest.files")
	tl.Pipe.DockerManifest.Images = tl.CliContext.StringSlice("docker_manifest.images")

	// extend image names with registry
	images := tl.CliContext.StringSlice("docker_manifest.targets")

	if login.TL.Pipe.DockerRegistry.Registry != "" {
		for i, image := range images {
			images[i] = fmt.Sprintf("%s/%s", login.TL.Pipe.DockerRegistry.Registry, image)
		}
	}

	tl.Pipe.DockerManifest.Targets = images

	return nil
}
