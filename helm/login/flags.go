package login

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

//revive:disable:line-length-limit

const (
	CATEGORY_HELM_REGISTRY = "Helm Registry"
)

var Spec = registry.Spec{
	Category:  CATEGORY_HELM_REGISTRY,
	Label:     "Helm registry",
	Prefix:    "helm",
	Command:   "login",
	LegacyEnv: "HELM",
}

var Flags = []cli.Flag(registry.NewFlags(Spec, P))
