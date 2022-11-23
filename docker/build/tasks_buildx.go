package build

import (
	"gitlab.kilic.dev/devops/pipes/docker/setup"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func DockerBuildXParent(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("buildx", "parent").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !setup.TL.Pipe.Docker.UseBuildx
		}).
		SetJobWrapper(func(job Job) Job {
			return tl.JobSequence(
				DockerBuildXCreate(tl).Job(),
				DockerBuildXUse(tl).Job(),
				DockerBuildxSetupQemu(tl).Job(),
				DockerBuildX(tl).Job(),
			)
		})
}

func DockerBuildXCreate(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("buildx", "create").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				setup.DOCKER_EXE,
				"buildx",
				"create",
				"--use",
				"--name",
				setup.TL.Pipe.Docker.BuildxInstance,
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				SetIgnoreError().
				Set(func(c *Command[Pipe]) error {
					c.Log.Infoln("Creating a new instance of docker buildx.")

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			if err := t.RunCommandJobAsJobSequence(); err != nil {
				t.Pipe.Ctx.UseExistingBuildXInstance = true
			}

			return nil
		})
}

func DockerBuildXUse(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("buildx", "use").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.Ctx.UseExistingBuildXInstance
		}).
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				setup.DOCKER_EXE,
				"buildx",
				"use",
				setup.TL.Pipe.Docker.BuildxInstance,
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				Set(func(c *Command[Pipe]) error {
					c.Log.Warnln(
						"Creating a new docker buildx instance failed, trying to use the existing one.",
					)

					return nil
				}).AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func DockerBuildxSetupQemu(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("buildx", "qemu").
		Set(func(t *Task[Pipe]) error {
			// spawn virtual machine
			t.CreateCommand(
				setup.DOCKER_EXE,
				"run",
				"--rm",
				"--privileged",
				"multiarch/qemu-user-static",
				"--reset",
				"-p",
				"yes",
			).
				SetIgnoreError().
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT).
				AddSelfToTheTask()

			// check virtual machine
			t.CreateCommand(
				setup.DOCKER_EXE,
				"buildx",
				"inspect",
				"--bootstrap",
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func DockerBuildX(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("buildx").
		Set(func(t *Task[Pipe]) error {
			t.Log.Infof(
				"Building Docker image: %s in %s",
				t.Pipe.DockerFile.Name,
				t.Pipe.DockerFile.Context,
			)

			t.Log.Infoln("Using Docker Buildx for building the Docker image.")

			// build image
			t.CreateCommand(
				setup.DOCKER_EXE,
				"buildx",
				"build",
			).
				Set(func(c *Command[Pipe]) error {
					for _, v := range t.Pipe.DockerImage.BuildArgs {
						c.AppendArgs("--build-arg", v)
					}

					if t.Pipe.DockerImage.Pull {
						c.AppendArgs("--pull")
					}

					c.AppendArgs("--push")

					if t.Pipe.Docker.BuildxPlatforms != "" {
						c.AppendArgs("--platform", t.Pipe.Docker.BuildxPlatforms)
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
