package setup

import (
	"context"
	"fmt"

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
			t.CreateCommand("pulumi", "version").
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				CaptureOutput(&C.Version).
				ShouldRunAfter(func(_ context.Context, c *Command) error {
					c.Log.Info(fmt.Sprintf("pulumi version: %s", C.Version))

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
