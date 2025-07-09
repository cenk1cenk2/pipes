package build

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/common/flags"
	"gitlab.kilic.dev/devops/pipes/docker/setup"

	. "github.com/cenk1cenk2/plumber/v6"
)

//revive:disable:line-length-limit

var Flags = CombineFlags(flags.NewGitFlags(
	flags.GitFlagsSetup{
		GitBranchDestination: &P.Git.Branch,
		GitTagDestination:    &P.Git.Tag,
	},
), flags.NewTagsFileFlags(
	flags.TagsFileFlagsSetup{
		TagsFileDestination: &P.DockerImage.TagsFile,
		TagsFileRequired:    false,
	},
), flags.NewTagsFileStrictFlags(
	flags.TagsFileStrictFlagsSetup{
		TagsFileStrictDestination: &P.DockerImage.TagsFileStrict,
		TagsFileStrictRequired:    false,
	},
), []cli.Flag{

	// CATEGORY_DOCKER

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER,
		Name:     "docker.buildx-platforms",
		Usage:    "Platform arguments for Docker BuildX.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_BUILDX_PLATFORMS"),
		),
		Value:       "linux/amd64",
		Destination: &P.Docker.BuildxPlatforms,
	},

	// setup.CATEGORY_DOCKER_IMAGE

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.name",
		Usage:    "Image name for the will be built Docker image.",
		Required: true,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_NAME"),
		),
		Destination: &P.DockerImage.Name,
	},

	&cli.StringSliceFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.tags",
		Usage:    "Image tag for the will be built Docker image.",
		Required: true,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_TAGS"),
		),
		Destination: &P.DockerImage.Tags,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-file.context",
		Usage:    "Dockerfile context argument for build operation.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_FILE_CONTEXT"),
		),
		Value:       ".",
		Destination: &P.DockerFile.Context,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-file.name",
		Usage:    "Dockerfile path for the build operation",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKERFILE_NAME"),
		),
		Value:       "Dockerfile",
		Destination: &P.DockerFile.Name,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.tag-as-latest",
		Usage: `Regex pattern to tag the image as latest.
      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags.
      json(RegExp[])`,
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_TAGS_AS_LATEST"),
		),
		Value: flags.FLAG_DEFAULT_DOCKER_IMAGE_TAG_AS_LATEST,
		Action: func(_ context.Context, command *cli.Command, v string) error {
			if err := json.Unmarshal([]byte(v), &P.DockerImage.TagAsLatest); err != nil {
				return fmt.Errorf("Can not unmarshal Docker image tags for latest: %w", err)
			}

			return nil
		},
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.sanitize-tags",
		Usage: `Sanitizes the given regex pattern out of tag name.
      Template is interpolated with the given matches in the regular expression.
      json([]struct { match: RegExp, template: Template[string](RegExpMatch) })`,
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_SANITIZE_TAGS"),
		),
		Value: flags.FLAG_DEFAULT_DOCKER_IMAGE_SANITIZE_TAGS,
		Action: func(_ context.Context, command *cli.Command, v string) error {
			if err := json.Unmarshal([]byte(v), &P.DockerImage.TagsSanitize); err != nil {
				return fmt.Errorf("Can not unmarshal Docker image sanitizing tag conditions: %w", err)
			}

			return nil
		},
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.tags-template",
		Usage: `Modifies every tag that matches a certain condition.
      Template is interpolated with the given matches in the regular expression.
      json([]struct { match: RegExp, template: Template[string](RegExpMatch) })`,
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_TAGS_TEMPLATE"),
		),
		Value: "[]",
		Action: func(_ context.Context, command *cli.Command, v string) error {
			if err := json.Unmarshal([]byte(v), &P.DockerImage.TagsTemplate); err != nil {
				return fmt.Errorf("Can not unmarshal Docker image templating tag conditions: %w", err)
			}

			return nil
		},
	},

	&cli.BoolFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.inspect",
		Usage:    "Inspect after pushing the image.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_INSPECT"),
		),
		Value:       true,
		Destination: &P.DockerImage.Inspect,
	},

	&cli.StringSliceFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.build-args",
		Usage: `Pass in extra build arguments for image.
      You can use it as a template with environment variables as the context.
      format(map[string]Template[string]())`,
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_BUILD_ARGS"),
		),
		Destination: &P.DockerImage.BuildArgs,
	},

	&cli.BoolFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.pull",
		Usage:    "Pull before building the image.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_PULL"),
		),
		Value:       true,
		Destination: &P.DockerImage.Pull,
	},

	// CATEGORY_DOCKER_MANIFEST

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker_manifest.target",
		Usage:    "Target image names for patching the manifest. format(Template[string]([]string))",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_MANIFEST_TARGET"),
		),
		Destination: &P.DockerManifest.Target,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-manifest.output-file",
		Usage:    "Write all the images that are published in to a file for later use. format(Template[string]([]string))",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_MANIFEST_OUTPUT_FILE"),
		),
		Value:       `.published-docker-images_{{ $ | join "," | sha256sum }}`,
		Destination: &P.DockerManifest.OutputFile,
	},
})
