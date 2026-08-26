package lint

import (
	"fmt"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/go/setup"
)

func GoLint(tl *TaskList) *Task {
	return tl.CreateTask("lint").
		Set(func(t *Task) error {
			if P.Workspace {
				t.CreateCommand(
					"go",
					"list",
					"-m",
					"-f",
					"{{.Dir}}",
				).
					AppendEnvironment(setup.C.EnvVars).
					SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
					SetDir(setup.P.Cwd).
					EnableStreamRecording().
					ShouldRunAfter(func(c *Command) error {
						for _, module := range c.GetStdoutStream() {
							if module := strings.TrimSpace(module); module != "" {
								C.Modules = append(C.Modules, module)
							}
						}

						if len(C.Modules) == 0 {
							return fmt.Errorf("Can not resolve any modules of the go workspace.")
						}

						t.Log.Infof("Linting modules of the workspace: %s", strings.Join(C.Modules, ", "))

						return nil
					}).
					AddSelfToTheTask()
			}

			t.CreateCommand(
				"golangci-lint",
				"run",
				"-v",
				"--timeout",
				P.Timeout.String(),
			).
				Set(func(c *Command) error {
					for _, module := range C.Modules {
						c.AppendArgs(filepath.Join(module, "..."))
					}

					return nil
				}).
				AppendEnvironment(setup.C.EnvVars).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
