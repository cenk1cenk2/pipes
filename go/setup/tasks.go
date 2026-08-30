package setup

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cenk1cenk2/plumber/v6"
)

func GoEnv(tl *plumber.TaskList) *plumber.Task {
	return tl.CreateTask("env").
		Set(func(t *plumber.Task) error {
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

func GoWorkspace(tl *plumber.TaskList) *plumber.Task {
	return tl.CreateTask("workspace").
		Set(func(t *plumber.Task) error {
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
				SetLogLevel(plumber.LOG_LEVEL_DEBUG, plumber.LOG_LEVEL_DEBUG, plumber.LOG_LEVEL_DEBUG).
				SetDir(C.Cwd).
				EnableStreamRecording().
				ShouldRunAfter(func(c *plumber.Command) error {
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
		ShouldRunAfter(func(t *plumber.Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
