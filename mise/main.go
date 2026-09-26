package main

import (
	"context"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/mise/install"
	"gitlab.kilic.dev/devops/pipes/mise/setup"
)

var version = "latest"

func main() {
	NewPlumber(func(p *Plumber) *cli.Command {
		return &cli.Command{
			Name:    "pipe-mise",
			Version: version,
			Description: strings.TrimSpace(`
Pipe for installing tools with mise.

The image ships in two variants, and the variant decides the libc of every tool mise installs.

| Tag | Base | Flavor |
| --- | --- | --- |
| latest | alpine (musl), the default | musl-amd64 |
| latest-glibc | debian (glibc), for glibc consumers | glibc-amd64 |

The flavor has to match the libc of the image that consumes the cache: a musl cache only serves static or musl-linked tools on a glibc consumer, and a glibc cache serves nothing on alpine. The flavor is part of the cache key, so the two variants never share a cache.
`),
			Commands: []*cli.Command{
				{
					Name: "install",
					Description: strings.TrimSpace(`
Install the tools of the mise configuration into the data directory.

The data directory carries its own copy of the mise binary and the shims link to that copy, so a job that restores the directory runs the tools without mise in its image.
`),
					Flags: CombineFlags(setup.Flags, install.Flags),
					Action: func(_ context.Context, _ *cli.Command) error {
						return p.RunJobs(CombineTaskLists(
							setup.New(p),
							install.New(p),
						))
					},
				},

				DocsCommand(p),
			},
		}
	}).
		SetDocumentationOptions(DocumentationOptions{
			ExcludeFlags: true,
		}).
		Run()
}
