package pipe

import (
	"gitlab.kilic.dev/devops/pipes/common/utils"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func AddNodeModules(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("packages", "node").
		Set(func(t *Task[Pipe]) error {
			t.CreateCommand(
				setup.TL.Pipe.Ctx.PackageManager.Exe,
			).
				Set(func(c *Command[Pipe]) error {
					ctx := environment.EnvironmentTemplate{
						Environment: environment.TL.Pipe.Ctx.Environment,
						EnvVars:     environment.TL.Pipe.Ctx.EnvVars,
					}

					if t.Pipe.NodeAdd.Global {
						c.AppendArgs(setup.TL.Pipe.Ctx.PackageManager.Commands.Global...)
					}

					c.AppendArgs(setup.TL.Pipe.Ctx.PackageManager.Commands.Add...)

					if t.Pipe.NodeAdd.ScriptArgs != "" {
						tmpl, err := utils.InlineTemplate(t.Pipe.NodeAdd.ScriptArgs, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.AppendArgs(t.Pipe.NodeAdd.Packages...)

					c.SetDir(t.Pipe.NodeAdd.Cwd)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
