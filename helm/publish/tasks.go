package publish

import (
	"context"
	"fmt"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v7"
	"gitlab.kilic.dev/devops/pipes/helm/setup"
)

func packageTask(tl *TaskList) *Task {
	return tl.CreateTask("package").
		ShouldDisable(func(t *Task) bool {
			if len(C.Versions) == 0 {
				t.Log.Warn("No version to package.")

				return true
			}

			return false
		}).
		Set(func(ctx context.Context, t *Task) error {
			for _, version := range C.Versions {
				t.CreateSubtask(fmt.Sprintf("%s@%s", setup.C.Chart.Name(), version)).
					Set(func(_ context.Context, t *Task) error {
						t.Log.Info(fmt.Sprintf("Packaging Helm Chart with version: %s@%s", setup.C.Chart.Name(), version))

						t.CreateCommand(
							"helm",
							"package",
							"-d",
							P.Chart.Destination,
							".",
							"--version",
							version,
						).
							SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDefault).
							SetDir(setup.C.Cwd).
							Set(func(_ context.Context, c *Command) error {
								if P.Chart.AppVersion != "" {
									c.AppendArgs("--app-version", P.Chart.AppVersion)
								}

								return nil
							}).
							AddSelfToTheTask()

						return nil
					}).
					ShouldRunAfter(func(ctx context.Context, t *Task) error {
						return t.RunCommandJobAsJobSequence(ctx)
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunSubtasks(ctx)
		})
}

func publish(tl *TaskList) *Task {
	return tl.CreateTask("publish").
		ShouldDisable(func(t *Task) bool {
			if len(C.Versions) == 0 {
				t.Log.Warn("No version to publish.")

				return true
			}

			return false
		}).
		Set(func(ctx context.Context, t *Task) error {
			for _, version := range C.Versions {
				t.Log.Info(fmt.Sprintf("Publishing Helm Chart with version: %s to %s", version, P.Chart.Target))

				t.CreateSubtask(fmt.Sprintf("%s@%s", setup.C.Chart.Name(), version)).
					Set(func(_ context.Context, t *Task) error {
						t.CreateCommand(
							"helm",
							"push",
							filepath.Join(P.Chart.Destination, fmt.Sprintf("%s-%s.tgz", setup.C.Chart.Name(), version)),
							P.Chart.Target,
						).
							SetLogLevel(LogLevelDefault, LogLevelDefault, LogLevelDefault).
							SetDir(setup.C.Cwd).
							AddSelfToTheTask()

						return nil
					}).
					ShouldRunAfter(func(ctx context.Context, t *Task) error {
						return t.RunCommandJobAsJobSequence(ctx)
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunSubtasks(ctx)
		})
}
