package login

import (
	"github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

// P is the chart registry the pipe authenticates against.
var P = &registry.Credentials{}

// Flags are declared once for the whole pipe, so every command that logs in
// registers the same flags rather than its own copy of them.
var Flags = registry.NewFlags(Spec, P)

// New is the login stage of the helm commands that reach the registry, which is
// every command that pulls a dependency or pushes a chart.
func New(p *plumber.Plumber) *plumber.TaskList {
	return registry.LoginTaskList(p, P, "helm", "registry", "login")
}
