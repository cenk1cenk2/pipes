package publish

import (
	"fmt"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/helm/login"
	"gitlab.kilic.dev/devops/pipes/helm/setup"
)

func HelmPackage(tl *TaskList) *Task {
	return tl.CreateTask("package").Set(func(t *Task) error {
		for _, version := range C.Versions {

			t.CreateSubtask(fmt.Sprintf("%s@%s", setup.C.Chart.Name(), version)).
				Set(func(t *Task) error {
					t.Log.Infof("Packaging Helm Chart with version: %s@%s", setup.C.Chart.Name(), version)

					t.CreateCommand(
						"helm",
						"package",
						"-d",
						P.HelmChart.Destination,
						".",
						"--version",
						version,
					).
						SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
						SetDir(setup.P.Cwd).
						Set(func(c *Command) error {
							if P.HelmChart.AppVersion != "" {
								c.AppendArgs("--app-version", P.HelmChart.AppVersion)
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

func HelmPublish(tl *TaskList) *Task {
	return tl.CreateTask("publish").Set(func(t *Task) error {
		for _, version := range C.Versions {
			t.Log.Infof("Publishing Helm Chart with version: %s to %s", version, login.P.HelmRegistry.Uri)

			t.CreateSubtask(fmt.Sprintf("%s@%s", setup.C.Chart.Name(), version)).
				Set(func(t *Task) error {
					t.CreateCommand(
						"helm",
						"push",
						filepath.Join(P.HelmChart.Destination, fmt.Sprintf("%s-%s.tgz", setup.C.Chart.Name(), version)),
						login.P.HelmRegistry.Uri,
					).
						SetLogLevel(LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEFAULT).
						SetDir(setup.P.Cwd).
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
