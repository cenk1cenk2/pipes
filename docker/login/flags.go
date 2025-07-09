package login

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{

	// CATEGORY_DOCKER_REGISTRY

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_REGISTRY,
		Name:     "docker-registry.registry",
		Usage:    "Docker registry url to login to.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_REGISTRY"),
		),
		Destination: &P.DockerRegistry.Registry,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_REGISTRY,
		Name:     "docker-registry.username",
		Usage:    "Docker registry username for the given registry.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_REGISTRY_USERNAME"),
		),
		Destination: &P.DockerRegistry.Username,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_REGISTRY,
		Name:     "docker-registry.password",
		Usage:    "Docker registry password for the given registry.",
		Required: false,
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_REGISTRY_PASSWORD"),
		),
		Destination: &P.DockerRegistry.Password,
	},
}
