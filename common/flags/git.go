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
			Category:    CATEGORY_GIT,
			Name:        "git.branch",
			Usage:       "Source control branch.",
			Required:    false,
			EnvVars:     []string{"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
			Value:       "",
			Destination: setup.GitBranchDestination,
		},

		&cli.StringFlag{
			Category:    CATEGORY_GIT,
			Name:        "git.tag",
			Usage:       "Source control tag.",
			Required:    false,
			EnvVars:     []string{"CI_COMMIT_TAG", "BITBUCKET_TAG"},
			Value:       "",
			Destination: setup.GitTagDestination,
		},
	}
}
