package pipe

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v44/github"
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
	"golang.org/x/oauth2"
)

type Ctx struct {
	Client    *github.Client
	LatestTag *github.RepositoryTag
}

var Context Ctx

func VerifyVariables() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "verify"},
		Task: func(t *utils.Task) error {
			err := utils.ValidateAndSetDefaults(t.Metadata, &Pipe)

			if err != nil {
				return err
			}

			// create client in any case
			Context.Client = github.NewClient(nil)

			return nil
		},
	}
}

func GithubLogin() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "login", Skip: Pipe.Github.Token == ""},
		Task: func(t *utils.Task) error {
			ctx := context.Background()
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: Pipe.Github.Token},
			)
			tc := oauth2.NewClient(ctx, ts)

			Context.Client = github.NewClient(tc)

			return nil
		},
	}
}

func FetchLatestTag() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "fetch"},
		Task: func(t *utils.Task) error {
			target := strings.Split(Pipe.Github.Repository, "/")

			if len(target) != 2 {
				return fmt.Errorf(
					`Repository name should be valid where "%s" was provided.`,
					Pipe.Github.Repository,
				)
			}

			owner := target[0]
			repository := target[1]

			tags, _, err := Context.Client.Repositories.ListTags(
				context.Background(),
				owner,
				repository,
				&github.ListOptions{PerPage: 1},
			)

			if err != nil {
				return err
			}

			if len(tags) != 1 {
				return fmt.Errorf(
					"Repository does not contain any tags: %s",
					Pipe.Github.Repository,
				)
			}

			latest := tags[0]

			Context.LatestTag = latest

			t.Log.Infoln(
				"Latest tag for repository: %s > %s",
				Pipe.Github.Repository,
				latest.GetName(),
			)

			return nil
		},
	}
}

func WriteTagsFile() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "write"},
		Task: func(t *utils.Task) error {
			f, err := os.Create(Pipe.DockerImage.TagsFile)

			if err != nil {
				return err
			}

			defer f.Close()

			_, err = f.WriteString(
				Context.LatestTag.GetName(),
			)

			if err != nil {
				return err
			}

			return nil
		},
	}
}
