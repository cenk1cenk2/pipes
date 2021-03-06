package pipe

import (
	"github.com/urfave/cli/v2"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:        "git.branch",
		Usage:       "Source control management branch.",
		Required:    false,
		EnvVars:     []string{"CI_COMMIT_REF_NAME"},
		Value:       "",
		Destination: &TL.Pipe.Git.Branch,
	},

	&cli.StringFlag{
		Name:        "git.tag",
		Usage:       "Source control management tag.",
		Required:    false,
		EnvVars:     []string{"CI_COMMIT_TAG"},
		Value:       "",
		Destination: &TL.Pipe.Git.Tag,
	},

	&cli.StringFlag{
		Name:        "docker_image.name",
		Usage:       "Image name for the to be built Docker image.",
		Required:    true,
		EnvVars:     []string{"IMAGE_NAME"},
		Destination: &TL.Pipe.DockerImage.Name,
	},

	&cli.StringSliceFlag{
		Name:        "docker_image.tags",
		Usage:       "Image tag for the to be built Docker image.",
		Required:    true,
		EnvVars:     []string{"IMAGE_TAGS"},
		Destination: &TL.Pipe.DockerImage.Tags,
	},

	&cli.StringFlag{
		Name:        "docker_file.context",
		Usage:       "Context for Dockerfile.",
		Required:    false,
		EnvVars:     []string{"DOCKERFILE_CONTEXT"},
		Value:       ".",
		Destination: &TL.Pipe.DockerFile.Context,
	},

	&cli.StringFlag{
		Name:        "docker_file.name",
		Usage:       "Dockerfile name to build from.",
		Required:    false,
		EnvVars:     []string{"DOCKERFILE_NAME"},
		Value:       "Dockerfile",
		Destination: &TL.Pipe.DockerFile.Name,
	},

	&cli.StringFlag{
		Name:        "docker_registry.registry",
		Usage:       "Docker registry to login to.",
		Required:    false,
		EnvVars:     []string{"DOCKER_REGISTRY"},
		Destination: &TL.Pipe.DockerRegistry.Registry,
	},

	&cli.StringFlag{
		Name:        "docker_registry.username",
		Usage:       "Docker registry username.",
		Required:    false,
		EnvVars:     []string{"DOCKER_REGISTRY_USERNAME"},
		Destination: &TL.Pipe.DockerRegistry.Username,
	},

	&cli.StringFlag{
		Name:        "docker_registry.password",
		Usage:       "Docker registry password.",
		Required:    false,
		EnvVars:     []string{"DOCKER_REGISTRY_PASSWORD"},
		Destination: &TL.Pipe.DockerRegistry.Password,
	},

	&cli.BoolFlag{
		Name:        "docker.use_buildx",
		Usage:       "Use docker buildx builder.",
		Required:    false,
		EnvVars:     []string{"DOCKER_USE_BUILDX"},
		Value:       false,
		Destination: &TL.Pipe.Docker.UseBuildx,
	},

	&cli.StringFlag{
		Name:        "docker.buildx_platforms",
		Usage:       "Platform arguments for docker buildx.",
		Required:    false,
		EnvVars:     []string{"DOCKER_BUILDX_PLATFORMS"},
		Value:       "linux/amd64",
		Destination: &TL.Pipe.Docker.BuildxPlatforms,
	},

	&cli.StringFlag{
		Name:        "docker_image.tag_as_latest_for_tags_regex",
		Usage:       "Regex pattern to tag the image as latest. (format: json(string[]))",
		Required:    false,
		EnvVars:     []string{"TAG_AS_LATEST_FOR_TAGS_REGEX"},
		Value:       `["^v\\d*\\.\\d*\\.\\d*$"]`,
		Destination: &TL.Pipe.DockerImage.TagAsLatestForTagsRegex,
	},

	&cli.StringFlag{
		Name:        "docker_image.tag_as_latest_for_branches_regex",
		Usage:       "Regex pattern to tag the image as latest. (format: json(string[]))",
		Required:    false,
		EnvVars:     []string{"TAG_AS_LATEST_FOR_BRANCHES_REGEX"},
		Value:       `[]`,
		Destination: &TL.Pipe.DockerImage.TagAsLatestForBranchesRegex,
	},

	&cli.BoolFlag{
		Name:        "docker_image.pull",
		Usage:       "Pull while building the image.",
		Required:    false,
		EnvVars:     []string{"IMAGE_PULL"},
		Value:       true,
		Destination: &TL.Pipe.DockerImage.Pull,
	},

	&cli.StringFlag{
		Name:        "docker_image.tags_file",
		Usage:       "Read tags from a file.",
		Required:    false,
		EnvVars:     []string{"TAGS_FILE"},
		Value:       "",
		Destination: &TL.Pipe.DockerImage.TagsFile,
	},

	&cli.BoolFlag{
		Name:        "docker_image.tags_file_ignore_missing",
		Usage:       "Dont finish the task if tags file is set and missing.",
		Required:    false,
		EnvVars:     []string{"TAGS_FILE_IGNORE_MISSING"},
		Value:       false,
		Destination: &TL.Pipe.DockerImage.TagsFileIgnoreMissing,
	},

	&cli.BoolFlag{
		Name:        "docker_image.inspect",
		Usage:       "Inspect after pushing the image.",
		Required:    false,
		EnvVars:     []string{"IMAGE_INSPECT"},
		Value:       true,
		Destination: &TL.Pipe.DockerImage.Inspect,
	},

	&cli.StringSliceFlag{
		Name:        "docker_image.build_args",
		Usage:       "Pass in extra build arguments for image.",
		Required:    false,
		EnvVars:     []string{"BUILD_ARGS"},
		Destination: &TL.Pipe.DockerImage.BuildArgs,
	},
}
