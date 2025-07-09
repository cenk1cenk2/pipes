package manifest

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{

	// CATEGORY_DOCKER_MANIFEST

	&cli.StringSliceFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker_manifest.files",
		Usage:    "Read published tags from a file. format(glob)",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_MANIFEST_FILES"),
		),
		Value:       []string{"**/.published-docker-images*"},
		Destination: &TL.Pipe.DockerManifest.Files,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker_manifest.target",
		Usage:    "Target image names for patching the manifest. format(Template[string]())",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_MANIFEST_TARGET"),
		),
		Destination: &TL.Pipe.DockerManifest.Target,
	},

	&cli.StringSliceFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker_manifest.images",
		Usage:    "Image names for patching the manifest with the given target.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_MANIFEST_IMAGES"),
		),
		Destination: &TL.Pipe.DockerManifest.Images,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker_manifest.matrix",
		Usage:    "Matrix of all the images that should be manifested. json([]struct { target: string, images: []string })",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_MANIFEST_MATRIX"),
		),
		Action: func(_ context.Context, command *cli.Command, v string) error {
			if err := json.Unmarshal([]byte(v), &TL.Pipe.DockerManifest.Matrix); err != nil {
				return fmt.Errorf("Can not unmarshal Docker manifest matrix: %w", err)

			}
			return nil
		},
	},
}
