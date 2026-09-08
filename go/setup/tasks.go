package setup

import (
	"fmt"
	"path/filepath"
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
			t.CreateCommand("go", "version").
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				ShouldRunAfter(func(c *Command) error {
					C.Version = strings.TrimSpace(strings.Join(c.GetCombinedStream(), "\n"))

					c.Log.Infof("go version: %s", C.Version)

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

func env(tl *TaskList) *Task {
	return tl.CreateTask("env").
		Set(func(t *Task) error {
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

func workspace(tl *TaskList) *Task {
	return tl.CreateTask("workspace").
		Set(func(t *Task) error {
			C.Workspace = P.Workspace

			if C.Workspace {
				t.Log.Debugf("Go workspace mode is enabled: %s", C.Cwd)

				return nil
			}

			t.CreateCommand(
				"go",
				"env",
				"GOWORK",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				SetDir(C.Cwd).
				EnableStreamRecording().
				ShouldRunAfter(func(c *Command) error {
					stream := c.GetStdoutStream()

					if len(stream) == 0 {
						return nil
					}

					// go env reports "off" instead of an empty value when workspace mode is explicitly disabled.
					if gowork := strings.TrimSpace(stream[0]); gowork != "" && gowork != "off" {
						C.Workspace = true

						t.Log.Debugf("Go workspace detected: %s", gowork)
					}

					return nil
				}).
				AppendEnvironment(C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func modules(tl *TaskList) *Task {
	return tl.CreateTask("modules").
		ShouldDisable(func(_ *Task) bool {
			return !C.Workspace
		}).
		Set(func(t *Task) error {
			t.CreateCommand(
				"go",
				"list",
				"-m",
				"-f",
				"{{.Dir}}",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				SetDir(C.Cwd).
				EnableStreamRecording().
				ShouldRunAfter(func(c *Command) error {
					C.Modules = nil

					for _, module := range c.GetStdoutStream() {
						if module := strings.TrimSpace(module); module != "" {
							C.Modules = append(C.Modules, module)
						}
					}

					if len(C.Modules) == 0 {
						return fmt.Errorf("Can not resolve any modules of the go workspace.")
					}

					t.Log.Infof("Modules of the workspace: %s", strings.Join(C.Modules, ", "))

					return nil
				}).
				AppendEnvironment(C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
