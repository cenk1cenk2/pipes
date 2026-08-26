package build

import (
	"os"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
)

func BuildNodeApplication(tl *TaskList, deps Deps) *Task {
	return tl.CreateTask("build").
		Set(func(t *Task) error {
			t.CreateCommand(
				deps.Node.PackageManager.Exe,
			).
				Set(func(c *Command) error {
					ctx := environment.Template{
						Environment: deps.Environment.Environment,
						EnvVars:     deps.Environment.EnvVars,
					}

					c.AppendArgs(deps.Node.PackageManager.Commands.Run...)

					if P.NodeBuild.Script != "" {
						tmpl, err := InlineTemplate(P.NodeBuild.Script, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.AppendArgs(deps.Node.PackageManager.Commands.RunDelimiter...)

					if P.NodeBuild.ScriptArgs != "" {
						tmpl, err := InlineTemplate(P.NodeBuild.ScriptArgs, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.SetDir(P.NodeBuild.Cwd)

					c.AppendDirectEnvironment(os.Environ()...).
						AppendEnvironment(deps.Environment.EnvVars)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
