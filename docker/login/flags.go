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
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_REGISTRY"),
		),
		Usage:       "Docker registry url to login to.",
		Required:    false,
		Destination: &P.DockerRegistry.Registry,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_REGISTRY,
		Name:     "docker-registry.username",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_REGISTRY_USERNAME"),
		),
		Usage:       "Docker registry username for the given registry.",
		Required:    false,
		Destination: &P.DockerRegistry.Username,
	},

	&cli.StringFlag{
		Category: setup.CATEGORY_DOCKER_REGISTRY,
		Name:     "docker-registry.password",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("DOCKER_REGISTRY_PASSWORD"),
		),
		Usage:       "Docker registry password for the given registry.",
		Required:    false,
		Destination: &P.DockerRegistry.Password,
	},
}
