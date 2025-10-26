package setup

import (
	"fmt"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v6"
)

func GoVersion(tl *TaskList) *Task {
	return tl.CreateTask("version").
		Set(func(t *Task) error {
			t.CreateCommand(
				"go",
				"version",
			).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func GoEnv(tl *TaskList) *Task {
	return tl.CreateTask("env").
		Set(func(t *Task) error {
			if P.Cache != "" {
				cache, err := filepath.Abs(P.Cache)
				if err != nil {
					return fmt.Errorf("Cannot get absolute path of cache dir: %s -> %w", P.Cache, err)
				}

				C.EnvVars["GOPATH"] = cache
				C.EnvVars["GOCACHE"] = filepath.Join(cache, "cache")
				C.EnvVars["GOLANGCI_LINT_CACHE"] = filepath.Join(C.EnvVars["GOCACHE"], "golangci-lint")
			}

			return nil
		})
}
