package login

import (
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

const (
	CategoryContainerRegistry = "Container Registry"
)

var Flags = registry.NewLoginFlags(registry.LoginFlagsOptions{
	Category:    CategoryContainerRegistry,
	Destination: &P.Container,
})
