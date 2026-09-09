package login

import (
	"context"
	"fmt"
	"io"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
)

func loginParent(tl *TaskList) *Task {
	return tl.CreateTask("login", "parent").
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobParallel(
				login(tl).Job(),
				loginVerify(tl).Job(),
			)
		})
}

func login(tl *TaskList) *Task {
	return tl.CreateTask("login").
		ShouldDisable(func(t *Task) bool {
			return P.Username == "" ||
				P.Password == ""
		}).
		ShouldRunBefore(func(_ context.Context, t *Task) error {
			t.Plumber.AppendSecrets(P.Password)

			return nil
		}).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"buildah",
				"login",
				P.Uri,
				"--username",
				P.Username,
				"--password-stdin",
			).
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDefault).
				Set(func(_ context.Context, c *Command) error {
					c.Log.Info(fmt.Sprintf("Logging in to container registry: %s",
						P.Uri),
					)

					return nil
				}).
				SetStdin(func(c *Command) io.Reader {
					return strings.NewReader(P.Password)
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

// Proves that the ambient login a credential-less pipeline relies on actually
// exists, before the build spends time on the image.
func loginVerify(tl *TaskList) *Task {
	return tl.CreateTask("login", "verify").
		ShouldDisable(func(_ *Task) bool {
			return P.Username != "" && P.Password != ""
		}).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"buildah",
				"login",
				P.Uri,
			).
				SetLogLevel(LogLevelDebug, LogLevelDefault, LogLevelDefault).
				Set(func(_ context.Context, c *Command) error {
					c.Log.Debug(fmt.Sprintf("Will verify authentication in to container registry: %s",
						P.Uri),
					)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
