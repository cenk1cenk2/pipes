package pipe

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v45/github"
	. "gitlab.kilic.dev/libraries/plumber/v3"
	"golang.org/x/oauth2"
)

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		Set(func(t *Task[Pipe]) error {
			// create client in any case
			t.Pipe.Ctx.Client = github.NewClient(nil)

			return nil
		})
}

func GithubLogin(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("login").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Github.Token == ""
		}).
		Set(func(t *Task[Pipe]) error {
			ctx := context.Background()
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: t.Pipe.Github.Token},
			)
			tc := oauth2.NewClient(ctx, ts)

			t.Pipe.Ctx.Client = github.NewClient(tc)

			return nil
		})
}

func FetchLatestTag(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("fetch").
		ShouldRunBefore(func(t *Task[Pipe]) error {
			target := strings.Split(t.Pipe.Github.Repository, "/")

			if len(target) != 2 {
				return fmt.Errorf(
					`Repository name should be valid where "%s" was provided.`,
					t.Pipe.Github.Repository,
				)
			}

			t.Pipe.Ctx.Owner = target[0]
			t.Pipe.Ctx.Repository = target[1]

			return nil

		}).
		Set(func(t *Task[Pipe]) error {
			tags, _, err := t.Pipe.Ctx.Client.Repositories.ListTags(
				context.Background(),
				t.Pipe.Ctx.Owner,
				t.Pipe.Ctx.Repository,
				&github.ListOptions{PerPage: 1},
			)

			if err != nil {
				return err
			}

			if len(tags) != 1 {
				return fmt.Errorf(
					"Repository does not contain any tags: %s",
					t.Pipe.Github.Repository,
				)
			}

			latest := tags[0]

			t.Pipe.Ctx.LatestTag = latest

			t.Log.Infoln(
				"Latest tag for repository: %s > %s",
				t.Pipe.Github.Repository,
				latest.GetName(),
			)

			return nil

		})
}

func WriteTagsFile(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("write").
		Set(func(t *Task[Pipe]) error {
			f, err := os.Create(t.Pipe.DockerImage.TagsFile)

			if err != nil {
				return err
			}

			defer f.Close()

			if _, err = f.WriteString(t.Pipe.Ctx.LatestTag.GetName()); err != nil {
				return err
			}

			return nil

		})
}
