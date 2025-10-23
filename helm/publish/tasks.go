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
			t.Log.Infof("Packaging Helm Chart with version: %s", version)
			t.CreateSubtask(version).
				Set(func(t *Task) error {
					t.CreateCommand(
						"helm",
						"package",
						"-d",
						P.HelmChart.Destination,
						".",
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

func HelmPublish(tl *TaskList) *Task {
	return tl.CreateTask("publish").Set(func(t *Task) error {
		for _, version := range C.Versions {
			t.Log.Infof("Publishing Helm Chart with version: %s to %s", version, login.P.HelmRegistry.Uri)

			t.CreateSubtask(version).
				Set(func(t *Task) error {
					t.CreateCommand(
						"helm",
						"push",
						filepath.Join(P.HelmChart.Destination, fmt.Sprintf("%s-%s.tgz", P.HelmChart.Name, version)),
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
