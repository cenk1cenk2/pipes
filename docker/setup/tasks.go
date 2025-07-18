package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
)

func DockerVersion(tl *TaskList) *Task {
	return tl.CreateTask("version").
		Set(func(t *Task) error {
			t.CreateCommand(
				DOCKER_EXE,
				"--version",
			).
				SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobParallel()
		})
}

func DockerBuildXVersion(tl *TaskList) *Task {
	return tl.CreateTask("version", "buildx").
		ShouldDisable(func(t *Task) bool {
			return !P.Docker.UseBuildx
		}).
		Set(func(t *Task) error {
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
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobParallel()
		})
}
