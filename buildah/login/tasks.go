package login

import (
	"context"
	"fmt"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/internal/registry"
)

// Proves that the ambient login a credential-less pipeline relies on actually
// exists, before the build spends time on the image.
func loginVerify(tl *TaskList, container *registry.Login, image string) *Task {
	return tl.CreateTask("login", "verify").
		ShouldDisable(func(_ *Task) bool {
			return image == ""
		}).
		Set(func(_ context.Context, t *Task) error {
			target, err := InlineTemplate[any](image, nil)
			if err != nil {
				return err
			}

			parsed, err := registry.ParseImage(target)
			if err != nil {
				return err
			}

			if container.Covers(parsed) {
				t.Log.Debug(fmt.Sprintf("Container login covers the image, skipping the verification: %s", target))

				return nil
			}

			t.CreateCommand(
				"buildah",
				"login",
				parsed.Host,
			).
				SetLogLevel(LogLevelDebug, LogLevelDefault, LogLevelDefault).
				Set(func(_ context.Context, c *Command) error {
					c.Log.Debug(fmt.Sprintf("Will verify authentication in to container registry: %s",
						parsed.Host),
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
