package build

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/common/flags"
	"gitlab.kilic.dev/devops/pipes/docker/setup"

	. "gitlab.kilic.dev/libraries/plumber/v4"
)

//revive:disable:line-length-limit

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

	&cli.StringFlag{
		Category:    setup.CATEGORY_DOCKER,
		Name:        "docker.buildx_platforms",
		Usage:       "Platform arguments for Docker BuildX.",
		Required:    false,
		EnvVars:     []string{"DOCKER_BUILDX_PLATFORMS"},
		Value:       "linux/amd64",
		Destination: &TL.Pipe.Docker.BuildxPlatforms,
	},

	// setup.CATEGORY_DOCKER_IMAGE

	&cli.StringFlag{
		Category:    setup.CATEGORY_DOCKER_IMAGE,
		Name:        "docker_image.name",
		Usage:       "Image name for the will be built Docker image.",
		Required:    true,
		EnvVars:     []string{"DOCKER_IMAGE_NAME"},
		Destination: &TL.Pipe.DockerImage.Name,
	},

	&cli.StringSliceFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker_image.tags",
		Usage:    "Image tag for the will be built Docker image.",
		Required: true,
		EnvVars:  []string{"DOCKER_IMAGE_TAGS"},
	},

	&cli.StringFlag{
		Category:    setup.CATEGORY_DOCKER_IMAGE,
		Name:        "docker_file.context",
		Usage:       "Dockerfile context argument for build operation.",
		Required:    false,
		EnvVars:     []string{"DOCKERFILE_CONTEXT"},
		Value:       ".",
		Destination: &TL.Pipe.DockerFile.Context,
	},

	&cli.StringFlag{
		Category:    setup.CATEGORY_DOCKER_IMAGE,
		Name:        "docker_file.name",
		Usage:       "Dockerfile path for the build operation",
		Required:    false,
		EnvVars:     []string{"DOCKERFILE_NAME"},
		Value:       "Dockerfile",
		Destination: &TL.Pipe.DockerFile.Name,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker_image.tag_as_latest",
		Usage: `Regex pattern to tag the image as latest.
      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags.
      json(RegExp[])`,
		Required: false,
		EnvVars:  []string{"DOCKER_IMAGE_TAG_AS_LATEST"},
		Value:    flags.FLAG_DEFAULT_DOCKER_IMAGE_TAG_AS_LATEST,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker_image.sanitize_tags",
		Usage: `Sanitizes the given regex pattern out of tag name.
      Template is interpolated with the given matches in the regular expression.
      json([]struct{ match: RegExp, template: Template(map[string]string) })`,
		Required: false,
		EnvVars:  []string{"DOCKER_IMAGE_SANITIZE_TAGS"},
		Value:    flags.FLAG_DEFAULT_DOCKER_IMAGE_SANITIZE_TAGS,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker_image.tags_template",
		Usage: `Modifies every tag that matches a certain condition.
      Template is interpolated with the given matches in the regular expression.
      json([]struct{ match: RegExp, template: Template(map[string]string) })`,
		Required: false,
		EnvVars:  []string{"DOCKER_IMAGE_TAGS_TEMPLATE"},
		Value:    "[]",
	},

	&cli.BoolFlag{
		Category:    setup.CATEGORY_DOCKER_IMAGE,
		Name:        "docker_image.inspect",
		Usage:       "Inspect after pushing the image.",
		Required:    false,
		EnvVars:     []string{"DOCKER_IMAGE_INSPECT"},
		Value:       true,
		Destination: &TL.Pipe.DockerImage.Inspect,
	},

	&cli.StringSliceFlag{
		Category: setup.CATEGORY_DOCKER_IMAGE,
		Name:     "docker_image.build_args",
		Usage:    "Pass in extra build arguments for image.",
		Required: false,
		EnvVars:  []string{"DOCKER_IMAGE_BUILD_ARGS"},
	},

	&cli.BoolFlag{
		Category:    setup.CATEGORY_DOCKER_IMAGE,
		Name:        "docker_image.pull",
		Usage:       "Pull before building the image.",
		Required:    false,
		EnvVars:     []string{"DOCKER_IMAGE_PULL"},
		Value:       true,
		Destination: &TL.Pipe.DockerImage.Pull,
	},

	// CATEGORY_DOCKER_MANIFEST

	&cli.StringFlag{
		Category:    setup.CATEGORY_DOCKER_MANIFEST,
		Name:        "docker_manifest.target",
		Usage:       "Target image names for patching the manifest.",
		Required:    false,
		EnvVars:     []string{"DOCKER_MANIFEST_TARGET"},
		Destination: &TL.Pipe.DockerManifest.Target,
	},

	&cli.StringFlag{
		Category:    setup.CATEGORY_DOCKER_IMAGE,
		Name:        "docker_manifest.output-file",
		Usage:       "Write all the images that are published in to a file for later use. Template(string)",
		Required:    false,
		EnvVars:     []string{"DOCKER_MANIFEST_OUTPUT_FILE"},
		Value:       `.published-docker-images_{{ sha256sum $ }}`,
		Destination: &TL.Pipe.DockerManifest.OutputFile,
	},
})

func ProcessFlags(tl *TaskList[Pipe]) error {
	if v := tl.CliContext.String("docker_image.tag_as_latest"); v != "" {
		if err := json.Unmarshal([]byte(v), &tl.Pipe.DockerImage.TagAsLatest); err != nil {
			return fmt.Errorf("Can not unmarshal Docker image tags for latest: %w", err)
		}
	}

	if v := tl.CliContext.String("docker_image.sanitize_tags"); v != "" {
		if err := json.Unmarshal([]byte(v), &tl.Pipe.DockerImage.TagsSanitize); err != nil {
			return fmt.Errorf("Can not unmarshal Docker image sanitizing tag conditions: %w", err)
		}
	}

	if v := tl.CliContext.String("docker_image.tags_template"); v != "" {
		if err := json.Unmarshal([]byte(v), &tl.Pipe.DockerImage.TagsTemplate); err != nil {
			return fmt.Errorf("Can not unmarshal Docker image templating tag conditions: %w", err)
		}
	}

	tl.Pipe.DockerImage.Tags = tl.CliContext.StringSlice("docker_image.tags")
	tl.Pipe.DockerImage.BuildArgs = tl.CliContext.StringSlice("docker_image.build_args")

	return nil
}

var DeprecationNotices = []DeprecationNotice{
	{
		Level:       LOG_LEVEL_ERROR,
		Environment: []string{"TAG_AS_LATEST_FOR_BRANCHES_REGEX", "TAG_AS_LATEST_FOR_TAGS_REGEX"},
		Flag:        []string{"--docker_image.tag_as_latest_for_branches_regex", "--docker_image.tag_as_latest_for_tags_regex"},
		Message:     `"%s" is deprecated, please use the "TAG_AS_LATEST" format.`,
	},
	{
		Level:       LOG_LEVEL_ERROR,
		Environment: []string{"IMAGE_NAME", "IMAGE_TAGS", "IMAGE_TAG_AS_LATEST", "IMAGE_SANITIZE_TAGS", "IMAGE_SANITIZE_TAGS", "IMAGE_INSPECT", "BUILD_ARGS", "IMAGE_PULL"},
		Message:     `"%s" is deprecated, please use the environment variable with the "DOCKER_" prefix instead.`,
	},
}
