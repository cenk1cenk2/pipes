package pipe

import (
	"strings"

	"github.com/workanator/go-floc/v3"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type Ctx struct {
	Exe string
}

func InstallPackages(tl *TaskList[Pipe, Ctx]) *Task[Pipe, Ctx] {
	t := Task[Pipe, Ctx]{}

	return t.New(tl, "install").Set(func(t *Task[Pipe, Ctx], c floc.Control) error {
		apks := t.Pipe.Apk.Value()

		if len(apks) > 0 {
			t.Log.Debugf(
				"Will install packages from APK repository: %s",
				strings.Join(apks, ", "),
			)

			cmd := Command[Pipe, Ctx]{}

			cmd.New(t, "apk", "--no-cache").Set(func(c *Command[Pipe, Ctx]) error {
				c.AppendArgs(apks...)

				return nil
			})

			t.AddCommands(cmd)
		}

		packages := t.Pipe.Node.Value()

		if len(packages) > 0 {
			t.Log.Debugf(
				"Will install packages from NPM repository: %s",
				strings.Join(packages, ", "),
			)

			cmd := Command[Pipe, Ctx]{}

			cmd.New(t, "yarn", "global", "add").Set(func(c *Command[Pipe, Ctx]) error {
				c.AppendArgs(packages...)

				return nil
			})

			t.AddCommands(cmd)
		}

		return nil
	}).ShouldRunAfter(func(t *Task[Pipe, Ctx], c floc.Control) error {
		return t.RunCommandJobAsJobParallel()
	})
}

func RunSemanticRelease(tl *TaskList[Pipe, Ctx]) *Task[Pipe, Ctx] {
	t := Task[Pipe, Ctx]{}

	return t.New(tl, "release").Set(func(t *Task[Pipe, Ctx], c floc.Control) error {
		if t.Pipe.UseMulti {
			t.Context.Exe = MULTI_SEMANTIC_RELEASE_EXE
		} else {
			t.Context.Exe = SEMANTIC_RELEASE_EXE
		}

		cmd := Command[Pipe, Ctx]{}

		cmd.New(t, t.Context.Exe).Set(func(c *Command[Pipe, Ctx]) error {
			if t.Pipe.SemanticRelease.IsDryRun {
				c.AppendArgs("--dry-run")
			}

			if t.App.Environment.Debug {
				c.AppendArgs("--debug")
			}

			return nil
		})

		t.AddCommands(cmd)

		return nil
	}).ShouldRunAfter(func(t *Task[Pipe, Ctx], c floc.Control) error {
		return t.RunCommandJobAsJobSequence()
	})
}
