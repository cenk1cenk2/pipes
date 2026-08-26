package login

import (
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

// P is the chart registry the pipe authenticates against.
var P = &registry.Credentials{}

// Step is the login stage of the helm commands that reach the registry, which is
// every command that pulls a dependency or pushes a chart.
var Step = registry.LoginStep(Spec, P, "helm", "registry", "login")
