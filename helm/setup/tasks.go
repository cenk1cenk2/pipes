package setup

import (
	"fmt"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	helmv2loader "helm.sh/helm/v4/pkg/chart/v2/loader"
)

func initialize(tl *TaskList) *Task {
	return tl.CreateTask("init").
		Set(func(t *Task) error {
			C.Cwd = P.Cwd

			t.Log.Debugf("Working directory: %s", C.Cwd)

			return nil
		})
}

func version(tl *TaskList) *Task {
	return tl.CreateTask("version").
		Set(func(t *Task) error {
			t.CreateCommand("helm", "version").
				SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG, LOG_LEVEL_DEBUG).
				ShouldRunAfter(func(c *Command) error {
					C.Version = strings.TrimSpace(strings.Join(c.GetCombinedStream(), "\n"))

					c.Log.Infof("helm version: %s", C.Version)

					return nil
				}).
				EnableStreamRecording().
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

func read(tl *TaskList) *Task {
	return tl.CreateTask("read").
		Set(func(t *Task) error {
			chart, err := helmv2loader.Load(C.Cwd)
			if err != nil {
				return fmt.Errorf("Error loading helm chart: %v in %s", err, C.Cwd)
			} else if chart == nil {
				return fmt.Errorf("Can not load helm chart: %s", C.Cwd)
			}

			t.Log.Infof("Chart Name: %s", chart.Metadata.Name)

			C.Chart = chart

			return nil
		})
}
