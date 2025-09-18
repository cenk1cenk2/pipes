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
		Name:     "container-manifest.files",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_MANIFEST_FILES"),
		),
		Usage:       "Read published tags from a file. format(glob)",
		Required:    false,
		Value:       []string{"**/.published-container-images*"},
		Destination: &P.ContainerManifest.Files,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_MANIFEST,
		Name:     "container-manifest.target",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_MANIFEST_TARGET"),
		),
		Usage:       "Target image names for patching the manifest. format(Template())",
		Required:    false,
		Destination: &P.ContainerManifest.Target,
	},

	&cli.StringSliceFlag{
		Category: CATEGORY_CONTAINER_MANIFEST,
		Name:     "container-manifest.images",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_MANIFEST_IMAGES"),
		),
		Usage:       "Image names for patching the manifest with the given target.",
		Required:    false,
		Destination: &P.ContainerManifest.Images,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_MANIFEST,
		Name:     "container-manifest.matrix",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_MANIFEST_MATRIX"),
		),
		Usage:            "Matrix of all the images that should be manifested. format(yaml([]struct { target: string, images: []string }))",
		Required:         false,
		ValidateDefaults: true,
		Validator: func(v string) error {
			if v == "" {
				return nil
			}

			if err := yaml.Unmarshal([]byte(v), &P.ContainerManifest.Matrix); err != nil {
				return fmt.Errorf("Cannot unmarshal container manifest matrix: %w", err)
			}

			return nil
		},
	},
}
