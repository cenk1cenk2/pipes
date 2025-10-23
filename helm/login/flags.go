package login

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_HELM_REGISTRY = "Helm Registry"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_HELM_REGISTRY,
		Name:     "helm-registry.uri",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_REGISTRY_URI"),
		),
		Usage:       "Helm registry url to login to.",
		Required:    false,
		Value:       "docker.io",
		Destination: &P.HelmRegistry.Uri,
	},

	&cli.StringFlag{
		Category: CATEGORY_HELM_REGISTRY,
		Name:     "helm-registry.username",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_REGISTRY_USERNAME"),
		),
		Usage:       "Helm registry username for the given registry.",
		Required:    false,
		Destination: &P.HelmRegistry.Username,
	},

	&cli.StringFlag{
		Category: CATEGORY_HELM_REGISTRY,
		Name:     "helm-registry.password",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_REGISTRY_PASSWORD"),
		),
		Usage:       "Helm registry password for the given registry.",
		Required:    false,
		Destination: &P.HelmRegistry.Password,
	},
}
