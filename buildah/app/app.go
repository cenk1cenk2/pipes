// Package app is the pipe itself: the command tree, its name and everything it
// is wired to. main only supplies the version and runs it, so anything that can
// run a *ucli.Command can run the pipe.
package app

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/buildah/build"
	"gitlab.kilic.dev/devops/pipes/buildah/login"
	"gitlab.kilic.dev/devops/pipes/buildah/manifest"
	"gitlab.kilic.dev/devops/pipes/buildah/setup"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

const Name = "pipe-buildah"

const Description = "Builds and publishes container images from CI/CD with buildah.io"

// Options is where a service the pipe reaches outside the machine for would be
// injected. This pipe drives only its own tool, so there is nothing to swap.
type Options struct{}

func New(p *plumber.Plumber, version string, _ Options) *ucli.Command {
	tool := setup.Step
	credentials := login.Step

	return cli.App(Name, Description, version,
		cli.Command(p, "login", "Login to the given container registries.", tool, credentials),
		cli.Command(p, "build", "Build container images.", tool, credentials, build.Step(build.Deps{Registry: login.P})),
		cli.Command(p, "manifest", "Update manifests of the container images.", tool, credentials, manifest.Step(manifest.Deps{Registry: login.P})),
	)
}
