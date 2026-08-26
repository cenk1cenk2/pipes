package main

import (
	ucli "github.com/urfave/cli/v3"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/kustomize/app"
)

var VERSION = "latest"

func main() {
	NewPlumber(func(p *Plumber) *ucli.Command {
		return app.New(p, VERSION, app.Options{})
	}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
