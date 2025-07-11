package pipe

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	environment "gitlab.kilic.dev/devops/pipes/select-env/setup"
)

func AddNodeModules(tl *TaskList) *Task {
	return tl.CreateTask("packages", "node").
		Set(func(t *Task) error {
			t.CreateCommand(
				setup.C.PackageManager.Exe,
			).
				Set(func(c *Command) error {
					ctx := environment.EnvironmentTemplate{
						Environment: environment.C.Environment,
						EnvVars:     environment.C.EnvVars,
					}

					if P.NodeAdd.Global {
						c.AppendArgs(setup.C.PackageManager.Commands.Global...)
					}

					c.AppendArgs(setup.C.PackageManager.Commands.Add...)

					if P.NodeAdd.ScriptArgs != "" {
						tmpl, err := InlineTemplate(P.NodeAdd.ScriptArgs, ctx)

						if err != nil {
							return err
						}

						c.AppendArgs(tmpl)
					}

					c.AppendArgs(P.NodeAdd.Packages...)

					c.SetDir(P.NodeAdd.Cwd)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
