package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func Setup(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		Set(func(t *Task[Pipe]) error {
			t.Pipe.Ctx.Tags = []string{}

			return nil
		})
}

func DockerVersion(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("version").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				DOCKER_EXE,
				"--version",
			).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobParallel()
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
