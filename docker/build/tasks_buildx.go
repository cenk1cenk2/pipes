package build

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"gitlab.kilic.dev/devops/pipes/docker/setup"
	. "github.com/cenk1cenk2/plumber/v6"
)

func DockerBuildXParent(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("buildx", "parent").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !setup.TL.Pipe.Docker.UseBuildx
		}).
		SetJobWrapper(func(job Job, t *Task[Pipe]) Job {
			return tl.JobSequence(
				DockerBuildXCreate(tl).Job(),
				DockerBuildxSetupQemu(tl).Job(),
				DockerBuildX(tl).Job(),
			)
		})
}

func DockerBuildXCreate(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("buildx", "create").
		Set(func(t *Task[Pipe]) error {
			r, err := rand.Int(rand.Reader, big.NewInt(1000))

			if err != nil {
				return err
			}

			instance := fmt.Sprintf("%s_%d", setup.TL.Pipe.Docker.BuildxInstance, r)

			t.CreateCommand(
				setup.DOCKER_EXE,
				"buildx",
				"create",
				"--use",
				"--name",
				instance,
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				Set(func(c *Command[Pipe]) error {
					c.Log.Infoln("Creating a new instance of docker buildx.")

					return nil
				}).
				AddSelfToTheTask()

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
				"--provenance=false",
			).
				Set(func(c *Command[Pipe]) error {
					var err error
					if t.Pipe.DockerImage.BuildArgs, err = InlineTemplates[any](t.Pipe.DockerImage.BuildArgs, nil); err != nil {
						return err
					}

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
