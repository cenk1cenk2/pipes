package git

import (
	ucli "github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

func NewFlags(dst *Refs) []ucli.Flag {
	return []ucli.Flag{
		&ucli.StringFlag{
			Category:    cli.CATEGORY_GIT,
			Name:        "git.branch",
			Sources:     cli.EnvVars("CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"),
			Usage:       "Source control branch.",
			Required:    false,
			Value:       "",
			Destination: &dst.Branch,
		},

		&ucli.StringFlag{
			Category:    cli.CATEGORY_GIT,
			Name:        "git.tag",
			Sources:     cli.EnvVars("CI_COMMIT_TAG", "BITBUCKET_TAG"),
			Usage:       "Source control tag.",
			Required:    false,
			Value:       "",
			Destination: &dst.Tag,
		},
	}
}
