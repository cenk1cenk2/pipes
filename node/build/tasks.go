package build

import (
	"os"

	"gitlab.kilic.dev/devops/pipes/node/setup"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func BuildNodeApplication(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("build").
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

					if t.Pipe.NodeBuild.Script != "" {
						tmpl, err := InlineTemplate(t.Pipe.NodeBuild.Script, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.AppendArgs(setup.TL.Pipe.Ctx.PackageManager.Commands.RunDelimitter...)

					if t.Pipe.NodeBuild.ScriptArgs != "" {
						tmpl, err := InlineTemplate(t.Pipe.NodeBuild.ScriptArgs, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.SetDir(t.Pipe.NodeBuild.Cwd)

					c.AppendDirectEnvironment(os.Environ()...).
						AppendEnvironment(environment.TL.Pipe.Ctx.EnvVars)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
