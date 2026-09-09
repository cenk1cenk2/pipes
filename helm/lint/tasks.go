package lint

import (
	"context"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/helm/setup"
)

func lint(tl *TaskList) *Task {
	return tl.CreateTask("lint").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"helm",
				"lint",
				".",
			).
				Set(func(_ context.Context, c *Command) error {
					if P.Kubernetes.Version != "" {
						c.AppendArgs("--kube-version", P.Kubernetes.Version)
					}

					return nil
				}).
				SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDefault).
				SetDir(setup.C.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func template(tl *TaskList) *Task {
	return tl.CreateTask("template").
		ShouldDisable(func(t *Task) bool {
			return !P.ShouldTemplate
		}).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"helm",
				"template",
				".",
			).
				SetLogLevel(LogLevelDebug, LogLevelDefault, LogLevelDefault).
				SetDir(setup.C.Cwd).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
