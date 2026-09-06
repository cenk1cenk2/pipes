// Command pipe-buildah builds and publishes container images from CI/CD.
package main

import (
	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/buildah/build"
	"gitlab.kilic.dev/devops/pipes/buildah/login"
	"gitlab.kilic.dev/devops/pipes/buildah/manifest"
	"gitlab.kilic.dev/devops/pipes/buildah/setup"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

const name = "pipe-buildah"

const description = "Builds and publishes container images from CI/CD with buildah.io"

var VERSION = "latest"

// options is where a service the pipe reaches outside the machine for would be
// injected. This pipe drives only its own tool, so there is nothing to swap.
type options struct{}

// newCommand builds the command tree. The version is a parameter so a spec can
// build the same tree main does without the build stamp.
func newCommand(p *plumber.Plumber, version string, _ options) *ucli.Command {
	tool := setup.Step
	credentials := login.Step

	return cli.App(name, description, version,
		cli.Command(p, "login", "Login to the given container registries.", tool, credentials),
		cli.Command(p, "build", "Build container images.", tool, credentials, build.Step(build.Deps{Registry: login.P})),
		cli.Command(p, "manifest", "Update manifests of the container images.", tool, credentials, manifest.Step(manifest.Deps{Registry: login.P})),
	)
}

func main() {
	plumber.NewPlumber(func(p *plumber.Plumber) *ucli.Command {
		return newCommand(p, VERSION, options{})
	}).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
