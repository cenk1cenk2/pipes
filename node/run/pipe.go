package run

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type (
	NodeCommand struct {
		Script     string
		ScriptArgs string
		Cwd        string `validate:"dir"`
	}

	Pipe struct {
		NodeCommand NodeCommand
	}
)

var P = TaskList[Pipe, Ctx]{
	Pipe:    Pipe{},
	Context: Ctx{},
}

func New(a *App) *TaskList[Pipe, Ctx] {
	return P.New(a).ShouldRunBefore(func(tl *TaskList[Pipe, Ctx], ctx *cli.Context) error {
		args := ctx.Args().Slice()
		if len(args) < 1 {
			return fmt.Errorf("Arguments are needed to run a specific script.")
		} else {
			P.Pipe.NodeCommand.Script, P.Pipe.NodeCommand.ScriptArgs = args[0], strings.Join(args[1:], " ")
		}

		return nil
	}).SetTasks(
		P.JobSequence(
			RunNodeScript(&P).Job(),
		),
	)
}
