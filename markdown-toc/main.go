package main

import (
	"github.com/urfave/cli/v2"

	"gitlab.kilic.dev/devops/pipes/markdown-toc/pipe"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func main() {
	NewPlumber(
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       pipe.Flags,
				Before: func(ctx *cli.Context) error {
					p.SetDeprecationNotices(pipe.DeprecationNotices)

					return nil
				},
				Action: func(c *cli.Context) error {
					tl := &pipe.TL

					return tl.RunJobs(
						pipe.New(p).SetCliContext(c).Job(),
					)
				},
			}
		}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags:       true,
			ExcludeHelpCommand: true,
		}).
		Run()
}
