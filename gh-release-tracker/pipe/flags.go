package pipe

import (
	"github.com/urfave/cli/v2"
	"gitlab.kilic.dev/devops/pipes/common/flags"
)

//revive:disable:line-length-limit

var Flags = TL.Plumber.AppendFlags(flags.NewTagsFileFlags(
	flags.TagsFileFlagsSetup{
		TagsFileDestination: &TL.Pipe.Tags.File,
		TagsFileRequired:    true,
		TagsFileValue:       ".tags",
	},
), []cli.Flag{
	&cli.StringFlag{
		Category:    flags.CATEGORY_GITHUB,
		Name:        "gh.token",
		Usage:       "Github token for the API requests.",
		Required:    false,
		EnvVars:     []string{"GH_TOKEN", "GITHUB_TOKEN"},
		Value:       "",
		Destination: &TL.Pipe.Github.Token,
	},
	&cli.StringFlag{
		Category:    flags.CATEGORY_GITHUB,
		Name:        "gh.repository",
		Usage:       "Target repository to fetch the latest tag.",
		Required:    true,
		EnvVars:     []string{"GH_REPOSITORY", "GH_REPOSITORY"},
		Value:       "",
		Destination: &TL.Pipe.Github.Repository,
	},
})
