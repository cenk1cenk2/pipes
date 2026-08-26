package publish

import (
	"gitlab.kilic.dev/devops/pipes/internal/versions"
)

// The collector reads the parsed flags and the working directory the setup task
// list resolved, so it is only built from inside a task list.
func HelmChartVersions(deps Deps) *versions.Collector {
	return &versions.Collector{
		Name:  "versions",
		Label: "Helm Chart versions",

		FromUser: P.Chart.Versions,

		File:       P.Chart.VersionFile,
		FileStrict: P.Chart.VersionFileStrict,
		FileDir:    deps.Tool.Cwd,

		Templates: P.Chart.VersionsTemplate,
		Sanitize:  P.Chart.VersionsSanitize,
	}
}
