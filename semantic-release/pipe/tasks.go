package pipe

import (
	"strings"

	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func InstallApkPackages(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("apks").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(t.Pipe.Apk) == 0
		}).
		Set(func(t *Task[Pipe]) error {
			apks := t.Pipe.Apk

			t.Log.Debugf(
				"Will install packages from APK repository: %s",
				strings.Join(apks, ", "),
			)

			t.CreateCommand(
				"apk",
				"--no-cache",
			).
				Set(func(c *Command[Pipe]) error {
					c.AppendArgs(apks...)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobParallel()
		})
}

func RunSemanticRelease(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("release").
		Set(func(t *Task[Pipe]) error {
			if t.Pipe.Workspace {
				t.Pipe.Ctx.Exe = MULTI_SEMANTIC_RELEASE_EXE
			} else {
				t.Pipe.Ctx.Exe = SEMANTIC_RELEASE_EXE
			}

			t.CreateCommand(
				t.Pipe.Ctx.Exe,
			).
				Set(func(c *Command[Pipe]) error {
					if t.Pipe.SemanticRelease.IsDryRun {
						c.AppendArgs("--dry-run")
					}

					if t.Plumber.Environment.Debug {
						c.AppendArgs("--debug")
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
