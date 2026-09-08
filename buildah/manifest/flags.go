package manifest

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

//revive:disable:line-length-limit

const (
	CATEGORY_CONTAINER_MANIFEST = "Container Manifest"
)

var Flags = []cli.Flag{

	// CATEGORY_CONTAINER_MANIFEST

	&cli.StringSliceFlag{
		Category: CATEGORY_CONTAINER_MANIFEST,
		Name:     "buildah.manifest.files",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("BUILDAH_MANIFEST_FILES"),
			cli.EnvVar("CONTAINER_MANIFEST_FILES"),
		),
		Usage:       "Read published tags from a file. format(glob)",
		Required:    false,
		Value:       []string{"**/.published-container-images*"},
		Destination: &P.Manifest.Files,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_MANIFEST,
		Name:     "buildah.manifest.target",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("BUILDAH_MANIFEST_TARGET"),
			cli.EnvVar("CONTAINER_MANIFEST_TARGET"),
		),
		Usage:       "Target image names for patching the manifest. format(Template())",
		Required:    false,
		Destination: &P.Manifest.Target,
	},

	&cli.StringSliceFlag{
		Category: CATEGORY_CONTAINER_MANIFEST,
		Name:     "buildah.manifest.images",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("BUILDAH_MANIFEST_IMAGES"),
			cli.EnvVar("CONTAINER_MANIFEST_IMAGES"),
		),
		Usage:       "Image names for patching the manifest with the given target.",
		Required:    false,
		Destination: &P.Manifest.Images,
	},

	flags.YAMLFlag(&P.Manifest.Matrix, &cli.StringFlag{
		Category: CATEGORY_CONTAINER_MANIFEST,
		Name:     "buildah.manifest.matrix",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("BUILDAH_MANIFEST_MATRIX"),
			cli.EnvVar("CONTAINER_MANIFEST_MATRIX"),
		),
		Usage:    "Matrix of all the images that should be manifested. format(yaml([]struct{ target: string, images: []string }))",
		Required: false,
	}),
}
