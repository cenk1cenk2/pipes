package main

import (
	"github.com/urfave/cli/v2"

	"gitlab.kilic.dev/devops/pipes/docker/pipe"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func main() {
	p := Plumber{
		DocsExcludeFlags:       true,
		DocsExcludeHelpCommand: true,
		DeprecationNotices: []DeprecationNotice{
			{
				Level:       LOG_LEVEL_ERROR,
				Environment: []string{"TAG_AS_LATEST_FOR_TAGS_REGEX"},
				Flag:        []string{"--docker_image.tag_as_latest_for_tags_regex"},
			},
			{
				Level:       LOG_LEVEL_ERROR,
				Environment: []string{"TAG_AS_LATEST_FOR_BRANCHES_REGEX"},
				Flag:        []string{"--docker_image.tag_as_latest_for_branches_regex"},
			},
		},
	}

	p.New(
		func(p *Plumber) *cli.App {
			return &cli.App{
				Name:        CLI_NAME,
				Version:     VERSION,
				Usage:       DESCRIPTION,
				Description: DESCRIPTION,
				Flags:       pipe.Flags,
				Action: func(c *cli.Context) error {
					return pipe.TL.RunJobs(
						pipe.New(p).SetCliContext(c).Job(),
					)
				},
			}
		}).Run()
}
