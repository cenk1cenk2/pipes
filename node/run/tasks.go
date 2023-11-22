package run

import (
	"gitlab.kilic.dev/devops/pipes/node/setup"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func RunNodeScript(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("run", tl.Pipe.Ctx.Script).
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

					if t.Pipe.Ctx.Script != "" {
						tmpl, err := InlineTemplate(t.Pipe.Ctx.Script, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.AppendArgs(setup.TL.Pipe.Ctx.PackageManager.Commands.RunDelimitter...)

					if t.Pipe.Ctx.ScriptArgs != "" {
						tmpl, err := InlineTemplate(t.Pipe.Ctx.ScriptArgs, ctx)

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
