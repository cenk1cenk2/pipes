package login

import (
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

// P is the registry the pipe authenticates against. The build and the manifest
// commands read its uri back, since that is what they prefix every image they
// publish with.
var P = &registry.Credentials{}

// Step is the login stage of every buildah command, since an image that can not
// be pushed is not worth the time it takes to build.
var Step = registry.LoginStep(Spec, P, "buildah", "login")
