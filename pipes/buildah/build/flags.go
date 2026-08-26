package build

import (
	"strings"

	ucli "github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
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
	tagsfile.NewFlags(&P.ContainerImage.TagsFile, "", &P.ContainerImage.TagsFileStrict, false),
	[]ucli.Flag{

		// CATEGORY_CONTAINER_IMAGE

		&ucli.StringSliceFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "container-image.platforms",
			Sources:     cli.EnvVars("CONTAINER_IMAGE_PLATFORMS"),
			Usage:       "Container image platforms to be built.",
			Required:    false,
			Value:       []string{},
			Destination: &P.ContainerImage.Platforms,
		},

		&ucli.StringFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "container-image.name",
			Sources:     cli.EnvVars("CONTAINER_IMAGE_NAME"),
			Usage:       "Image name for the container image to be built.",
			Required:    true,
			Destination: &P.ContainerImage.Name,
		},

		&ucli.StringSliceFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "container-image.tags",
			Sources:     cli.EnvVars("CONTAINER_IMAGE_TAGS"),
			Usage:       "Image tags for the container image to be built.",
			Required:    true,
			Destination: &P.ContainerImage.Tags,
		},

		cli.YAMLFlag(&ucli.StringFlag{
			Category: CATEGORY_CONTAINER_IMAGE,
			Name:     "container-image.tags-template",
			Sources:  cli.EnvVars("CONTAINER_IMAGE_TAGS_TEMPLATE"),
			Usage: strings.TrimSpace(`
    Modifies every tag that matches a certain condition.
    Template is interpolated with the given matches in the regular expression.

    format(yaml([]struct{ match: RegExp, template: Template(match) }))
    `),
			Required: false,
			Value:    "[]",
		}, &P.ContainerImage.TagsTemplate),

		cli.YAMLFlag(&ucli.StringFlag{
			Category: CATEGORY_CONTAINER_IMAGE,
			Name:     "container-image.sanitize-tags",
			Sources:  cli.EnvVars("CONTAINER_IMAGE_SANITIZE_TAGS"),
			Usage: strings.TrimSpace(`
    Sanitizes the given regex pattern out of tag name.
    Template is interpolated with the given matches in the regular expression.

    format(yaml([]struct{ match: RegExp, template: Template(match) }))
    `),
			Required: false,
			Value:    DEFAULT_SANITIZE_TAGS,
		}, &P.ContainerImage.TagsSanitize),

		cli.YAMLFlag(&ucli.StringFlag{
			Category: CATEGORY_CONTAINER_IMAGE,
			Name:     "container-image.tag-as-latest",
			Sources:  cli.EnvVars("CONTAINER_IMAGE_TAGS_AS_LATEST"),
			Usage: strings.TrimSpace(`
    Regex pattern to tag the image as latest.
    Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags.

    format(yaml([]RegExp))
    `),
			Required: false,
			Value:    DEFAULT_TAG_AS_LATEST,
		}, &P.ContainerImage.TagAsLatest),

		&ucli.BoolFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "container-image.pull",
			Sources:     cli.EnvVars("CONTAINER_IMAGE_PULL"),
			Usage:       "Pull before building the image.",
			Required:    false,
			Value:       true,
			Destination: &P.ContainerImage.Pull,
		},

		&ucli.BoolFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "container-image.push",
			Sources:     cli.EnvVars("CONTAINER_IMAGE_PUSH"),
			Usage:       "Push the image after building.",
			Required:    false,
			Value:       true,
			Destination: &P.ContainerImage.Push,
		},

		cli.YAMLFlag(&ucli.StringFlag{
			Category: CATEGORY_CONTAINER_IMAGE,
			Name:     "container-image.build-args",
			Sources:  cli.EnvVars("CONTAINER_IMAGE_BUILD_ARGS"),
			Usage: strings.TrimSpace(`
    Pass in extra build arguments for image.
    You can use it as a template with environment variables as the context.

    format(yaml(map[string]Template()))
    `),
			Required: false,
			Value:    "",
		}, &P.ContainerImage.BuildArgs),

		&ucli.StringFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "container-image.latest-tag",
			Sources:     cli.EnvVars("CONTAINER_IMAGE_LATEST_TAG"),
			Usage:       "Latest tag for the container image where it is marked as latest.",
			Required:    false,
			Value:       "latest",
			Destination: &P.ContainerImage.LatestTag,
		},

		&ucli.StringFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "container-image.cache",
			Sources:     cli.EnvVars("CONTAINER_IMAGE_CACHE"),
			Usage:       "Specify the cache for the container image.",
			Required:    false,
			Value:       "",
			Destination: &P.ContainerImage.Cache,
		},

		&ucli.StringFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "container-image.format",
			Sources:     cli.EnvVars("CONTAINER_IMAGE_FORMAT"),
			Usage:       "Specify the format for Container Image.",
			Required:    false,
			Value:       "oci",
			Destination: &P.ContainerImage.Format,
		},

		&ucli.StringFlag{
			Category:    CATEGORY_CONTAINER_IMAGE,
			Name:        "container-image.storage-driver",
			Sources:     cli.EnvVars("CONTAINER_IMAGE_STORAGE_DRIVER", "BUILDAH_STORAGE_DRIVER"),
			Usage:       "Specify the storage driver for Buildah.",
			Required:    false,
			Value:       "vfs",
			Destination: &P.ContainerImage.StorageDriver,
		},

		// CATEGORY_CONTAINER_FILE

		&ucli.StringFlag{
			Category:    CATEGORY_CONTAINER_FILE,
			Name:        "containerfile.context",
			Sources:     cli.EnvVars("CONTAINER_FILE_CONTEXT"),
			Usage:       "Containerfile context argument for build operation.",
			Required:    false,
			Value:       ".",
			Destination: &P.ContainerFile.Context,
		},

		&ucli.StringFlag{
			Category:    CATEGORY_CONTAINER_FILE,
			Name:        "containerfile.name",
			Sources:     cli.EnvVars("CONTAINER_FILE_NAME"),
			Usage:       "Containerfile path for the build operation",
			Required:    false,
			Value:       "Dockerfile",
			Destination: &P.ContainerFile.Name,
		},

		// CATEGORY_CONTAINER_MANIFEST

		&ucli.StringFlag{
			Category:    CATEGORY_CONTAINER_MANIFEST,
			Name:        "container-manifest.target",
			Sources:     cli.EnvVars("CONTAINER_MANIFEST_TARGET"),
			Usage:       "Target image names for patching the manifest. format(Template([]string))",
			Required:    false,
			Destination: &P.ContainerManifest.Target,
		},

		&ucli.StringFlag{
			Category:    CATEGORY_CONTAINER_MANIFEST,
			Name:        "container-manifest.file",
			Sources:     cli.EnvVars("CONTAINER_MANIFEST_FILE"),
			Usage:       "Write all the images that are published in to a file for later use. format(Template([]string))",
			Value:       `.published-container-images_{{ $ | join "," | sha256sum }}`,
			Destination: &P.ContainerManifest.File,
		},
	})
