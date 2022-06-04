package run

import (
	"strings"

	"gitlab.kilic.dev/devops/pipes/node/setup"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

type Ctx struct {
}

func RunNodeScript(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("run").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(setup.TL.Pipe.Ctx.PackageManager.Exe).Set(func(c *Command[Pipe]) error {
				c.AppendArgs(setup.TL.Pipe.Ctx.PackageManager.Commands.Run...).
					AppendArgs(t.Pipe.NodeCommand.Script).
					AppendArgs(setup.TL.Pipe.Ctx.PackageManager.Commands.RunDelimitter...).
					AppendArgs(strings.Split(t.Pipe.NodeCommand.ScriptArgs, " ")...)

				c.SetDir(t.Pipe.NodeCommand.Cwd)

				return nil
			}).AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
