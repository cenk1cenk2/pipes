package pipe

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/common/flags"
)

//revive:disable:line-length-limit

const (
	CATEGORY_DOCKER          = "Docker"
	CATEGORY_DOCKER_REGISTRY = "Registry"
	CATEGORY_DOCKER_IMAGE    = "Image"
)

var Flags = TL.Plumber.AppendFlags(flags.NewGitFlags(flags.GitFlagsDestination{
	GitBranch: &TL.Pipe.Git.Branch,
	GitTag:    &TL.Pipe.Git.Tag,
}), []cli.Flag{

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
		EnvVars:     []string{"IMAGE_NAME"},
		Destination: &TL.Pipe.DockerImage.Name,
	},

	&cli.StringSliceFlag{
		Category:    CATEGORY_DOCKER_IMAGE,
		Name:        "docker_image.tags",
		Usage:       "Image tag for will be built Docker image.",
		Required:    true,
		EnvVars:     []string{"IMAGE_TAGS"},
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
		EnvVars:  []string{"IMAGE_TAG_AS_LATEST"},
		Value:    `[ "^tags/v?\\d.\\d.\\d$" ]`,
		Action: func(ctx *cli.Context, s string) error {
			if err := json.Unmarshal([]byte(s), &TL.Pipe.DockerImage.TagAsLatest); err != nil {
				return fmt.Errorf("Can not unmarshal Docker image tags for latest: %w", err)
			}

			return nil
		},
	},

	&cli.StringFlag{
		Category: CATEGORY_DOCKER_IMAGE,
		Name:     "docker_image.sanitize_tags",
		Usage:    `Sanitizes the given regex pattern out of tag name. Template is interpolated with the given matches in the regular expression. json(map[RegExp]Template[[]string])`,
		Required: false,
		EnvVars:  []string{"IMAGE_SANITIZE_TAGS"},
		Value:    `{ "([^/]*)/(.*)": "{{ index $ 1 | to_upper_case }}_{{ index $ 2 }}" }`,
		Action: func(ctx *cli.Context, s string) error {
			if err := json.Unmarshal([]byte(s), &TL.Pipe.DockerImage.TagsSanitize); err != nil {
				return fmt.Errorf("Can not unmarshal Docker image sanitizing tag conditions: %w", err)
			}

			return nil
		},
	},

	&cli.BoolFlag{
		Category:    CATEGORY_DOCKER_IMAGE,
		Name:        "docker_image.pull",
		Usage:       "Pull before building the image.",
		Required:    false,
		EnvVars:     []string{"IMAGE_PULL"},
		Value:       true,
		Destination: &TL.Pipe.DockerImage.Pull,
	},

	&cli.StringFlag{
		Category:    CATEGORY_DOCKER_IMAGE,
		Name:        "docker_image.tags_file",
		Usage:       "Read tags from a file.",
		Required:    false,
		EnvVars:     []string{"TAGS_FILE"},
		Value:       "",
		Destination: &TL.Pipe.DockerImage.TagsFile,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_DOCKER_IMAGE,
		Name:        "docker_image.tags_file_ignore_missing",
		Usage:       "Ignore the missing tags file and contunie operation as expected in that case.",
		Required:    false,
		EnvVars:     []string{"TAGS_FILE_IGNORE_MISSING"},
		Value:       false,
		Destination: &TL.Pipe.DockerImage.TagsFileIgnoreMissing,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_DOCKER_IMAGE,
		Name:        "docker_image.inspect",
		Usage:       "Inspect after pushing the image.",
		Required:    false,
		EnvVars:     []string{"IMAGE_INSPECT"},
		Value:       true,
		Destination: &TL.Pipe.DockerImage.Inspect,
	},

	&cli.StringSliceFlag{
		Category:    CATEGORY_DOCKER_IMAGE,
		Name:        "docker_image.build_args",
		Usage:       "Pass in extra build arguments for image.",
		Required:    false,
		EnvVars:     []string{"BUILD_ARGS"},
		Destination: &TL.Pipe.DockerImage.BuildArgs,
	},
})
