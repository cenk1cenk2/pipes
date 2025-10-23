package login

import (
	"io"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

func HelmLogin(tl *TaskList) *Task {
	return tl.CreateTask("login").
		ShouldDisable(func(t *Task) bool {
			return P.HelmRegistry.Username == "" ||
				P.HelmRegistry.Password == ""
		}).
		ShouldRunBefore(func(t *Task) error {
			t.Plumber.AppendSecrets(P.HelmRegistry.Password)

			return nil
		}).
		Set(func(t *Task) error {
			t.CreateCommand(
				"helm",
				"registry",
				"login",
				P.HelmRegistry.Uri,
				"--username",
				P.HelmRegistry.Username,
				"--password-stdin",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT).
				Set(func(c *Command) error {
					c.Log.Infof(
						"Logging in to chart repository: %s",
						P.HelmRegistry.Uri,
					)

					return nil
				}).
				SetStdin(func(c *Command) io.Reader {
					return strings.NewReader(P.HelmRegistry.Password)
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
