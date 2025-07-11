package flags

import (
	"github.com/urfave/cli/v3"
)

type GitFlagsSetup struct {
	GitBranchDestination *string
	GitTagDestination    *string
}

func NewGitFlags(setup GitFlagsSetup) []cli.Flag {
	return []cli.Flag{
		// CATEGORY_GIT
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
			Destination: setup.GitBranchDestination,
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
			Destination: setup.GitTagDestination,
		},
	}
}
