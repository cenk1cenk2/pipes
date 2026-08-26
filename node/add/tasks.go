package pipe

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/environment"
)

func AddNodeModules(tl *TaskList, deps Deps) *Task {
	return tl.CreateTask("packages", "node").
		Set(func(t *Task) error {
			t.CreateCommand(
				deps.Node.PackageManager.Exe,
			).
				Set(func(c *Command) error {
					ctx := environment.Template{
						Environment: deps.Environment.Environment,
						EnvVars:     deps.Environment.EnvVars,
					}

					if P.Add.Global {
						c.AppendArgs(deps.Node.PackageManager.Commands.Global...)
					}

					c.AppendArgs(deps.Node.PackageManager.Commands.Add...)

					if P.Add.ScriptArgs != "" {
						tmpl, err := InlineTemplate(P.Add.ScriptArgs, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.AppendArgs(P.Add.Packages...)

					c.SetDir(P.Add.Cwd)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
