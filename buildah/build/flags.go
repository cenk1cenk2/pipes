package build

import (
	"strings"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
	"gitlab.kilic.dev/devops/pipes/internal/git"
	"gitlab.kilic.dev/devops/pipes/internal/tagsfile"

	. "github.com/cenk1cenk2/plumber/v6"
)

const (
	CategoryContainerImage    = "Container Image"
	CategoryContainerFile     = "Containerfile"
	CategoryContainerManifest = "Container Manifest"

	DefaultTagAsLatest  = `[ "^tags/v?\\d+.\\d+.\\d+$" ]`
	DefaultSanitizeTags = `[
    { "match": "([^/]*)/(.*)", "template": "{{ index $ 1 | upper }}_{{ index $ 2 }}" }
]`
)

//revive:disable:line-length-limit

var Flags = CombineFlags(
	git.NewFlags(git.Options{Destination: &P.Git}),
	tagsfile.NewFlags(tagsfile.Options{Destination: &P.Image.TagsFile, Strict: &P.Image.TagsFileStrict}),
	[]cli.Flag{

		// CategoryContainerImage

		&cli.StringSliceFlag{
			Category: CategoryContainerImage,
			Name:     "buildah.build.image.platforms",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_IMAGE_PLATFORMS"),
				cli.EnvVar("CONTAINER_IMAGE_PLATFORMS"),
			),
			Usage:       "Container image platforms to be built.",
			Required:    false,
			Value:       []string{},
			Destination: &P.Image.Platforms,
		},

		&cli.StringFlag{
			Category: CategoryContainerImage,
			Name:     "buildah.build.image.name",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_IMAGE_NAME"),
				cli.EnvVar("CONTAINER_IMAGE_NAME"),
			),
			Usage:       "Image name for the container image to be built.",
			Required:    true,
			Destination: &P.Image.Name,
		},

		&cli.StringSliceFlag{
			Category: CategoryContainerImage,
			Name:     "buildah.build.image.tags",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_IMAGE_TAGS"),
				cli.EnvVar("CONTAINER_IMAGE_TAGS"),
			),
			Usage:       "Image tags for the container image to be built.",
			Required:    true,
			Destination: &P.Image.Tags,
		},

		flags.YAMLFlag(&P.Image.TagsTemplate, &cli.StringFlag{
			Category: CategoryContainerImage,
			Name:     "buildah.build.image.tags-template",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_IMAGE_TAGS_TEMPLATE"),
				cli.EnvVar("CONTAINER_IMAGE_TAGS_TEMPLATE"),
			),
			Usage: strings.TrimSpace(`
Modifies every tag that matches a certain condition.
Template is interpolated with the given matches in the regular expression.

format(yaml([]struct{ match: RegExp, template: Template(match) }))
`),
			Required: false,
			Value:    "[]",
		}),

		flags.YAMLFlag(&P.Image.TagsSanitize, &cli.StringFlag{
			Category: CategoryContainerImage,
			Name:     "buildah.build.image.tags-sanitize",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_IMAGE_TAGS_SANITIZE"),
				cli.EnvVar("CONTAINER_IMAGE_SANITIZE_TAGS"),
			),
			Usage: strings.TrimSpace(`
Sanitizes the given regex pattern out of tag name.
Template is interpolated with the given matches in the regular expression.

format(yaml([]struct{ match: RegExp, template: Template(match) }))
`),
			Required: false,
			Value:    DefaultSanitizeTags,
		}),

		flags.YAMLFlag(&P.Image.TagAsLatest, &cli.StringFlag{
			Category: CategoryContainerImage,
			Name:     "buildah.build.image.tag-as-latest",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_IMAGE_TAG_AS_LATEST"),
				cli.EnvVar("CONTAINER_IMAGE_TAGS_AS_LATEST"),
			),
			Usage: strings.TrimSpace(`
Regex pattern to tag the image as latest.
Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags.

format(yaml([]RegExp))
`),
			Required: false,
			Value:    DefaultTagAsLatest,
		}),

		&cli.BoolFlag{
			Category: CategoryContainerImage,
			Name:     "buildah.build.image.pull",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_IMAGE_PULL"),
				cli.EnvVar("CONTAINER_IMAGE_PULL"),
			),
			Usage:       "Pull before building the image.",
			Required:    false,
			Value:       true,
			Destination: &P.Image.Pull,
		},

		&cli.BoolFlag{
			Category: CategoryContainerImage,
			Name:     "buildah.build.image.push",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_IMAGE_PUSH"),
				cli.EnvVar("CONTAINER_IMAGE_PUSH"),
			),
			Usage:       "Push the image after building.",
			Required:    false,
			Value:       true,
			Destination: &P.Image.Push,
		},

		flags.YAMLFlag(&P.Image.BuildArgs, &cli.StringFlag{
			Category: CategoryContainerImage,
			Name:     "buildah.build.image.build-args",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_IMAGE_BUILD_ARGS"),
				cli.EnvVar("CONTAINER_IMAGE_BUILD_ARGS"),
			),
			Usage: strings.TrimSpace(`
Pass in extra build arguments for image.
You can use it as a template with environment variables as the context.

format(yaml(map[string]Template()))
`),
			Required: false,
			Value:    "",
		}),

		&cli.StringFlag{
			Category: CategoryContainerImage,
			Name:     "buildah.build.image.latest-tag",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_IMAGE_LATEST_TAG"),
				cli.EnvVar("CONTAINER_IMAGE_LATEST_TAG"),
			),
			Usage:       "Latest tag for the container image where it is marked as latest.",
			Required:    false,
			Value:       "latest",
			Destination: &P.Image.LatestTag,
		},

		&cli.StringFlag{
			Category: CategoryContainerImage,
			Name:     "buildah.build.image.cache",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_IMAGE_CACHE"),
				cli.EnvVar("CONTAINER_IMAGE_CACHE"),
			),
			Usage:       "Specify the cache for the container image.",
			Required:    false,
			Value:       "",
			Destination: &P.Image.Cache,
		},

		&cli.StringFlag{
			Category: CategoryContainerImage,
			Name:     "buildah.build.image.format",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_IMAGE_FORMAT"),
				cli.EnvVar("CONTAINER_IMAGE_FORMAT"),
			),
			Usage:       `Specify the format for Container Image. format(enum("oci", "docker"))`,
			Required:    false,
			Value:       "oci",
			Destination: &P.Image.Format,
		},

		&cli.StringFlag{
			Category: CategoryContainerImage,
			Name:     "buildah.build.image.storage-driver",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_IMAGE_STORAGE_DRIVER"),
				cli.EnvVar("CONTAINER_IMAGE_STORAGE_DRIVER"),
				cli.EnvVar("BUILDAH_STORAGE_DRIVER"),
			),
			Usage:       `Specify the storage driver for Buildah. format(enum("overlay", "overlay2", "vfs"))`,
			Required:    false,
			Value:       "vfs",
			Destination: &P.Image.StorageDriver,
		},

		// CategoryContainerFile

		&cli.StringFlag{
			Category: CategoryContainerFile,
			Name:     "buildah.build.file.context",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_FILE_CONTEXT"),
				cli.EnvVar("CONTAINER_FILE_CONTEXT"),
			),
			Usage:       "Containerfile context argument for build operation.",
			Required:    false,
			Value:       ".",
			Destination: &P.File.Context,
		},

		&cli.StringFlag{
			Category: CategoryContainerFile,
			Name:     "buildah.build.file.name",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_FILE_NAME"),
				cli.EnvVar("CONTAINER_FILE_NAME"),
			),
			Usage:       "Containerfile path for the build operation.",
			Required:    false,
			Value:       "Dockerfile",
			Destination: &P.File.Name,
		},

		// CategoryContainerManifest

		&cli.StringFlag{
			Category: CategoryContainerManifest,
			Name:     "buildah.build.manifest.target",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_MANIFEST_TARGET"),
				cli.EnvVar("CONTAINER_MANIFEST_TARGET"),
			),
			Usage:       "Target image names for patching the manifest. format(Template([]string))",
			Required:    false,
			Destination: &P.Manifest.Target,
		},

		&cli.StringFlag{
			Category: CategoryContainerManifest,
			Name:     "buildah.build.manifest.file",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("BUILDAH_BUILD_MANIFEST_FILE"),
				cli.EnvVar("CONTAINER_MANIFEST_FILE"),
			),
			Usage:       "Write all the published images into a file for later use. format(Template([]string))",
			Value:       `.published-container-images_{{ $ | join "," | sha256sum }}`,
			Destination: &P.Manifest.File,
		},
	})
