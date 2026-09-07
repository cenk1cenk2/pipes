package login

import (
	"io"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
)

func ContainerRegistryLoginParent(tl *TaskList) *Task {
	return tl.CreateTask("login", "parent").
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobParallel(
				ContainerRegistryLogin(tl).Job(),
				ContainerRegistryLoginVerify(tl).Job(),
			)
		})
}

func ContainerRegistryLogin(tl *TaskList) *Task {
	return tl.CreateTask("login").
		ShouldDisable(func(t *Task) bool {
			return P.Username == "" ||
				P.Password == ""
		}).
		ShouldRunBefore(func(t *Task) error {
			t.Plumber.AppendSecrets(P.Password)

			return nil
		}).
		Set(func(t *Task) error {
			t.CreateCommand(
				"buildah",
				"login",
				P.Uri,
				"--username",
				P.Username,
				"--password-stdin",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT).
				Set(func(c *Command) error {
					c.Log.Infof(
						"Logging in to container registry: %s",
						P.Uri,
					)

					return nil
				}).
				SetStdin(func(c *Command) io.Reader {
					return strings.NewReader(P.Password)
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

// ContainerRegistryLoginVerify is the counterpart of the login task, for a
// pipeline that carries no credentials and is relying on an ambient login: it
// proves that login actually exists before the build spends time on an image it
// cannot push.
func ContainerRegistryLoginVerify(tl *TaskList) *Task {
	return tl.CreateTask("login", "verify").
		ShouldDisable(func(_ *Task) bool {
			return P.Username != "" && P.Password != ""
		}).
		Set(func(t *Task) error {
			t.CreateCommand(
				"buildah",
				"login",
				P.Uri,
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
				Set(func(c *Command) error {
					c.Log.Debugf(
						"Will verify authentication in to container registry: %s",
						P.Uri,
					)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
