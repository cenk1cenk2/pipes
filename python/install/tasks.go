package install

import (
	"context"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/python/setup"
)

func install(tl *TaskList) *Task {
	return tl.CreateTask("install").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"uv",
				"sync",
			).
				Set(func(_ context.Context, c *Command) error {
					if P.Locked {
						c.AppendArgs("--locked")

						t.Log.Info("Installing against a lockfile that has to be up to date.")
					} else {
						c.AppendArgs("--frozen")

						t.Log.Info("Installing the lockfile as it is.")
					}

					if P.Args != "" {
						c.AppendArgs(strings.Split(P.Args, " ")...)
					}

					return nil
				}).
				SetDir(setup.C.Cwd).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

// the cache keeps every wheel uv ever fetched, so the cached directory only stays the
// size of the project when the wheels it can download again are dropped.
func prune(tl *TaskList) *Task {
	return tl.CreateTask("prune").
		ShouldDisable(func(_ *Task) bool {
			return setup.C.Env["UV_CACHE_DIR"] == ""
		}).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"uv",
				"cache",
				"prune",
				"--ci",
			).
				SetDir(setup.C.Cwd).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
