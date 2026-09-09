package login

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategoryContainerRegistry = "Container Registry"
)

// Flags are declared once for the whole pipe, so every command that logs in registers the same ones.
var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CategoryContainerRegistry,
		Name:     "buildah.login.registry.uri",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("BUILDAH_LOGIN_REGISTRY_URI"),
			cli.EnvVar("CONTAINER_REGISTRY_URI"),
		),
		Usage:       "Container registry URL to login to.",
		Required:    false,
		Value:       "docker.io",
		Destination: &P.Uri,
	},

	&cli.StringFlag{
		Category: CategoryContainerRegistry,
		Name:     "buildah.login.registry.username",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("BUILDAH_LOGIN_REGISTRY_USERNAME"),
			cli.EnvVar("CONTAINER_REGISTRY_USERNAME"),
		),
		Usage:       "Container registry username for the given registry.",
		Required:    false,
		Destination: &P.Username,
	},

	&cli.StringFlag{
		Category: CategoryContainerRegistry,
		Name:     "buildah.login.registry.password",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("BUILDAH_LOGIN_REGISTRY_PASSWORD"),
			cli.EnvVar("CONTAINER_REGISTRY_PASSWORD"),
		),
		Usage:       "Container registry password for the given registry.",
		Required:    false,
		Destination: &P.Password,
	},
}
