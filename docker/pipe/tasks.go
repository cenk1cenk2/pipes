package pipe

import (
	"strings"

	. "gitlab.kilic.dev/libraries/plumber/v3"
)

type Ctx struct {
	Tags                           []string
	TryToUseExistingBuildXInstance bool
}

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		ShouldRunBefore(func(t *Task[Pipe]) error {

			return nil
		}).
		Set(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.Tags = []string{}

			return nil
		})
}

func DockerVersion(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("version").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(DOCKER_EXE, "--version").
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobParallel()
		})
}

func DockerBuildXVersion(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("version:buildx").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.Docker.UseBuildx
		}).
		Set(func(t *Task[Pipe]) error {
			t.Log.Infoln("Docker Buildx is enabled.")

			t.CreateCommand(DOCKER_EXE, "buildx", "version").
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobParallel()
		})
}

func DockerLogin(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("login").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.DockerRegistry.Username == "" ||
				t.Pipe.DockerRegistry.Password == ""
		}).
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
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
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
		}).ShouldRunAfter(func(t *Task[Pipe]) error {
		return t.RunCommandJobAsJobParallel()
	})
}

func DockerLoginVerify(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("login:verify").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.DockerRegistry.Username != "" &&
				t.Pipe.DockerRegistry.Password != ""
		}).
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

func DockerInspect(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("inspect").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.DockerImage.Inspect
		}).
		Set(func(t *Task[Pipe]) error {
			for _, tag := range t.Pipe.Ctx.Tags {
				func(tag string) {
					t.CreateCommand(
						DOCKER_EXE,
						"manifest",
						"inspect",
						tag,
					).
						SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
						Set(func(c *Command[Pipe]) error {
							c.Log.Infof(
								"Inspecting Docker image: %s",
								tag,
							)

							return nil
						}).
						AddSelfToTheTask()
				}(tag)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobParallel()
		})
}
