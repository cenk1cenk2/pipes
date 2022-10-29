package pipe

import (
	"github.com/urfave/cli/v2"
)

//revive:disable:line-length-limit

const (
	category_git             = "GIT"
	category_docker          = "Docker"
	category_docker_registry = "Registry"
	category_docker_image    = "Image"
)

var Flags = []cli.Flag{
	// category_docker

	&cli.StringFlag{
		Category:    category_git,
		Name:        "git.branch",
		Usage:       "Source control branch.",
		Required:    false,
		EnvVars:     []string{"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
		Value:       "",
		Destination: &TL.Pipe.Git.Branch,
	},

	&cli.StringFlag{
		Category:    category_git,
		Name:        "git.tag",
		Usage:       "Source control tag.",
		Required:    false,
		EnvVars:     []string{"CI_COMMIT_TAG", "BITBUCKET_TAG"},
		Value:       "",
		Destination: &TL.Pipe.Git.Tag,
	},

	// category_docker

	&cli.BoolFlag{
		Category:    category_docker,
		Name:        "docker.use_buildkit",
		Usage:       "Use Docker build kit for building images.",
		Required:    false,
		EnvVars:     []string{"DOCKER_USE_BUILDKIT"},
		Value:       true,
		Destination: &TL.Pipe.Docker.UseBuildKit,
	},

	&cli.BoolFlag{
		Category:    category_docker,
		Name:        "docker.use_buildx",
		Usage:       "Use Docker BuildX builder.",
		Required:    false,
		EnvVars:     []string{"DOCKER_USE_BUILDX"},
		Value:       false,
		Destination: &TL.Pipe.Docker.UseBuildx,
	},

	&cli.StringFlag{
		Category:    category_docker,
		Name:        "docker.buildx_platforms",
		Usage:       "Platform arguments for Docker BuildX.",
		Required:    false,
		EnvVars:     []string{"DOCKER_BUILDX_PLATFORMS"},
		Value:       "linux/amd64",
		Destination: &TL.Pipe.Docker.BuildxPlatforms,
	},

	// category_docker_registry

	&cli.StringFlag{
		Category:    category_docker_registry,
		Name:        "docker_registry.registry",
		Usage:       "Docker registry url for logging in.",
		Required:    false,
		EnvVars:     []string{"DOCKER_REGISTRY"},
		Destination: &TL.Pipe.DockerRegistry.Registry,
	},

	&cli.StringFlag{
		Category:    category_docker_registry,
		Name:        "docker_registry.username",
		Usage:       "Docker registry username for the given registry.",
		Required:    true,
		EnvVars:     []string{"DOCKER_REGISTRY_USERNAME"},
		Destination: &TL.Pipe.DockerRegistry.Username,
	},

	&cli.StringFlag{
		Category:    category_docker_registry,
		Name:        "docker_registry.password",
		Usage:       "Docker registry password for the given registry.",
		Required:    true,
		EnvVars:     []string{"DOCKER_REGISTRY_PASSWORD"},
		Destination: &TL.Pipe.DockerRegistry.Password,
	},

	// category_docker_image

	&cli.StringFlag{
		Category:    category_docker_image,
		Name:        "docker_image.name",
		Usage:       "Image name for will be built Docker image.",
		Required:    true,
		EnvVars:     []string{"IMAGE_NAME"},
		Destination: &TL.Pipe.DockerImage.Name,
	},

	&cli.StringSliceFlag{
		Category:    category_docker_image,
		Name:        "docker_image.tags",
		Usage:       "Image tag for will be built Docker image.",
		Required:    true,
		EnvVars:     []string{"IMAGE_TAGS"},
		Destination: &TL.Pipe.DockerImage.Tags,
	},

	&cli.StringFlag{
		Category:    category_docker_image,
		Name:        "docker_file.context",
		Usage:       "Dockerfile context argument for build operation.",
		Required:    false,
		EnvVars:     []string{"DOCKERFILE_CONTEXT"},
		Value:       ".",
		Destination: &TL.Pipe.DockerFile.Context,
	},

	&cli.StringFlag{
		Category:    category_docker_image,
		Name:        "docker_file.name",
		Usage:       "Dockerfile path for the build operation",
		Required:    false,
		EnvVars:     []string{"DOCKERFILE_NAME"},
		Value:       "Dockerfile",
		Destination: &TL.Pipe.DockerFile.Name,
	},

	&cli.StringFlag{
		Category:    category_docker_image,
		Name:        "docker_image.tag_as_latest",
		Usage:       `Regex pattern to tag the image as latest. Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags. format(json(RegExp[]))`,
		Required:    false,
		EnvVars:     []string{"IMAGE_TAG_AS_LATEST"},
		Value:       `[]`,
		Destination: &TL.Pipe.DockerImage.TagAsLatest,
	},

	&cli.StringFlag{
		Category:    category_docker_image,
		Name:        "docker_image.sanitize_tags",
		Usage:       `Sanitizes the given regex pattern out of tag name. Template is interpolated with the given matches in the regular expression. format(json(map[RegExp]Template[[]string]))`,
		Required:    false,
		EnvVars:     []string{"IMAGE_SANITIZE_TAGS"},
		Value:       `{ "([^/]*/(.*))": "{{ .0 | to_upper_case }}_{{ .1 }}" }`,
		Destination: &TL.Pipe.DockerImage.TagsSanitize,
	},

	&cli.BoolFlag{
		Category:    category_docker_image,
		Name:        "docker_image.pull",
		Usage:       "Pull before building the image.",
		Required:    false,
		EnvVars:     []string{"IMAGE_PULL"},
		Value:       true,
		Destination: &TL.Pipe.DockerImage.Pull,
	},

	&cli.StringFlag{
		Category:    category_docker_image,
		Name:        "docker_image.tags_file",
		Usage:       "Read tags from a file.",
		Required:    false,
		EnvVars:     []string{"TAGS_FILE"},
		Value:       "",
		Destination: &TL.Pipe.DockerImage.TagsFile,
	},

	&cli.BoolFlag{
		Category:    category_docker_image,
		Name:        "docker_image.tags_file_ignore_missing",
		Usage:       "Ignore the missing tags file and contunie operation as expected in that case.",
		Required:    false,
		EnvVars:     []string{"TAGS_FILE_IGNORE_MISSING"},
		Value:       false,
		Destination: &TL.Pipe.DockerImage.TagsFileIgnoreMissing,
	},

	&cli.BoolFlag{
		Category:    category_docker_image,
		Name:        "docker_image.inspect",
		Usage:       "Inspect after pushing the image.",
		Required:    false,
		EnvVars:     []string{"IMAGE_INSPECT"},
		Value:       true,
		Destination: &TL.Pipe.DockerImage.Inspect,
	},

	&cli.StringSliceFlag{
		Category:    category_docker_image,
		Name:        "docker_image.build_args",
		Usage:       "Pass in extra build arguments for image.",
		Required:    false,
		EnvVars:     []string{"BUILD_ARGS"},
		Destination: &TL.Pipe.DockerImage.BuildArgs,
	},
}
