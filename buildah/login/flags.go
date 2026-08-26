package login

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

//revive:disable:line-length-limit

const (
	CATEGORY_CONTAINER_REGISTRY = "Container Registry"
)

var Spec = registry.Spec{
	Category:  CATEGORY_CONTAINER_REGISTRY,
	Label:     "Container registry",
	Prefix:    "buildah",
	Command:   "login",
	LegacyEnv: "CONTAINER",
}

var Flags = []cli.Flag(registry.NewFlags(Spec, P))
