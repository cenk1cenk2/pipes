package build

import (
	"fmt"
	"strings"

	"go.yaml.in/yaml/v4"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/common/flags"

	. "github.com/cenk1cenk2/plumber/v6"
)

const (
	CATEGORY_CONTAINER_IMAGE    = "Container Image"
	CATEGORY_CONTAINER_FILE     = "Containerfile"
	CATEGORY_CONTAINER_MANIFEST = "Container Manifest"
)

//revive:disable:line-length-limit

var Flags = CombineFlags(flags.NewGitFlags(
	flags.GitFlagsSetup{
		GitBranchDestination: &P.Git.Branch,
		GitTagDestination:    &P.Git.Tag,
	},
), flags.NewTagsFileFlags(
	flags.TagsFileFlagsSetup{
		TagsFileDestination: &P.ContainerImage.TagsFile,
		TagsFileRequired:    false,
	},
), flags.NewTagsFileStrictFlags(
	flags.TagsFileStrictFlagsSetup{
		TagsFileStrictDestination: &P.ContainerImage.TagsFileStrict,
		TagsFileStrictRequired:    false,
	},
), []cli.Flag{

	// CATEGORY_CONTAINER_IMAGE

	&cli.StringSliceFlag{
		Category: CATEGORY_CONTAINER_IMAGE,
		Name:     "container-image.platforms",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_IMAGE_PLATFORMS"),
		),
		Usage:       "Container image platforms to be built.",
		Required:    false,
		Value:       []string{},
		Destination: &P.ContainerImage.Platforms,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_IMAGE,
		Name:     "container-image.name",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_IMAGE_NAME"),
		),
		Usage:       "Image name for the container image to be built.",
		Required:    true,
		Destination: &P.ContainerImage.Name,
	},

	&cli.StringSliceFlag{
		Category: CATEGORY_CONTAINER_IMAGE,
		Name:     "container-image.tags",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_IMAGE_TAGS"),
		),
		Usage:       "Image tags for the container image to be built.",
		Required:    true,
		Destination: &P.ContainerImage.Tags,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_IMAGE,
		Name:     "container-image.tags-template",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_IMAGE_TAGS_TEMPLATE"),
		),
		Usage: strings.TrimSpace(`
    Modifies every tag that matches a certain condition.
    Template is interpolated with the given matches in the regular expression.

    format(yaml([]struct{ match: RegExp, template: Template(match) }))
    `),
		Required:         false,
		Value:            "[]",
		ValidateDefaults: true,
		Validator: func(v string) error {
			if v == "" {
				return nil
			}

			if err := yaml.Unmarshal([]byte(v), &P.ContainerImage.TagsTemplate); err != nil {
				return fmt.Errorf("Cannot unmarshal container image templating tag conditions: %w", err)
			}

			return nil
		},
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_IMAGE,
		Name:     "container-image.sanitize-tags",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_IMAGE_SANITIZE_TAGS"),
		),
		Usage: strings.TrimSpace(`
    Sanitizes the given regex pattern out of tag name.
    Template is interpolated with the given matches in the regular expression.

    format(yaml([]struct{ match: RegExp, template: Template(match) }))
    `),
		Required:         false,
		Value:            flags.FLAG_DEFAULT_DOCKER_IMAGE_SANITIZE_TAGS,
		ValidateDefaults: true,
		Validator: func(v string) error {
			if v == "" {
				return nil
			}

			if err := yaml.Unmarshal([]byte(v), &P.ContainerImage.TagsSanitize); err != nil {
				return fmt.Errorf("Cannot unmarshal container image sanitizing tag conditions: %w", err)
			}

			return nil
		},
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_IMAGE,
		Name:     "container-image.tag-as-latest",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_IMAGE_TAGS_AS_LATEST"),
		),
		Usage: strings.TrimSpace(`
    Regex pattern to tag the image as latest.
    Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags.

    format(yaml([]RegExp))
    `),
		Required:         false,
		Value:            flags.FLAG_DEFAULT_DOCKER_IMAGE_TAG_AS_LATEST,
		ValidateDefaults: true,
		Validator: func(v string) error {
			if v == "" {
				return nil
			}

			if err := yaml.Unmarshal([]byte(v), &P.ContainerImage.TagAsLatest); err != nil {
				return fmt.Errorf("Cannot unmarshal container image tags for latest: %w", err)
			}

			return nil
		},
	},

	&cli.BoolFlag{
		Category: CATEGORY_CONTAINER_IMAGE,
		Name:     "container-image.pull",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_IMAGE_PULL"),
		),
		Usage:       "Pull before building the image.",
		Required:    false,
		Value:       true,
		Destination: &P.ContainerImage.Pull,
	},

	&cli.BoolFlag{
		Category: CATEGORY_CONTAINER_IMAGE,
		Name:     "container-image.push",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_IMAGE_PUSH"),
		),
		Usage:       "Push the image after building.",
		Required:    false,
		Value:       true,
		Destination: &P.ContainerImage.Push,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_IMAGE,
		Name:     "container-image.build-args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_IMAGE_BUILD_ARGS"),
		),
		Usage: strings.TrimSpace(`
    Pass in extra build arguments for image.
    You can use it as a template with environment variables as the context.

    format(yaml(map[string]Template()))
    `),
		Required:         false,
		Value:            "",
		ValidateDefaults: true,
		Validator: func(v string) error {
			if v == "" {
				return nil
			}

			if err := yaml.Unmarshal([]byte(v), &P.ContainerImage.BuildArgs); err != nil {
				return fmt.Errorf("Cannot unmarshal build arguments: %w", err)
			}

			return nil
		},
	},

	&cli.BoolFlag{
		Category: CATEGORY_CONTAINER_IMAGE,
		Name:     "container-image.inspect",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_IMAGE_INSPECT"),
		),
		Usage:       "Inspect after pushing the image.",
		Required:    false,
		Value:       true,
		Destination: &P.ContainerImage.Inspect,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_IMAGE,
		Name:     "container-image.latest-tag",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_IMAGE_LATEST_TAG"),
		),
		Usage:       "Latest tag for the container image where it is marked as latest.",
		Required:    false,
		Value:       "latest",
		Destination: &P.ContainerImage.LatestTag,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_IMAGE,
		Name:     "container-image.cache",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_IMAGE_CACHE"),
		),
		Usage:       "Specify the cache for the container image.",
		Required:    false,
		Value:       "",
		Destination: &P.ContainerImage.Cache,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_IMAGE,
		Name:     "container-image.format",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_IMAGE_FORMAT"),
		),
		Usage:       "Specify the format for Container Image.",
		Required:    false,
		Value:       "docker",
		Destination: &P.ContainerImage.Format,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_IMAGE,
		Name:     "container-image.storage-driver",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_IMAGE_STORAGE_DRIVER"),
			cli.EnvVar("BUILDAH_STORAGE_DRIVER"),
		),
		Usage:       "Specify the storage driver for Buildah.",
		Required:    false,
		Value:       "vfs",
		Destination: &P.ContainerImage.StorageDriver,
	},

	// CATEGORY_CONTAINER_FILE

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_FILE,
		Name:     "containerfile.context",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_FILE_CONTEXT"),
		),
		Usage:       "Containerfile context argument for build operation.",
		Required:    false,
		Value:       ".",
		Destination: &P.ContainerFile.Context,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_FILE,
		Name:     "containerfile.name",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_FILE_NAME"),
		),
		Usage:       "Containerfile path for the build operation",
		Required:    false,
		Value:       "Dockerfile",
		Destination: &P.ContainerFile.Name,
	},

	// CATEGORY_CONTAINER_MANIFEST

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_MANIFEST,
		Name:     "container-manifest.target",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_MANIFEST_TARGET"),
		),
		Usage:       "Target image names for patching the manifest. format(Template([]string))",
		Required:    false,
		Destination: &P.ContainerManifest.Target,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_MANIFEST,
		Name:     "container-manifest.file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_MANIFEST_FILE"),
		),
		Usage:       "Write all the images that are published in to a file for later use. format(Template([]string))",
		Value:       `.published-container-images_{{ $ | join "," | sha256sum }}`,
		Destination: &P.ContainerManifest.File,
	},
})
