// Command pipe-semantic-release releases applications through the
// semantic-release library.
package main

import (
	"context"

	"github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/semantic-release/release"
	"gitlab.kilic.dev/devops/pipes/semantic-release/setup"
)

func newCommand(p *plumber.Plumber) *ucli.Command {
	// The environment feature is opt-in for this pipe, unlike the pipes that own
	// their environment. The flags are shared package level values, so this runs
	// before the command tree reads them.
	plumber.OverwriteCliFlag(setup.EnvironmentFlags, func(f *ucli.BoolFlag) bool {
		return f.Name == "environment.enable"
	}, func(f *ucli.BoolFlag) *ucli.BoolFlag {
		f.Hidden = false
		f.Value = false

		return f
	})

	return &ucli.Command{
		Name:        CLI_NAME,
		Version:     VERSION,
		Usage:       DESCRIPTION,
		Description: DESCRIPTION,
		Flags:       plumber.CombineFlags(setup.EnvironmentFlags, setup.NodeFlags, setup.LoginFlags, release.Flags),
		Action: func(_ context.Context, _ *ucli.Command) error {
			return p.RunJobs(plumber.CombineTaskLists(
				setup.NewEnvironment(p),
				setup.New(p),
				setup.NewLogin(p),
				release.New(p),
			))
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
