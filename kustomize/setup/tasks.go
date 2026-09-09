package setup

import (
	"path/filepath"
	"slices"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

func initialize(tl *TaskList) *Task {
	return tl.CreateTask("init").
		Set(func(t *Task) error {
			C.Cwd = P.Cwd

			t.Log.Debugf("Working directory: %s", C.Cwd)

			return nil
		})
}

func version(tl *TaskList) *Task {
	return tl.CreateTask("version").
		Set(func(t *Task) error {
			t.CreateCommand("kustomize", "version").
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				ShouldRunAfter(func(c *Command) error {
					C.Version = strings.TrimSpace(strings.Join(c.GetCombinedStream(), "\n"))

					c.Log.Infof("kustomize version: %s", C.Version)

					return nil
				}).
				EnableStreamRecording().
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func resolve(tl *TaskList) *Task {
	return tl.CreateTask("resolve").
		Set(func(t *Task) error {
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
				t.Log.Debugf("Using explicit overlay paths: %v", C.Overlays)

				return nil
			}

			C.Overlays = []string{cwd}
			t.Log.Debugf("Using overlay path: %s", cwd)

			return nil
		})
}
