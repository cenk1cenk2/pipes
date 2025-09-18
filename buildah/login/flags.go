package login

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_CONTAINER_REGISTRY = "Container Registry"
)

var Flags = []cli.Flag{

	// CATEGORY_CONTAINER_REGISTRY

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_REGISTRY,
		Name:     "container-registry.uri",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_REGISTRY_URI"),
		),
		Usage:       "Container registry url to login to.",
		Required:    false,
		Destination: &P.ContainerRegistry.Uri,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_REGISTRY,
		Name:     "container-registry.username",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_REGISTRY_USERNAME"),
		),
		Usage:       "Container registry username for the given registry.",
		Required:    false,
		Destination: &P.ContainerRegistry.Username,
	},

	&cli.StringFlag{
		Category: CATEGORY_CONTAINER_REGISTRY,
		Name:     "container-registry.password",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("CONTAINER_REGISTRY_PASSWORD"),
		),
		Usage:       "Container registry password for the given registry.",
		Required:    false,
		Destination: &P.ContainerRegistry.Password,
	},
}
