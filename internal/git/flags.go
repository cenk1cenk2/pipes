package git

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_GIT = "GIT"
)

// Options is what the git flags are built onto. Every flag constructor takes a
// struct, so a new knob is a new field and not a signature change at every call site.
type Options struct {
	Destination *Refs
}

func NewFlags(opts Options) []cli.Flag {
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
			Destination: &opts.Destination.Branch,
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
			Destination: &opts.Destination.Tag,
		},
	}
}
