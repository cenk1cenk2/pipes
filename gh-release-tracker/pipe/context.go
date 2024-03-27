package pipe

import (
	"github.com/google/go-github/v60/github"
)

type Ctx struct {
	Client     *github.Client
	LatestTag  *github.RepositoryTag
	Owner      string
	Repository string
}
