package publish

import (
	"gitlab.kilic.dev/devops/pipes/internal/versions"
	"gitlab.kilic.dev/devops/pipes/pipes/helm/setup"
)

// The collector reads the parsed flags and the working directory the setup task
// list resolved, so it is only built from inside a task list.
func HelmChartVersions() *versions.Collector {
	return &versions.Collector{
		Name:  "versions",
		Label: "Helm Chart versions",

		FromUser: P.HelmChart.Versions,

		File:       P.HelmChart.VersionFile,
		FileStrict: P.HelmChart.VersionFileStrict,
		FileDir:    setup.C.Cwd,

		Templates: P.HelmChart.VersionsTemplate,
		Sanitize:  P.HelmChart.VersionsSanitize,
	}
}
