package publish

import (
	"slices"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/helm/setup"
	"gitlab.kilic.dev/devops/pipes/internal/versions"
)

// The collector reads the parsed flags and the resolved working directory, so it is only built from inside a task list.
func HelmChartVersions() *versions.Collector {
	return &versions.Collector{
		Name: "versions",

		FromUser: P.Chart.Versions,

		File:       P.Chart.VersionFile,
		FileStrict: P.Chart.VersionFileStrict,
		FileDir:    setup.C.Cwd,

		Templates: P.Chart.VersionsTemplate,
		Sanitize:  P.Chart.VersionsSanitize,
	}
}

// A chart has no notion of a latest version, so the parent waits on the other two sources only.
func HelmChartVersionsParent(tl *TaskList) *Task {
	collector := HelmChartVersions()

	return tl.CreateTask("versions").
		SetJobWrapper(func(job Job, t *Task) Job {
			return JobSequence(
				JobParallel(
					collector.UserTask(tl, &C.Versions).Job(),
					collector.FileTask(tl, &C.Versions).Job(),
				),
				job,
			)
		}).
		Set(func(t *Task) error {
			C.Versions = slices.Compact(C.Versions)

			t.Log.Infof("Helm Chart versions: %s", strings.Join(C.Versions, ", "))

			return nil
		})
}
