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

func newCommand(p *plumber.Plumber) *ucli.Command {
	tool := setup.Step
	credentials := login.Step

	return cli.App(CLI_NAME, DESCRIPTION, VERSION,
		cli.Command(p, "login", "Login to the given container registries.", tool, credentials),
		cli.Command(p, "build", "Build container images.", tool, credentials, build.Step(build.Deps{Registry: login.P})),
		cli.Command(p, "manifest", "Update manifests of the container images.", tool, credentials, manifest.Step(manifest.Deps{Registry: login.P})),
	)
}

func main() {
	plumber.NewPlumber(newCommand).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
