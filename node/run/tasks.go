package run

import (
	"strings"

	"github.com/workanator/go-floc/v3"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type Ctx struct {
}

func RunNodeScript(tl *TaskList[Pipe, Ctx]) *Task[Pipe, Ctx] {
	t := Task[Pipe, Ctx]{}

	return t.New(tl, "run").Set(func(t *Task[Pipe, Ctx], c floc.Control) error {
		cmd := Command[Pipe, Ctx]{}

		cmd.New(t, setup.P.Context.PackageManager.Exe).Set(func(c *Command[Pipe, Ctx]) error {
			c.AppendArgs(setup.P.Context.PackageManager.Commands.Run...).AppendArgs(t.Pipe.NodeCommand.Script).
				AppendArgs(setup.P.Context.PackageManager.Commands.RunDelimitter...).
				AppendArgs(strings.Split(t.Pipe.NodeCommand.ScriptArgs, " ")...)

			c.SetDir(t.Pipe.NodeCommand.Cwd)

			return nil
		})

		t.AddCommands(cmd)

		return nil
	}).ShouldRunAfter(func(t *Task[Pipe, Ctx], c floc.Control) error {
		return t.TaskList.RunJobs(t.GetCommandJobAsJobSequence())
	})
}
