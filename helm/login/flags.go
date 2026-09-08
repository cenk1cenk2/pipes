package login

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_HELM_REGISTRY = "Helm Registry"
)

// Flags are declared once for the whole pipe, so every command that logs in registers the same ones.
var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_HELM_REGISTRY,
		Name:     "helm.login.registry.uri",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_LOGIN_REGISTRY_URI"),
			cli.EnvVar("HELM_REGISTRY_URI"),
		),
		Usage:       "Helm registry URL to login to.",
		Required:    false,
		Value:       "docker.io",
		Destination: &P.Uri,
	},

	&cli.StringFlag{
		Category: CATEGORY_HELM_REGISTRY,
		Name:     "helm.login.registry.username",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_LOGIN_REGISTRY_USERNAME"),
			cli.EnvVar("HELM_REGISTRY_USERNAME"),
		),
		Usage:       "Helm registry username for the given registry.",
		Required:    false,
		Destination: &P.Username,
	},

	&cli.StringFlag{
		Category: CATEGORY_HELM_REGISTRY,
		Name:     "helm.login.registry.password",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_LOGIN_REGISTRY_PASSWORD"),
			cli.EnvVar("HELM_REGISTRY_PASSWORD"),
		),
		Usage:       "Helm registry password for the given registry.",
		Required:    false,
		Destination: &P.Password,
	},
}
