package publish

import (
	"path"
	"slices"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/common/parser"
	"gitlab.kilic.dev/devops/pipes/helm/setup"
)

func HelmChartVersionsParent(tl *TaskList) *Task {
	return tl.CreateTask("versions").
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobSequence(
				JobParallel(
					HelmChartVersionsFromUser(tl).Job(),
					HelmChartVersionsFromFile(tl).Job(),
				),
				job,
			)
		}).
		Set(func(t *Task) error {
			C.Versions = slices.Compact(C.Versions)

			t.Log.Infof(
				"Helm Chart versions: %s", strings.Join(C.Versions, ", "),
			)

			return nil
		})
}

func HelmChartVersionsFromUser(tl *TaskList) *Task {
	return tl.CreateTask("versions", "user").
		Set(func(t *Task) error {
			// add all the specified version
			for _, v := range slices.Compact(P.HelmChart.Versions) {
				if err := AddHelmChartVersion(t, v); err != nil {
					return err
				}
			}

			return nil
		})
}

func HelmChartVersionsFromFile(tl *TaskList) *Task {
	return tl.CreateTask("versions", "file").
		ShouldDisable(func(t *Task) bool {
			return P.HelmChart.VersionFile == ""
		}).
		Set(func(t *Task) error {
			// add versions through versions file
			versions, err := parser.ParseTagsFile(t.Log, path.Join(setup.P.Cwd, P.HelmChart.VersionFile), P.HelmChart.VersionFileStrict)

			if err != nil {
				return err
			}

			for _, v := range versions {
				t.CreateSubtask(v).
					Set(func(t *Task) error {
						return AddHelmChartVersion(t, v)
					}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}
