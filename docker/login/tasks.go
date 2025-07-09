package login

import (
	"io"
	"strings"

	"gitlab.kilic.dev/devops/pipes/docker/setup"
	. "github.com/cenk1cenk2/plumber/v6"
)

func DockerLoginParent(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("login", "parent").
		SetJobWrapper(func(job Job, t *Task[Pipe]) Job {
			return tl.JobParallel(
				DockerLogin(tl).Job(),
				DockerLoginVerify(tl).Job(),
			)
		})
}

func DockerLogin(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("login").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.DockerRegistry.Username == "" ||
				t.Pipe.DockerRegistry.Password == ""
		}).
		Set(func(t *Task[Pipe]) error {
			t.Plumber.AppendSecrets(t.Pipe.DockerRegistry.Password)

			// login task
			t.CreateCommand(
				setup.DOCKER_EXE,
				"login",
				t.Pipe.DockerRegistry.Registry,
				"--username",
				t.Pipe.DockerRegistry.Username,
				"--password-stdin",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT).
				Set(func(c *Command[Pipe]) error {
					c.Log.Infof(
						"Logging in to Docker registry: %s",
						t.Pipe.DockerRegistry.Registry,
					)

					return nil
				}).
				SetStdin(func(c *Command[Pipe]) io.Reader {
					return strings.NewReader(t.Pipe.DockerRegistry.Password)
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func DockerLoginVerify(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("login", "verify").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.DockerRegistry.Username != "" &&
				t.Pipe.DockerRegistry.Password != ""
		}).
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				setup.DOCKER_EXE,
				"login",
				t.Pipe.DockerRegistry.Registry,
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
				Set(func(c *Command[Pipe]) error {
					c.Log.Debugf(
						"Will verify authentication in to Docker registry: %s",
						t.Pipe.DockerRegistry.Registry,
					)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
