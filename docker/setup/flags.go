package setup

import (
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/common/flags"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

const (
	CATEGORY_DOCKER          = "Docker"
	CATEGORY_DOCKER_REGISTRY = "Docker Registry"
	CATEGORY_DOCKER_IMAGE    = "Docker Image"
	CATEGORY_DOCKER_MANIFEST = "Docker Manifest"
)

var Flags = []cli.Flag{

	// CATEGORY_DOCKER

	&cli.BoolFlag{
		Category:    CATEGORY_DOCKER,
		Name:        "docker.use_buildkit",
		Usage:       "Use Docker BuildKit for building images.",
		Required:    false,
		EnvVars:     []string{"DOCKER_USE_BUILDKIT", "DOCKER_BUILDKIT"},
		Value:       flags.FLAG_DOCKER_USE_BUILD_KIT,
		Destination: &TL.Pipe.Docker.UseBuildKit,
	},

	&cli.BoolFlag{
		Category:    CATEGORY_DOCKER,
		Name:        "docker.use_buildx",
		Usage:       "Use Docker BuildX builder for multi-platform builds.",
		Required:    false,
		EnvVars:     []string{"DOCKER_USE_BUILDX"},
		Value:       false,
		Destination: &TL.Pipe.Docker.UseBuildx,
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
}

func ProcessFlags(_ *TaskList[Pipe]) error {
	return nil
}
