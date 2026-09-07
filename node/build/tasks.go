package build

import (
	"os"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/node/setup"
)

func BuildNodeApplication(tl *TaskList) *Task {
	return tl.CreateTask("build").
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

					if P.Build.Script != "" {
						tmpl, err := InlineTemplate(P.Build.Script, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.AppendArgs(setup.NodeCtx.PackageManager.Commands.RunDelimiter...)

					if P.Build.ScriptArgs != "" {
						tmpl, err := InlineTemplate(P.Build.ScriptArgs, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.SetDir(P.Build.Cwd)

					c.AppendDirectEnvironment(os.Environ()...).
						AppendEnvironment(setup.EnvironmentCtx.EnvVars)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
