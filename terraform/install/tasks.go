package install

import (
	"context"
	"fmt"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/terraform/setup"
)

func install(tl *TaskList) *Task {
	return tl.CreateTask("install").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand(
				"terraform",
				"init",
				"-input=false",
			).
				Set(func(_ context.Context, c *Command) error {
					if P.Install.Reconfigure {
						t.Log.Info("Will reconfigure state.")

						c.AppendArgs("-reconfigure")
					}

					if P.Install.UseLockfile {
						t.Log.Info("Using lockfile.")

						c.AppendArgs("-lockfile=readonly")
					}

					if P.Install.Args != "" {
						c.AppendArgs(P.Install.Args)
					}

					return nil
				}).
				SetDir(setup.C.Cwd).
				AppendEnvironment(setup.C.Env).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

// terraform never removes a provider version from the plugin cache, so the cached
// directory and the artifact that links into it grow with every provider bump unless the
// versions this root did not select are dropped.
func prune(tl *TaskList) *Task {
	return tl.CreateTask("prune").
		ShouldDisable(func(_ *Task) bool {
			return setup.C.PluginCache == ""
		}).
		Set(func(_ context.Context, t *Task) error {
			removed, freed, err := pruneCache(setup.C.PluginCache, filepath.Join(setup.C.Cwd, ".terraform", "providers"))
			if err != nil {
				return err
			}

			for _, version := range removed {
				t.Log.Info(fmt.Sprintf("Pruned from the plugin cache: %s", version))
			}

			t.Log.Info(fmt.Sprintf("Pruned %d provider versions from the plugin cache, freed %d bytes.", len(removed), freed))

			return nil
		})
}
