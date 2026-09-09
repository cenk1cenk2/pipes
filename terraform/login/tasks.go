package login

import (
	"context"
	"fmt"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
)

func environmentCredentials(tl *TaskList) *Task {
	return tl.CreateTask("environment", "credentials").
		ShouldDisable(func(t *Task) bool {
			return len(P.Registry.Credentials) == 0
		}).
		Set(func(_ context.Context, t *Task) error {
			for _, c := range P.Registry.Credentials {
				t.Log.Info(fmt.Sprintf("Generating registry token: %s", c.Registry))

				sanitized := strings.ReplaceAll(c.Registry, ".", "_")
				sanitized = strings.ReplaceAll(sanitized, "-", "__")

				setup.C.Env["TF_TOKEN_"+sanitized] = c.Token
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
