package setup

import (
	"context"
	"fmt"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
)

func initialize(tl *TaskList) *Task {
	return tl.CreateTask("init").
		Set(func(_ context.Context, t *Task) error {
			C.Cwd = P.Cwd

			t.Log.Debug(fmt.Sprintf("Working directory: %s", C.Cwd))

			return nil
		})
}

func version(tl *TaskList) *Task {
	return tl.CreateTask("version").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand("buildah", "--version").
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				ShouldRunAfter(func(_ context.Context, c *Command) error {
					C.Version = strings.TrimSpace(strings.Join(c.GetCombinedStream(), "\n"))

					c.Log.Info(fmt.Sprintf("buildah version: %s", C.Version))

					return nil
				}).
				EnableStreamRecording().
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
