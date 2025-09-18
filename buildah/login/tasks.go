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
			return P.ContainerRegistry.Username == "" ||
				P.ContainerRegistry.Password == ""
		}).
		ShouldRunBefore(func(t *Task) error {
			t.Plumber.AppendSecrets(P.ContainerRegistry.Password)

			return nil
		}).
		Set(func(t *Task) error {

			// login task
			t.CreateCommand(
				"buildah",
				"login",
				P.ContainerRegistry.Uri,
				"--username",
				P.ContainerRegistry.Username,
				"--password-stdin",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT).
				Set(func(c *Command) error {
					c.Log.Infof(
						"Logging in to container registry: %s",
						P.ContainerRegistry.Uri,
					)

					return nil
				}).
				SetStdin(func(c *Command) io.Reader {
					return strings.NewReader(P.ContainerRegistry.Password)
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func ContainerRegistryLoginVerify(tl *TaskList) *Task {
	return tl.CreateTask("login", "verify").
		ShouldDisable(func(t *Task) bool {
			return P.ContainerRegistry.Username != "" &&
				P.ContainerRegistry.Password != ""
		}).
		Set(func(t *Task) error {
			t.CreateCommand(
				"buildah",
				"login",
				P.ContainerRegistry.Uri,
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
				Set(func(c *Command) error {
					c.Log.Debugf(
						"Will verify authentication in to container registry: %s",
						P.ContainerRegistry.Uri,
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
