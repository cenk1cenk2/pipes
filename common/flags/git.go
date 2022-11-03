package flags

import "github.com/urfave/cli/v2"

type GitFlagsDestination struct {
	GitBranch *string
	GitTag    *string
}

func NewGitFlags(destination GitFlagsDestination) []cli.Flag {
	return []cli.Flag{
		// category_git
		&cli.StringFlag{
			Category:    category_git,
			Name:        "git.branch",
			Usage:       "Source control branch.",
			Required:    false,
			EnvVars:     []string{"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
			Value:       "",
			Destination: destination.GitBranch,
		},

		&cli.StringFlag{
			Category:    category_git,
			Name:        "git.tag",
			Usage:       "Source control tag.",
			Required:    false,
			EnvVars:     []string{"CI_COMMIT_TAG", "BITBUCKET_TAG"},
			Value:       "",
			Destination: destination.GitTag,
		},
	}
}
