package build

import (
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
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_BUILDX_PLATFORMS"),
		),
		Usage:       "Platform arguments for Docker BuildX.",
		Required:    false,
		Value:       "linux/amd64",
		Destination: &P.Docker.BuildxPlatforms,
	},

	// setup.CATEGORY_DOCKER_IMAGE

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.name",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_NAME"),
		),
		Usage:       "Image name for the will be built Docker image.",
		Required:    true,
		Destination: &P.DockerImage.Name,
	},

	&cli.StringSliceFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.tags",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_TAGS"),
		),
		Usage:       "Image tag for the will be built Docker image.",
		Required:    true,
		Destination: &P.DockerImage.Tags,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-file.context",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKERFILE_CONTEXT"),
		),
		Usage:       "Dockerfile context argument for build operation.",
		Required:    false,
		Value:       ".",
		Destination: &P.DockerFile.Context,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-file.name",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKERFILE_NAME"),
		),
		Usage:       "Dockerfile path for the build operation",
		Required:    false,
		Value:       "Dockerfile",
		Destination: &P.DockerFile.Name,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.tag-as-latest",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_TAGS_AS_LATEST"),
		),
		Usage: `Regex pattern to tag the image as latest.
      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags.
      json(RegExp[])`,
		Required: false,
		Value:    flags.FLAG_DEFAULT_DOCKER_IMAGE_TAG_AS_LATEST,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.sanitize-tags",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_SANITIZE_TAGS"),
		),
		Usage: `Sanitizes the given regex pattern out of tag name.
      Template is interpolated with the given matches in the regular expression.
      json([]struct { match: RegExp, template: Template[string](RegExpMatch) })`,
		Required: false,
		Value:    flags.FLAG_DEFAULT_DOCKER_IMAGE_SANITIZE_TAGS,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.tags-template",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_TAGS_TEMPLATE"),
		),
		Usage: `Modifies every tag that matches a certain condition.
      Template is interpolated with the given matches in the regular expression.
      json([]struct { match: RegExp, template: Template[string](RegExpMatch) })`,
		Required: false,
		Value:    "[]",
	},

	&cli.BoolFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.inspect",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_INSPECT"),
		),
		Usage:       "Inspect after pushing the image.",
		Required:    false,
		Value:       true,
		Destination: &P.DockerImage.Inspect,
	},

	&cli.StringSliceFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.build-args",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_BUILD_ARGS"),
		),
		Usage: `Pass in extra build arguments for image.
      You can use it as a template with environment variables as the context.
      format(map[string]Template[string]())`,
		Required:    false,
		Destination: &P.DockerImage.BuildArgs,
	},

	&cli.BoolFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-image.pull",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_IMAGE_PULL"),
		),
		Usage:       "Pull before building the image.",
		Required:    false,
		Value:       true,
		Destination: &P.DockerImage.Pull,
	},

	// CATEGORY_DOCKER_MANIFEST

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_MANIFEST,
		Name:     "docker_manifest.target",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_MANIFEST_TARGET"),
		),
		Usage:       "Target image names for patching the manifest. format(Template[string]([]string))",
		Required:    false,
		Destination: &P.DockerManifest.Target,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker-manifest.output-file",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_MANIFEST_OUTPUT_FILE"),
		),
		Usage:       "Write all the images that are published in to a file for later use. format(Template[string]([]string))",
		Required:    false,
		Value:       `.published-docker-images_{{ $ | join "," | sha256sum }}`,
		Destination: &P.DockerManifest.OutputFile,
	},
})
