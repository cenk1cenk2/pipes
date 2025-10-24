package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/go/setup"
	helmv2loader "helm.sh/helm/v4/pkg/chart/v2/loader"
)

func HelmVersion(tl *TaskList) *Task {
	return tl.CreateTask("version").
		Set(func(t *Task) error {
			t.CreateCommand(
				"helm",
				"version",
			).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func HelmLoadChart(tl *TaskList) *Task {
	return tl.CreateTask("read").
		Set(func(t *Task) error {
			chart, err := helmv2loader.Load(setup.P.Cwd)
			if err != nil {
				return err
			}

			t.Log.Infof("Chart Name: %s", chart.Metadata.Name)

			C.Chart = chart

			return nil
		})
}
