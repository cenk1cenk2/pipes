package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

func DockerBuildParent(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("build:parent").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Docker.UseBuildx
		}).
		Set(func(t *Task[Pipe]) error {
			t.SetSubtask(
				tl.JobSequence(
					DockerBuild(tl).Job(),
					DockerPush(tl).Job(),
				),
			)

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}

func DockerBuild(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("build").
		Set(func(t *Task[Pipe]) error {
			t.Log.Infof(
				"Building Docker image: %s in %s",
				t.Pipe.DockerFile.Name,
				t.Pipe.DockerFile.Context,
			)

			// build image
			t.CreateCommand(
				DOCKER_EXE,
				"build",
			).
				Set(func(c *Command[Pipe]) error {
					for _, v := range t.Pipe.DockerImage.BuildArgs.Value() {
						c.AppendArgs("--build-arg", v)
					}

					if t.Pipe.DockerImage.Pull {
						c.AppendArgs("--pull")
					}

					for _, tag := range t.Pipe.Ctx.Tags {
						c.AppendArgs("-t", tag)
					}

					c.AppendArgs(
						"--file",
						t.Pipe.DockerFile.Name,
						".",
					)

					c.SetDir(t.Pipe.DockerFile.Context)
					t.Log.Debugf("CWD set as: %s", c.Command.Dir)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func DockerPush(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("push").
		Set(func(t *Task[Pipe]) error {
			for _, tag := range t.Pipe.Ctx.Tags {
				func(tag string) {
					t.CreateCommand(
						DOCKER_EXE,
						"push",
						tag,
					).
						Set(func(c *Command[Pipe]) error {
							c.Log.Infof(
								"Pushing Docker image: %s",
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
