package setup

import (
	"fmt"

	"github.com/cenk1cenk2/plumber/v6"
	helmv2loader "helm.sh/helm/v4/pkg/chart/v2/loader"
)

func HelmLoadChart(tl *plumber.TaskList) *plumber.Task {
	return tl.CreateTask("read").
		Set(func(t *plumber.Task) error {
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
