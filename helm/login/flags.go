package login

import (
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

const (
	CategoryHelmRegistry = "Helm Registry"
)

var Flags = registry.NewLoginFlags(registry.LoginFlagsOptions{
	Category:    CategoryHelmRegistry,
	Alias:       "HELM_LOGIN",
	Destination: &P.Container,
})
