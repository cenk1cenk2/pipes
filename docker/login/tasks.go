package login

import (
	"io"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
)

func DockerLoginParent(tl *TaskList) *Task {
	return tl.CreateTask("login", "parent").
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobParallel(
				DockerLogin(tl).Job(),
				DockerLoginVerify(tl).Job(),
			)
		})
}

func DockerLogin(tl *TaskList) *Task {
	return tl.CreateTask("login").
		ShouldDisable(func(t *Task) bool {
			return P.DockerRegistry.Username == "" ||
				P.DockerRegistry.Password == ""
		}).
		Set(func(t *Task) error {
			t.Plumber.AppendSecrets(P.DockerRegistry.Password)

			// login task
			t.CreateCommand(
				setup.DOCKER_EXE,
				"login",
				P.DockerRegistry.Registry,
				"--username",
				P.DockerRegistry.Username,
				"--password-stdin",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT).
				Set(func(c *Command) error {
					c.Log.Infof(
						"Logging in to Docker registry: %s",
						P.DockerRegistry.Registry,
					)

					return nil
				}).
				SetStdin(func(c *Command) io.Reader {
					return strings.NewReader(P.DockerRegistry.Password)
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func DockerLoginVerify(tl *TaskList) *Task {
	return tl.CreateTask("login", "verify").
		ShouldDisable(func(t *Task) bool {
			return P.DockerRegistry.Username != "" &&
				P.DockerRegistry.Password != ""
		}).
		Set(func(t *Task) error {
			t.CreateCommand(
				setup.DOCKER_EXE,
				"login",
				P.DockerRegistry.Registry,
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
				Set(func(c *Command) error {
					c.Log.Debugf(
						"Will verify authentication in to Docker registry: %s",
						P.DockerRegistry.Registry,
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
