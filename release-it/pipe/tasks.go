package pipe

import (
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func RunReleaseIt(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("release").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				RELEASE_IT_EXE,
				"--ci",
			).
				Set(func(c *Command[Pipe]) error {
					if t.Plumber.Environment.Debug {
						c.AppendArgs("--verbose")
					}

					if t.Pipe.ReleaseIt.IsDryRun {
						c.AppendArgs("--dry-run")
					}

					if t.Pipe.ReleaseIt.ConfigFile != "" {
						c.AppendArgs("--config", t.Pipe.ReleaseIt.ConfigFile)
					}

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
