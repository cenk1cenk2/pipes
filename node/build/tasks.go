package build

import (
	"os"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
)

func BuildNodeApplication(tl *TaskList) *Task {
	return tl.CreateTask("build").
		Set(func(t *Task) error {
			t.CreateCommand(
				setup.C.PackageManager.Exe,
			).
				Set(func(c *Command) error {
					ctx := environment.EnvironmentTemplate{
						Environment: environment.C.Environment,
						EnvVars:     environment.C.EnvVars,
					}

					c.AppendArgs(setup.C.PackageManager.Commands.Run...)

					if P.NodeBuild.Script != "" {
						tmpl, err := InlineTemplate(P.NodeBuild.Script, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.AppendArgs(setup.C.PackageManager.Commands.RunDelimiter...)

					if P.NodeBuild.ScriptArgs != "" {
						tmpl, err := InlineTemplate(P.NodeBuild.ScriptArgs, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.SetDir(P.NodeBuild.Cwd)

					c.AppendDirectEnvironment(os.Environ()...).
						AppendEnvironment(environment.C.EnvVars)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
