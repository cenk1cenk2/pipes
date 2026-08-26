package publish

import (
	"fmt"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v6"
)

func HelmPackage(tl *TaskList, deps Deps) *Task {
	return tl.CreateTask("package").
		ShouldDisable(func(t *Task) bool {
			if len(C.Versions) == 0 {
				t.Log.Warnf("No version to package.")

				return true
			}

			return false
		}).
		Set(func(t *Task) error {
			for _, version := range C.Versions {
				t.CreateSubtask(fmt.Sprintf("%s@%s", deps.Tool.Chart.Name(), version)).
					Set(func(t *Task) error {
						t.Log.Infof("Packaging Helm Chart with version: %s@%s", deps.Tool.Chart.Name(), version)

						t.CreateCommand(
							"helm",
							"package",
							"-d",
							P.Chart.Destination,
							".",
							"--version",
							version,
						).
							SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
							SetDir(deps.Tool.Cwd).
							Set(func(c *Command) error {
								if P.Chart.AppVersion != "" {
									c.AppendArgs("--app-version", P.Chart.AppVersion)
								}

								return nil
							}).
							AddSelfToTheTask()

						return nil
					}).
					ShouldRunAfter(func(t *Task) error {
						return t.RunCommandJobAsJobSequence()
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}

func HelmPublish(tl *TaskList, deps Deps) *Task {
	return tl.CreateTask("publish").
		ShouldDisable(func(t *Task) bool {
			if len(C.Versions) == 0 {
				t.Log.Warnf("No version to publish.")

				return true
			}

			return false
		}).
		Set(func(t *Task) error {
			for _, version := range C.Versions {
				t.Log.Infof("Publishing Helm Chart with version: %s to %s", version, P.Chart.Target)

				t.CreateSubtask(fmt.Sprintf("%s@%s", deps.Tool.Chart.Name(), version)).
					Set(func(t *Task) error {
						t.CreateCommand(
							"helm",
							"push",
							filepath.Join(P.Chart.Destination, fmt.Sprintf("%s-%s.tgz", deps.Tool.Chart.Name(), version)),
							P.Chart.Target,
						).
							SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
							SetDir(deps.Tool.Cwd).
							AddSelfToTheTask()

						return nil
					}).
					ShouldRunAfter(func(t *Task) error {
						return t.RunCommandJobAsJobSequence()
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}
