// Command pipe-go builds Go applications from CI/CD.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/go/build"
	"gitlab.kilic.dev/devops/pipes/go/install"
	"gitlab.kilic.dev/devops/pipes/go/lint"
	"gitlab.kilic.dev/devops/pipes/go/setup"
	gotool "gitlab.kilic.dev/devops/pipes/go/tool"
)

func newCommand(p *plumber.Plumber) *ucli.Command {
	return &ucli.Command{
		Name:        CLI_NAME,
		Version:     VERSION,
		Usage:       DESCRIPTION,
		Description: DESCRIPTION,
		Commands: []*ucli.Command{
			{
				Name:        "install",
				Description: "Vendor go modules.",
				Flags:       plumber.CombineFlags(setup.Flags, install.Flags),
				Action: func(_ context.Context, _ *ucli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						install.New(p, install.Deps{Tool: setup.C}),
					))
				},
			},
			{
				Name:        "build",
				Description: "Build an application.",
				Flags:       plumber.CombineFlags(setup.Flags, build.Flags),
				Action: func(_ context.Context, _ *ucli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						build.New(p, build.Deps{Tool: setup.C.Ctx}),
					))
				},
			},
			{
				Name:        "lint",
				Description: "Run golangci-lint on the project.",
				Flags:       plumber.CombineFlags(setup.Flags, lint.Flags),
				Action: func(_ context.Context, _ *ucli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						lint.New(p, lint.Deps{Tool: setup.C}),
					))
				},
			},
			{
				Name:        "tool",
				Description: "Run a specified go tool.",
				Flags:       plumber.CombineFlags(setup.Flags, gotool.Flags),
				Arguments:   gotool.Arguments,
				Action: func(_ context.Context, _ *ucli.Command) error {
					return p.RunJobs(plumber.CombineTaskLists(
						setup.New(p),
						gotool.New(p, gotool.Deps{Tool: setup.C.Ctx}),
					))
				},
			},
		},
	}
}

func main() {
	plumber.NewPlumber(newCommand).
		SetDocumentationOptions(plumber.DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
