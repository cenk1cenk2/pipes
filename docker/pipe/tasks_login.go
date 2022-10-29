package pipe

import (
	"strings"

	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func DockerLoginParent(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("login", "parent").
		Set(func(t *Task[Pipe]) error {
			t.SetSubtask(
				tl.JobSequence(
					DockerLogin(tl).Job(),
					DockerLoginVerify(tl).Job(),
				),
			)

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}

func DockerLogin(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("login").
		Set(func(t *Task[Pipe]) error {
			// login task
			t.CreateCommand(
				DOCKER_EXE,
				"login",
				t.Pipe.DockerRegistry.Registry,
				"--username",
				t.Pipe.DockerRegistry.Username,
				"--password-stdin",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT).
				Set(func(c *Command[Pipe]) error {
					c.Command.Stdin = strings.NewReader(t.Pipe.DockerRegistry.Password)

					c.Log.Infof(
						"Logging in to Docker registry: %s",
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

func DockerLoginVerify(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("login", "verify").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				DOCKER_EXE,
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
