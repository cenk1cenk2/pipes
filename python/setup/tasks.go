package setup

import (
	"context"
	"fmt"
	"path/filepath"

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
			for _, exe := range []string{"python", "uv"} {
				var output string

				t.CreateCommand(exe, "--version").
					SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
					CaptureOutput(&output).
					ShouldRunAfter(func(_ context.Context, c *Command) error {
						c.Log.Info(fmt.Sprintf("%s version: %s", exe, output))

						return nil
					}).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobParallel(ctx)
		})
}

// uv resolves a relative cache directory against the directory it runs in, which
// differs between the project and the job root, so the cache is pinned absolute.
func env(tl *TaskList) *Task {
	return tl.CreateTask("env").
		Set(func(_ context.Context, t *Task) error {
			if P.Cache != "" {
				cache, err := filepath.Abs(P.Cache)
				if err != nil {
					return fmt.Errorf("Cannot get absolute path of cache dir: %s -> %w", P.Cache, err)
				}

				C.Env["UV_CACHE_DIR"] = cache
			}

			return nil
		})
}
