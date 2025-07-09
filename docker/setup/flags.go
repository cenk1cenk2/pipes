package setup

import (
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/common/flags"
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
		Category: CATEGORY_DOCKER,
		Name:     "docker.use-buildkit",
		Usage:    "Use Docker BuildKit for building images.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_BUILDKIT"),
			cli.EnvVar("DOCKER_USE_BUILDKIT"),
		),
		Value:       flags.FLAG_DOCKER_USE_BUILD_KIT,
		Destination: &P.Docker.UseBuildKit,
	},

	&cli.BoolFlag{
		Category: CATEGORY_DOCKER,
		Name:     "docker.use-buildx",
		Usage:    "Use Docker BuildX builder for multi-platform builds.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_USE_BUILDX"),
		),
		Value:       false,
		Destination: &P.Docker.UseBuildx,
	},

	&cli.StringFlag{
		Category: CATEGORY_DOCKER,
		Name:     "docker.buildx-instance",
		Usage:    "Docker BuildX instance to be started or to use.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_BUILDX_INSTANCE"),
		),
		Value:       "CI",
		Destination: &P.Docker.BuildxInstance,
	},
}
