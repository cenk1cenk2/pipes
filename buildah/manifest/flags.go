package manifest

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"go.yaml.in/yaml/v4"
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
			cli.EnvVar("CONTAINER_MANIFEST_FILES"),
			cli.EnvVar("BUILDAH_MANIFEST_FILES"),
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
			cli.EnvVar("CONTAINER_MANIFEST_TARGET"),
			cli.EnvVar("BUILDAH_MANIFEST_TARGET"),
		),
		Usage:       "Target image names for patching the manifest. format(Template())",
		Required:    false,
		Destination: &P.Manifest.Target,
	},

	&cli.StringSliceFlag{
		Category: CATEGORY_CONTAINER_MANIFEST,
		Name:     "buildah.manifest.images",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_MANIFEST_IMAGES"),
			cli.EnvVar("BUILDAH_MANIFEST_IMAGES"),
		),
		Usage:       "Image names for patching the manifest with the given target.",
		Required:    false,
		Destination: &P.Manifest.Images,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_MANIFEST,
		Name:     "buildah.manifest.matrix",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_MANIFEST_MATRIX"),
			cli.EnvVar("BUILDAH_MANIFEST_MATRIX"),
		),
		Usage:            "Matrix of all the images that should be manifested. format(yaml([]struct { target: string, images: []string }))",
		Required:         false,
		ValidateDefaults: true,
		Validator: func(v string) error {
			if v == "" {
				return nil
			}

			if err := yaml.Unmarshal([]byte(v), &P.Manifest.Matrix); err != nil {
				return fmt.Errorf("Cannot unmarshal container manifest matrix: %w", err)
			}

			return nil
		},
	},
}
