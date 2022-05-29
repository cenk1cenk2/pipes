package pipe

import (
	"strings"

	"github.com/workanator/go-floc/v3"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type Ctx struct {
	Exe string
}

func InstallApkPackages(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("apks").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(t.Pipe.Apk.Value()) == 0
		}).
		Set(func(t *Task[Pipe], c floc.Control) error {
			apks := t.Pipe.Apk.Value()

			t.Log.Debugf(
				"Will install packages from APK repository: %s",
				strings.Join(apks, ", "),
			)

			t.CreateCommand("apk", "--no-cache").Set(func(c *Command[Pipe]) error {
				c.AppendArgs(apks...)

				return nil
			}).AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe], c floc.Control) error {
			return t.RunCommandJobAsJobParallel()
		})
}

func InstallNodePackages(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("install").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return len(t.Pipe.Node.Value()) == 0
		}).
		Set(func(t *Task[Pipe], c floc.Control) error {
			packages := t.Pipe.Node.Value()

			t.Log.Debugf(
				"Will install packages from NPM repository: %s",
				strings.Join(packages, ", "),
			)

			t.CreateCommand("yarn", "global", "add").Set(func(c *Command[Pipe]) error {
				c.AppendArgs(packages...)

				return nil
			}).AddSelfToTheTask()

			return nil
		}).ShouldRunAfter(func(t *Task[Pipe], c floc.Control) error {
		return t.RunCommandJobAsJobParallel()
	})
}

func RunSemanticRelease(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("release").
		Set(func(t *Task[Pipe], c floc.Control) error {
			if t.Pipe.UseMulti {
				t.Pipe.Ctx.Exe = MULTI_SEMANTIC_RELEASE_EXE
			} else {
				t.Pipe.Ctx.Exe = SEMANTIC_RELEASE_EXE
			}

			t.CreateCommand(t.Pipe.Ctx.Exe).Set(func(c *Command[Pipe]) error {
				if t.Pipe.SemanticRelease.IsDryRun {
					c.AppendArgs("--dry-run")
				}

				if t.App.Environment.Debug {
					c.AppendArgs("--debug")
				}

				return nil
			}).AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe], c floc.Control) error {
			return t.RunCommandJobAsJobSequence()
		})
}
