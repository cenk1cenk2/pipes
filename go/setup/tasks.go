package setup

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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
			t.CreateCommand("go", "version").
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				ShouldRunAfter(func(_ context.Context, c *Command) error {
					C.Version = strings.TrimSpace(strings.Join(c.GetCombinedStream(), "\n"))

					c.Log.Info(fmt.Sprintf("go version: %s", C.Version))

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

func env(tl *TaskList) *Task {
	return tl.CreateTask("env").
		Set(func(_ context.Context, t *Task) error {
			if P.Cache != "" {
				cache, err := filepath.Abs(P.Cache)
				if err != nil {
					return fmt.Errorf("Cannot get absolute path of cache dir: %s -> %w", P.Cache, err)
				}

				C.Env["GOPATH"] = cache
				C.Env["GOCACHE"] = filepath.Join(cache, "go-build")
				C.Env["GOLANGCI_LINT_CACHE"] = filepath.Join(cache, "golangci-lint")
			}

			return nil
		})
}

// The detection is scoped to the working directory on purpose: "go env GOWORK"
// searches the parent directories too, which would turn every module that
// merely sits inside a workspace into a workspace run.
func workspace(tl *TaskList) *Task {
	return tl.CreateTask("workspace").
		Set(func(_ context.Context, t *Task) error {
			gowork := filepath.Join(C.Cwd, "go.work")

			if _, err := os.Stat(gowork); err != nil {
				if !errors.Is(err, fs.ErrNotExist) {
					return fmt.Errorf("Cannot probe for the workspace file: %s -> %w", gowork, err)
				}

				C.Workspace = false

				return nil
			}

			C.Workspace = true

			t.Log.Debug(fmt.Sprintf("Go workspace detected: %s", gowork))

			return nil
		})
}

func modules(tl *TaskList) *Task {
	return tl.CreateTask("modules").
		ShouldDisable(func(_ *Task) bool {
			return !C.Workspace
		}).
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"go",
				"list",
				"-m",
				"-f",
				"{{.Dir}}",
			).
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				SetDir(C.Cwd).
				EnableStreamRecording().
				ShouldRunAfter(func(_ context.Context, c *Command) error {
					C.Modules = nil

					for _, module := range c.GetStdoutStream() {
						if module := strings.TrimSpace(module); module != "" {
							C.Modules = append(C.Modules, module)
						}
					}

					if len(C.Modules) == 0 {
						return fmt.Errorf("Can not resolve any modules of the go workspace.")
					}

					t.Log.Info(fmt.Sprintf("Modules of the workspace: %s", strings.Join(C.Modules, ", ")))

					return nil
				}).
				AppendEnvironment(C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
