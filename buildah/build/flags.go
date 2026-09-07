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
	CATEGORY_CONTAINER_IMAGE    = "Container Image"
	CATEGORY_CONTAINER_FILE     = "Containerfile"
	CATEGORY_CONTAINER_MANIFEST = "Container Manifest"

	DEFAULT_TAG_AS_LATEST = `[ "^tags/v?\\d+.\\d+.\\d+$" ]`
	DEFAULT_SANITIZE_TAGS = `[
    { "match": "([^/]*)/(.*)", "template": "{{ index $ 1 | upper }}_{{ index $ 2 }}" }
]`
)

//revive:disable:line-length-limit

var Flags = CombineFlags(
	git.NewFlags(&P.Git),
	tagsfile.NewFlags(&P.Image.TagsFile, "", &P.Image.TagsFileStrict, false),
	[]cli.Flag{

		// CATEGORY_CONTAINER_IMAGE

		&cli.StringSliceFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "buildah.build.image.platforms",
			Sources:     flags.EnvVars("BUILDAH_BUILD_IMAGE_PLATFORMS"),
			Usage:       "Container image platforms to be built.",
			Required:    false,
			Value:       []string{},
			Destination: &P.Image.Platforms,
		},

		&cli.StringFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "buildah.build.image.name",
			Sources:     flags.EnvVars("BUILDAH_BUILD_IMAGE_NAME"),
			Usage:       "Image name for the container image to be built.",
			Required:    true,
			Destination: &P.Image.Name,
		},

		&cli.StringSliceFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "buildah.build.image.tags",
			Sources:     flags.EnvVars("BUILDAH_BUILD_IMAGE_TAGS"),
			Usage:       "Image tags for the container image to be built.",
			Required:    true,
			Destination: &P.Image.Tags,
		},

		flags.YAMLFlag(&cli.StringFlag{
			Category: CATEGORY_CONTAINER_IMAGE,
			Name:     "buildah.build.image.tags-template",
			Sources:  flags.EnvVars("BUILDAH_BUILD_IMAGE_TAGS_TEMPLATE"),
			Usage: strings.TrimSpace(`
    Modifies every tag that matches a certain condition.
    Template is interpolated with the given matches in the regular expression.

    format(yaml([]struct{ match: RegExp, template: Template(match) }))
    `),
			Required: false,
			Value:    "[]",
		}, &P.Image.TagsTemplate),

		flags.YAMLFlag(&cli.StringFlag{
			Category: CATEGORY_CONTAINER_IMAGE,
			Name:     "buildah.build.image.tags-sanitize",
			Sources:  flags.EnvVars("BUILDAH_BUILD_IMAGE_TAGS_SANITIZE"),
			Usage: strings.TrimSpace(`
    Sanitizes the given regex pattern out of tag name.
    Template is interpolated with the given matches in the regular expression.

    format(yaml([]struct{ match: RegExp, template: Template(match) }))
    `),
			Required: false,
			Value:    DEFAULT_SANITIZE_TAGS,
		}, &P.Image.TagsSanitize),

		flags.YAMLFlag(&cli.StringFlag{
			Category: CATEGORY_CONTAINER_IMAGE,
			Name:     "buildah.build.image.tag-as-latest",
			Sources:  flags.EnvVars("BUILDAH_BUILD_IMAGE_TAG_AS_LATEST"),
			Usage: strings.TrimSpace(`
    Regex pattern to tag the image as latest.
    Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags.

    format(yaml([]RegExp))
    `),
			Required: false,
			Value:    DEFAULT_TAG_AS_LATEST,
		}, &P.Image.TagAsLatest),

		&cli.BoolFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "buildah.build.image.pull",
			Sources:     flags.EnvVars("BUILDAH_BUILD_IMAGE_PULL"),
			Usage:       "Pull before building the image.",
			Required:    false,
			Value:       true,
			Destination: &P.Image.Pull,
		},

		&cli.BoolFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "buildah.build.image.push",
			Sources:     flags.EnvVars("BUILDAH_BUILD_IMAGE_PUSH"),
			Usage:       "Push the image after building.",
			Required:    false,
			Value:       true,
			Destination: &P.Image.Push,
		},

		flags.YAMLFlag(&cli.StringFlag{
			Category: CATEGORY_CONTAINER_IMAGE,
			Name:     "buildah.build.image.build-args",
			Sources:  flags.EnvVars("BUILDAH_BUILD_IMAGE_BUILD_ARGS"),
			Usage: strings.TrimSpace(`
    Pass in extra build arguments for image.
    You can use it as a template with environment variables as the context.

    format(yaml(map[string]Template()))
    `),
			Required: false,
			Value:    "",
		}, &P.Image.BuildArgs),

		&cli.StringFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "buildah.build.image.latest-tag",
			Sources:     flags.EnvVars("BUILDAH_BUILD_IMAGE_LATEST_TAG"),
			Usage:       "Latest tag for the container image where it is marked as latest.",
			Required:    false,
			Value:       "latest",
			Destination: &P.Image.LatestTag,
		},

		&cli.StringFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "buildah.build.image.cache",
			Sources:     flags.EnvVars("BUILDAH_BUILD_IMAGE_CACHE"),
			Usage:       "Specify the cache for the container image.",
			Required:    false,
			Value:       "",
			Destination: &P.Image.Cache,
		},

		&cli.StringFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "buildah.build.image.format",
			Sources:     flags.EnvVars("BUILDAH_BUILD_IMAGE_FORMAT"),
			Usage:       "Specify the format for Container Image.",
			Required:    false,
			Value:       "oci",
			Destination: &P.Image.Format,
		},

		&cli.StringFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "buildah.build.image.storage-driver",
			Sources:     flags.EnvVars("BUILDAH_BUILD_IMAGE_STORAGE_DRIVER"),
			Usage:       "Specify the storage driver for Buildah.",
			Required:    false,
			Value:       "vfs",
			Destination: &P.Image.StorageDriver,
		},

		// CATEGORY_CONTAINER_FILE

		&cli.StringFlag{
			Category:    CATEGORY_CONTAINER_FILE,
			Name:        "buildah.build.file.context",
			Sources:     flags.EnvVars("BUILDAH_BUILD_FILE_CONTEXT"),
			Usage:       "Containerfile context argument for build operation.",
			Required:    false,
			Value:       ".",
			Destination: &P.File.Context,
		},

		&cli.StringFlag{
			Category:    CATEGORY_CONTAINER_FILE,
			Name:        "buildah.build.file.name",
			Sources:     flags.EnvVars("BUILDAH_BUILD_FILE_NAME"),
			Usage:       "Containerfile path for the build operation",
			Required:    false,
			Value:       "Dockerfile",
			Destination: &P.File.Name,
		},

		// CATEGORY_CONTAINER_MANIFEST

		&cli.StringFlag{
			Category:    CATEGORY_CONTAINER_MANIFEST,
			Name:        "buildah.build.manifest.target",
			Sources:     flags.EnvVars("BUILDAH_BUILD_MANIFEST_TARGET"),
			Usage:       "Target image names for patching the manifest. format(Template([]string))",
			Required:    false,
			Destination: &P.Manifest.Target,
		},

		&cli.StringFlag{
			Category:    CATEGORY_CONTAINER_MANIFEST,
			Name:        "buildah.build.manifest.file",
			Sources:     flags.EnvVars("BUILDAH_BUILD_MANIFEST_FILE"),
			Usage:       "Write all the images that are published in to a file for later use. format(Template([]string))",
			Value:       `.published-container-images_{{ $ | join "," | sha256sum }}`,
			Destination: &P.Manifest.File,
		},
	})
