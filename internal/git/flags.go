package git

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

func NewFlags(dst *Refs) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Category:    flags.CATEGORY_GIT,
			Name:        "git.branch",
			Sources:     flags.EnvVars("CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"),
			Usage:       "Source control branch.",
			Required:    false,
			Value:       "",
			Destination: &dst.Branch,
		},

		&cli.StringFlag{
			Category:    flags.CATEGORY_GIT,
			Name:        "git.tag",
			Sources:     flags.EnvVars("CI_COMMIT_TAG", "BITBUCKET_TAG"),
			Usage:       "Source control tag.",
			Required:    false,
			Value:       "",
			Destination: &dst.Tag,
		},
	}
}
