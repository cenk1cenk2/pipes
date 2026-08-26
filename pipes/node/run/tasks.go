package run

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/pipes/node/setup"
)

func RunNodeScript(tl *TaskList) *Task {
	return tl.CreateTask("run", C.Script).
		Set(func(t *Task) error {
			t.CreateCommand(
				setup.NodeCtx.PackageManager.Exe,
			).
				Set(func(c *Command) error {
					ctx := environment.Template{
						Environment: setup.EnvironmentCtx.Environment,
						EnvVars:     setup.EnvironmentCtx.EnvVars,
					}

					c.AppendArgs(setup.NodeCtx.PackageManager.Commands.Run...)

					if C.Script != "" {
						tmpl, err := InlineTemplate(C.Script, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.AppendArgs(setup.NodeCtx.PackageManager.Commands.RunDelimiter...)

					if C.ScriptArgs != "" {
						tmpl, err := InlineTemplate(C.ScriptArgs, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.SetDir(P.NodeCommand.Cwd)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
