package manifest

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

//revive:disable:line-length-limit

const (
	CategoryContainerManifest = "Container Manifest"
)

// PublishedOptions is what the flag of the published prefix is built onto, so the
// build that writes under the prefix and the manifest that reads it agree on one.
type PublishedOptions struct {
	Destination *string
}

func NewPublishedFlags(opts PublishedOptions) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Category:    CategoryContainerManifest,
			Name:        "container-manifest.variable-prefix",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("CONTAINER_MANIFEST_VARIABLE_PREFIX")),
			Usage:       "Prefix of the dotenv variables the build jobs hand their published images over to the manifest job in.",
			Required:    false,
			Value:       "CONTAINER_MANIFEST_PUBLISHED_",
			Destination: opts.Destination,
		},
	}
}

var Flags = CombineFlags(
	NewPublishedFlags(PublishedOptions{Destination: &P.Manifest.PublishedPrefix}),
	[]cli.Flag{

		// CategoryContainerManifest

		&cli.StringFlag{
			Category: CategoryContainerManifest,
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
			Category: CategoryContainerManifest,
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
			Category: CategoryContainerManifest,
			Name:     "buildah.manifest.matrix",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_MANIFEST_MATRIX"),
				cli.EnvVar("CONTAINER_MANIFEST_MATRIX"),
			),
			Usage:    "Matrix of all the images that should be manifested. format(yaml([]struct{ target: string, images: []string }))",
			Required: false,
		}),
	})
