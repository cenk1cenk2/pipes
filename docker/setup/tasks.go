package setup

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

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

func DockerBuildXVersion(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("version", "buildx").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return !t.Pipe.Docker.UseBuildx
		}).
		Set(func(t *Task[Pipe]) error {
			t.Log.Infoln("Docker Buildx is enabled.")

			t.CreateCommand(
				DOCKER_EXE,
				"buildx",
				"version",
			).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobParallel()
		})
}
