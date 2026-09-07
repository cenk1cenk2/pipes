package git

import (
	"github.com/urfave/cli/v3"
)

const (
	CATEGORY_GIT = "GIT"
)

func NewFlags(dst *Refs) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Category: CATEGORY_GIT,
			Name:     "git.branch",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_COMMIT_REF_NAME"),
				cli.EnvVar("BITBUCKET_BRANCH"),
			),
			Usage:       "Source control branch.",
			Required:    false,
			Value:       "",
			Destination: &dst.Branch,
		},

		&cli.StringFlag{
			Category: CATEGORY_GIT,
			Name:     "git.tag",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("CI_COMMIT_TAG"),
				cli.EnvVar("BITBUCKET_TAG"),
			),
			Usage:       "Source control tag.",
			Required:    false,
			Value:       "",
			Destination: &dst.Tag,
		},
	}
}
