package run

import (
	"gitlab.kilic.dev/devops/pipes/common/utils"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func RunNodeScript(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("run").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				setup.TL.Pipe.Ctx.PackageManager.Exe,
			).
				Set(func(c *Command[Pipe]) error {
					ctx := environment.EnvironmentTemplate{
						Environment: environment.TL.Pipe.Ctx.Environment,
						EnvVars:     environment.TL.Pipe.Ctx.EnvVars,
					}

					c.AppendArgs(setup.TL.Pipe.Ctx.PackageManager.Commands.Run...)

					if t.Pipe.NodeCommand.Script != "" {
						tmpl, err := utils.InlineTemplate(t.Pipe.NodeCommand.Script, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.AppendArgs(setup.TL.Pipe.Ctx.PackageManager.Commands.RunDelimitter...)

					if t.Pipe.NodeCommand.ScriptArgs != "" {
						tmpl, err := utils.InlineTemplate(t.Pipe.NodeCommand.ScriptArgs, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.SetDir(t.Pipe.NodeCommand.Cwd)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
