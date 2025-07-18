package manifest

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{

	// CATEGORY_DOCKER_MANIFEST

	&cli.StringSliceFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker-manifest.files",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_MANIFEST_FILES"),
		),
		Usage:       "Read published tags from a file. format(glob)",
		Required:    false,
		Value:       []string{"**/.published-docker-images*"},
		Destination: &P.DockerManifest.Files,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker-manifest.target",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_MANIFEST_TARGET"),
		),
		Usage:       "Target image names for patching the manifest. format(Template[string]())",
		Required:    false,
		Destination: &P.DockerManifest.Target,
	},

	&cli.StringSliceFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker-manifest.images",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_MANIFEST_IMAGES"),
		),
		Usage:       "Image names for patching the manifest with the given target.",
		Required:    false,
		Destination: &P.DockerManifest.Images,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker-manifest.matrix",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_MANIFEST_MATRIX"),
		),
		Usage:       "Matrix of all the images that should be manifested. json([]struct { target: string, images: []string })",
		Required:    false,
		Destination: &raw.DockerManifestMatrix,
	},
}
