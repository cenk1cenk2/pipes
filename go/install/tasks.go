package install

import (
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/go/setup"
)

func GoModVendor(tl *TaskList) *Task {
	return tl.CreateTask("vendor").
		Set(func(t *Task) error {
			t.CreateCommand(
				"go",
				"env",
				"GOWORK",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				SetDir(setup.C.Cwd).
				EnableStreamRecording().
				ShouldRunAfter(func(c *Command) error {
					stream := c.GetStdoutStream()

					if len(stream) == 0 {
						return nil
					}

					// go env reports "off" instead of an empty value when workspace mode is explicitly disabled.
					if workspace := strings.TrimSpace(stream[0]); workspace != "off" {
						C.Workspace = workspace
					}

					return nil
				}).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			t.CreateCommand(
				"go",
			).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
				SetDir(setup.C.Cwd).
				Set(func(c *Command) error {
					if C.Workspace != "" {
						c.AppendArgs("work", "vendor")

						t.Log.Infof("Vendoring workspace: %s in %s", C.Workspace, setup.C.Cwd)
					} else {
						c.AppendArgs("mod", "vendor")

						t.Log.Infof("Vendoring: in %s", setup.C.Cwd)
					}

					if P.Args != "" {
						c.AppendArgs(strings.Split(P.Args, " ")...)
					}

					return nil
				}).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func GoModVerify(tl *TaskList) *Task {
	return tl.CreateTask("verify").
		ShouldDisable(func(t *Task) bool {
			return !P.Verify
		}).
		Set(func(t *Task) error {
			t.CreateCommand(
				"go",
				"mod",
				"verify",
			).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
				SetDir(setup.C.Cwd).
				Set(func(c *Command) error {
					t.Log.Infof("Verifying modules: in %s", setup.C.Cwd)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
