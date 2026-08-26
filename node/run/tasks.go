package run

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
)

func RunNodeScript(tl *TaskList, deps Deps) *Task {
	return tl.CreateTask("run", C.Script).
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

					if C.Script != "" {
						tmpl, err := InlineTemplate(C.Script, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.AppendArgs(deps.Node.PackageManager.Commands.RunDelimiter...)

					if C.ScriptArgs != "" {
						tmpl, err := InlineTemplate(C.ScriptArgs, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.SetDir(P.Run.Cwd)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
