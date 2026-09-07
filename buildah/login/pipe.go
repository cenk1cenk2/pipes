package login

import (
	"github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

// P is the registry the pipe authenticates against. The build and the manifest
// commands read its uri back, since that is what they prefix every image they
// publish with.
var P = &registry.Credentials{}

// Flags are declared once for the whole pipe, so every command that logs in
// registers the same flags rather than its own copy of them.
var Flags = registry.NewFlags(Spec, P)

// New is the login stage of every buildah command, since an image that can not
// be pushed is not worth the time it takes to build.
func New(p *plumber.Plumber) *plumber.TaskList {
	return registry.LoginTaskList(p, P, "buildah", "login")
}
