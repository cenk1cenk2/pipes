package build

import (
	"crypto/rand"
	"fmt"
	"math/big"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
)

func DockerBuildXParent(tl *TaskList) *Task {
	return tl.CreateTask("buildx", "parent").
		ShouldDisable(func(t *Task) bool {
			return !setup.P.Docker.UseBuildx
		}).
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobSequence(
				DockerBuildXCreate(tl).Job(),
				DockerBuildxSetupQemu(tl).Job(),
				DockerBuildX(tl).Job(),
			)
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func DockerBuildXCreate(tl *TaskList) *Task {
	return tl.CreateTask("buildx", "create").
		Set(func(t *Task) error {
			r, err := rand.Int(rand.Reader, big.NewInt(1000))

			if err != nil {
				return err
			}

			instance := fmt.Sprintf("%s_%d", setup.P.Docker.BuildxInstance, r)

			t.CreateCommand(
				setup.DOCKER_EXE,
				"buildx",
				"create",
				"--use",
				"--name",
				instance,
			).
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				Set(func(c *Command) error {
					c.Log.Infoln("Creating a new instance of docker buildx.")

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func DockerBuildxSetupQemu(tl *TaskList) *Task {
	return tl.CreateTask("buildx", "qemu").
		Set(func(t *Task) error {
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
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func DockerBuildX(tl *TaskList) *Task {
	return tl.CreateTask("buildx").
		Set(func(t *Task) error {
			t.Log.Infof(
				"Building Docker image: %s in %s",
				P.DockerFile.Name,
				P.DockerFile.Context,
			)

			t.Log.Infoln("Using Docker Buildx for building the Docker image.")

			// build image
			t.CreateCommand(
				setup.DOCKER_EXE,
				"buildx",
				"build",
				"--provenance=false",
			).
				Set(func(c *Command) error {
					var err error
					if P.DockerImage.BuildArgs, err = InlineTemplates[any](P.DockerImage.BuildArgs, nil); err != nil {
						return err
					}

					for _, v := range P.DockerImage.BuildArgs {
						c.AppendArgs("--build-arg", v)
					}

					if P.DockerImage.Pull {
						c.AppendArgs("--pull")
					}

					c.AppendArgs("--push")

					if P.Docker.BuildxPlatforms != "" {
						c.AppendArgs("--platform", P.Docker.BuildxPlatforms)
					}

					for _, tag := range C.Tags {
						c.AppendArgs("-t", tag)
					}

					c.AppendArgs(
						"--file",
						P.DockerFile.Name,
						".",
					)

					c.SetDir(P.DockerFile.Context)
					t.Log.Debugf("CWD set as: %s", c.Command.Dir)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
