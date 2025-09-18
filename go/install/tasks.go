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
				"mod",
				"vendor",
			).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
				SetDir(P.Cwd).
				Set(func(c *Command) error {
					c.Log.Info("Vendoring Go modules.")

					if P.Args != "" {
						c.AppendArgs(strings.Split(P.Args, " ")...)
					}

					if setup.P.Cache != "" {
						c.AppendEnvironment(map[string]string{
							"GOPATH": setup.P.Cache,
						})
					}

					return nil
				}).
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
				SetDir(P.Cwd).
				Set(func(c *Command) error {
					c.Log.Info("Verifying Go modules.")

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
