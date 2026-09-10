package setup

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"

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
			t.CreateCommand("kustomize", "version").
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				CaptureOutput(&C.Version).
				ShouldRunAfter(func(_ context.Context, c *Command) error {
					c.Log.Info(fmt.Sprintf("kustomize version: %s", C.Version))

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func resolve(tl *TaskList) *Task {
	return tl.CreateTask("resolve").
		Set(func(_ context.Context, t *Task) error {
			cwd := C.Cwd
			if cwd == "" {
				cwd = "."
			}

			if len(P.Paths) > 0 {
				overlays := make([]string, 0, len(P.Paths))
				for _, path := range P.Paths {
					overlays = append(overlays, filepath.Join(cwd, path))
				}

				C.Overlays = slices.Compact(slices.Sorted(slices.Values(overlays)))
				t.Log.Debug(fmt.Sprintf("Using explicit overlay paths: %v", C.Overlays))

				return nil
			}

			C.Overlays = []string{cwd}
			t.Log.Debug(fmt.Sprintf("Using overlay path: %s", cwd))

			return nil
		})
}
