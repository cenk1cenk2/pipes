package setup

import (
	"context"
	"fmt"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	helmv2loader "helm.sh/helm/v4/pkg/chart/v2/loader"
)

func initialize(tl *TaskList) *Task {
	return tl.CreateTask("init").
		Set(func(_ context.Context, t *Task) error {
			C.Cwd = P.Cwd

			t.Log.Debug(fmt.Sprintf("Working directory: %s", C.Cwd))

			return nil
		})
}

func version(tl *TaskList) *Task {
	return tl.CreateTask("version").
		Set(func(_ context.Context, t *Task) error {
			t.CreateCommand("helm", "version").
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				ShouldRunAfter(func(_ context.Context, c *Command) error {
					C.Version = strings.TrimSpace(strings.Join(c.GetCombinedStream(), "\n"))

					c.Log.Info(fmt.Sprintf("helm version: %s", C.Version))

					return nil
				}).
				EnableStreamRecording().
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

func read(tl *TaskList) *Task {
	return tl.CreateTask("read").
		Set(func(_ context.Context, t *Task) error {
			chart, err := helmv2loader.Load(C.Cwd)
			if err != nil {
				return fmt.Errorf("Error loading helm chart: %v in %s", err, C.Cwd)
			} else if chart == nil {
				return fmt.Errorf("Can not load helm chart: %s", C.Cwd)
			}

			t.Log.Info(fmt.Sprintf("Chart Name: %s", chart.Metadata.Name))

			C.Chart = chart

			return nil
		})
}
