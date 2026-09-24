package status

import (
	"context"
	"errors"
	"fmt"

	. "github.com/cenk1cenk2/plumber/v7"

	"gitlab.kilic.dev/devops/pipes/github/client"
)

// authenticate prefers minting over a given token, since a status posted at the
// end of a long pipeline can outlive a token minted at its start.
func authenticate(tl *TaskList) *Task {
	return tl.CreateTask("authenticate").
		Set(func(ctx context.Context, t *Task) error {
			if !P.App.Enabled() {
				C.Token = P.Status.Token

				return nil
			}

			t.Log.Info(fmt.Sprintf("Minting installation token for app: %s > %s", P.App.Id, P.App.InstallationId))

			token, err := P.App.Mint(ctx, t.Plumber, C.Client, client.TokenRequest{
				Repositories: []string{C.Repository},
				Permissions:  map[string]string{"statuses": "write"},
			})

			if err != nil {
				return err
			}

			C.Token = token

			return nil
		})
}

func post(tl *TaskList) *Task {
	return tl.CreateTask("post").
		Set(func(ctx context.Context, t *Task) error {
			t.Log.Info(fmt.Sprintf("Posting commit status: %s@%s > %s", P.Status.Project, P.Status.Sha, P.Status.State))

			err := C.Client.CreateCommitStatus(ctx, C.Token, P.Status.Project, P.Status.Sha, client.CommitStatus{
				State:       P.Status.State,
				TargetUrl:   P.Status.TargetUrl,
				Description: P.Status.Description,
				Context:     P.Status.Context,
			})

			var response *client.ResponseError
			if errors.As(err, &response) && response.CommitNotFound() {
				t.Log.Warn(fmt.Sprintf("Commit is not on GitHub, skipping status: %s@%s", P.Status.Project, P.Status.Sha))

				return nil
			}

			return err
		})
}
