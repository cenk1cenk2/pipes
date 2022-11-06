package pipe

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/common/flags"

	. "gitlab.kilic.dev/libraries/plumber/v4"
)

//revive:disable:line-length-limit

const (
	CATEGORY_DOCKER          = "Docker"
	CATEGORY_DOCKER_REGISTRY = "Registry"
	CATEGORY_DOCKER_IMAGE    = "Image"
)

var Flags = TL.Plumber.AppendFlags(flags.NewGitFlags(
	flags.GitFlagsSetup{
		GitBranchDestination: &TL.Pipe.Git.Branch,
		GitTagDestination:    &TL.Pipe.Git.Tag,
	},
), flags.NewTagsFileFlags(
	flags.TagsFileFlagsSetup{
		TagsFileDestination: &TL.Pipe.DockerImage.TagsFile,
		TagsFileRequired:    false,
	},
), flags.NewTagsFileStrictFlags(
	flags.TagsFileStrictFlagsSetup{
		TagsFileStrictDestination: &TL.Pipe.DockerImage.TagsFileStrict,
		TagsFileStrictRequired:    false,
	},
), []cli.Flag{

	// CATEGORY_DOCKER

	&cli.BoolFlag{
		Category:    CATEGORY_DOCKER,
		Name:        "docker.use_buildkit",
		Usage:       "Use Docker build kit for building images.",
		Required:    false,
		EnvVars:     []string{"DOCKER_USE_BUILDKIT"},
		Value:       true,
		Destination: &TL.Pipe.Docker.UseBuildKit,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_DOCKER,
		Name:        "docker.use_buildx",
		Usage:       "Use Docker BuildX builder.",
		Required:    false,
		EnvVars:     []string{"DOCKER_USE_BUILDX"},
		Value:       false,
		Destination: &TL.Pipe.Docker.UseBuildx,
	},

	&cli.StringFlag{
		Category:    CATEGORY_DOCKER,
		Name:        "docker.buildx_platforms",
		Usage:       "Platform arguments for Docker BuildX.",
		Required:    false,
		EnvVars:     []string{"DOCKER_BUILDX_PLATFORMS"},
		Value:       "linux/amd64",
		Destination: &TL.Pipe.Docker.BuildxPlatforms,
	},

	&cli.StringFlag{
		Category:    CATEGORY_DOCKER,
		Name:        "docker.buildx_instance",
		Usage:       "Docker BuildX instance to be started or to use.",
		Required:    false,
		EnvVars:     []string{"DOCKER_BUILDX_INSTANCE"},
		Value:       "CI",
		Destination: &TL.Pipe.Docker.BuildxInstance,
	},

	// CATEGORY_DOCKER_REGISTRY

	&cli.StringFlag{
		Category:    CATEGORY_DOCKER_REGISTRY,
		Name:        "docker_registry.registry",
		Usage:       "Docker registry url for logging in.",
		Required:    false,
		EnvVars:     []string{"DOCKER_REGISTRY"},
		Destination: &TL.Pipe.DockerRegistry.Registry,
	},

	&cli.StringFlag{
		Category:    CATEGORY_DOCKER_REGISTRY,
		Name:        "docker_registry.username",
		Usage:       "Docker registry username for the given registry.",
		Required:    false,
		EnvVars:     []string{"DOCKER_REGISTRY_USERNAME"},
		Destination: &TL.Pipe.DockerRegistry.Username,
	},

	&cli.StringFlag{
		Category:    CATEGORY_DOCKER_REGISTRY,
		Name:        "docker_registry.password",
		Usage:       "Docker registry password for the given registry.",
		Required:    false,
		EnvVars:     []string{"DOCKER_REGISTRY_PASSWORD"},
		Destination: &TL.Pipe.DockerRegistry.Password,
	},

	// CATEGORY_DOCKER_IMAGE

	&cli.StringFlag{
		Category:    CATEGORY_DOCKER_IMAGE,
		Name:        "docker_image.name",
		Usage:       "Image name for will be built Docker image.",
		Required:    true,
		EnvVars:     []string{"IMAGE_NAME", "DOCKER_IMAGE_NAME"},
		Destination: &TL.Pipe.DockerImage.Name,
	},

	&cli.StringSliceFlag{
		Category:    CATEGORY_DOCKER_IMAGE,
		Name:        "docker_image.tags",
		Usage:       "Image tag for will be built Docker image.",
		Required:    true,
		EnvVars:     []string{"IMAGE_TAGS", "DOCKER_IMAGE_TAGS"},
		Destination: &TL.Pipe.DockerImage.Tags,
	},

	&cli.StringFlag{
		Category:    CATEGORY_DOCKER_IMAGE,
		Name:        "docker_file.context",
		Usage:       "Dockerfile context argument for build operation.",
		Required:    false,
		EnvVars:     []string{"DOCKERFILE_CONTEXT"},
		Value:       ".",
		Destination: &TL.Pipe.DockerFile.Context,
	},

	&cli.StringFlag{
		Category:    CATEGORY_DOCKER_IMAGE,
		Name:        "docker_file.name",
		Usage:       "Dockerfile path for the build operation",
		Required:    false,
		EnvVars:     []string{"DOCKERFILE_NAME"},
		Value:       "Dockerfile",
		Destination: &TL.Pipe.DockerFile.Name,
	},

	&cli.StringFlag{
		Category: CATEGORY_DOCKER_IMAGE,
		Name:     "docker_image.tag_as_latest",
		Usage:    `Regex pattern to tag the image as latest. Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. json(RegExp[])`,
		Required: false,
		EnvVars:  []string{"IMAGE_TAG_AS_LATEST", "DOCKER_IMAGE_TAG_AS_LATEST"},
		Value:    flags.FLAG_DEFAULT_DOCKER_IMAGE_TAG_AS_LATEST,
	},

	&cli.StringFlag{
		Category: CATEGORY_DOCKER_IMAGE,
		Name:     "docker_image.sanitize_tags",
		Usage:    `Sanitizes the given regex pattern out of tag name. Template is interpolated with the given matches in the regular expression. json([]struct{ match: RegExp, template: Template(map[string]string) })`,
		Required: false,
		EnvVars:  []string{"IMAGE_SANITIZE_TAGS", "DOCKER_IMAGE_SANITIZE_TAGS"},
		Value:    flags.FLAG_DEFAULT_DOCKER_IMAGE_SANITIZE_TAGS,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_DOCKER_IMAGE,
		Name:        "docker_image.inspect",
		Usage:       "Inspect after pushing the image.",
		Required:    false,
		EnvVars:     []string{"IMAGE_INSPECT", "DOCKER_IMAGE_INSPECT"},
		Value:       true,
		Destination: &TL.Pipe.DockerImage.Inspect,
	},

	&cli.StringSliceFlag{
		Category:    CATEGORY_DOCKER_IMAGE,
		Name:        "docker_image.build_args",
		Usage:       "Pass in extra build arguments for image.",
		Required:    false,
		EnvVars:     []string{"BUILD_ARGS", "DOCKER_IMAGE_BUILD_ARGS"},
		Destination: &TL.Pipe.DockerImage.BuildArgs,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_DOCKER_IMAGE,
		Name:        "docker_image.pull",
		Usage:       "Pull before building the image.",
		Required:    false,
		EnvVars:     []string{"IMAGE_PULL", "DOCKER_IMAGE_PULL"},
		Value:       true,
		Destination: &TL.Pipe.DockerImage.Pull,
	},
})

func ProcessFlags(tl *TaskList[Pipe]) Job {
	return tl.CreateBasicJob(func() error {
		if err := json.Unmarshal([]byte(tl.CliContext.String("docker_image.tag_as_latest")), &tl.Pipe.DockerImage.TagAsLatest); err != nil {
			return fmt.Errorf("Can not unmarshal Docker image tags for latest: %w", err)
		}

		if err := json.Unmarshal([]byte(tl.CliContext.String("docker_image.sanitize_tags")), &tl.Pipe.DockerImage.TagsSanitize); err != nil {
			return fmt.Errorf("Can not unmarshal Docker image sanitizing tag conditions: %w", err)
		}

		return tl.Validate(tl.Pipe)
	})
}
