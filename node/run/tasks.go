package run

import (
	"strings"

	"github.com/workanator/go-floc/v3"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type Ctx struct {
}

func RunNodeScript(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("run").
		Set(func(t *Task[Pipe], c floc.Control) error {
			t.CreateCommand(setup.P.Pipe.Ctx.PackageManager.Exe).Set(func(c *Command[Pipe]) error {
				c.AppendArgs(setup.P.Pipe.Ctx.PackageManager.Commands.Run...).
					AppendArgs(t.Pipe.NodeCommand.Script).
					AppendArgs(setup.P.Pipe.Ctx.PackageManager.Commands.RunDelimitter...).
					AppendArgs(strings.Split(t.Pipe.NodeCommand.ScriptArgs, " ")...)

				c.SetDir(t.Pipe.NodeCommand.Cwd)

				return nil
			}).AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe], c floc.Control) error {
			return t.RunCommandJobAsJobSequence()
		})
}
