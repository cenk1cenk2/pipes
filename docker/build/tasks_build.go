package build

import (
	"time"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
)

func DockerBuildParent(tl *TaskList) *Task {
	return tl.CreateTask("build", "parent").
		ShouldDisable(func(t *Task) bool {
			return setup.P.Docker.UseBuildx
		}).
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobSequence(
				DockerBuild(tl).Job(),
				DockerPush(tl).Job(),
			)
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func DockerBuild(tl *TaskList) *Task {
	return tl.CreateTask("build").
		Set(func(t *Task) error {
			t.Log.Infof(
				"Building Docker image: %s in %s",
				P.DockerFile.Name,
				P.DockerFile.Context,
			)

			// build image
			t.CreateCommand(
				setup.DOCKER_EXE,
				"build",
			).
				Set(func(c *Command) error {
					if setup.P.Docker.UseBuildKit {
						t.Log.Infoln("Using Docker BuildKit for the build operation.")

						c.AppendEnvironment(
							map[string]string{
								"DOCKER_BUILDKIT": "1",
							},
						)
					} else {
						t.Log.Infoln("Forcing Docker to use legacy build mode.")

						c.AppendEnvironment(
							map[string]string{
								"DOCKER_BUILDKIT": "0",
							},
						)
					}

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

					if setup.P.Docker.UseBuildKit {
						c.AppendArgs("--push")
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

func DockerPush(tl *TaskList) *Task {
	return tl.CreateTask("push").
		ShouldDisable(func(t *Task) bool {
			return setup.P.Docker.UseBuildKit
		}).
		Set(func(t *Task) error {
			for _, tag := range C.Tags {
				t.CreateCommand(
					setup.DOCKER_EXE,
					"push",
					tag,
				).
					Set(func(c *Command) error {
						c.Log.Infof(
							"Pushing Docker image: %s",
							tag,
						)

						return nil
					}).
					SetRetries(&CommandRetry{
						Tries: 3,
						Delay: time.Second * 10,
					}).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobParallel()
		})
}
