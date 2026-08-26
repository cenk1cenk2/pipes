package publish

import (
	"strings"

	ucli "github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/git"
	"gitlab.kilic.dev/devops/pipes/internal/tagsfile"

	. "github.com/cenk1cenk2/plumber/v6"
)

//revive:disable:line-length-limit

const (
	CATEGORY_HELM_CHART = "Helm Chart"

	DEFAULT_SANITIZE_VERSIONS = `[
    { "match": "([^/]*)/(.*)", "template": "{{ index $ 1 | upper }}_{{ index $ 2 }}" }
]`
)

var Flags = CombineFlags(
	git.NewFlags(&P.Git),
	tagsfile.NewFlags(&P.HelmChart.VersionFile, "", &P.HelmChart.VersionFileStrict, false),
	[]ucli.Flag{
		&ucli.StringFlag{
			Category:    CATEGORY_HELM_CHART,
			Name:        "helm-chart.target",
			Sources:     cli.EnvVars("HELM_CHART_TARGET"),
			Usage:       "Helm chart repository target to publish to.",
			Required:    true,
			Destination: &P.HelmChart.Target,
		},

		&ucli.StringSliceFlag{
			Category:    CATEGORY_HELM_CHART,
			Name:        "helm-chart.versions",
			Sources:     cli.EnvVars("HELM_CHART_VERSIONS"),
			Usage:       "Versions for the helm chart to be published.",
			Required:    false,
			Destination: &P.HelmChart.Versions,
		},

		cli.YAMLFlag(&ucli.StringFlag{
			Category: CATEGORY_HELM_CHART,
			Name:     "helm-chart.versions-template",
			Sources:  cli.EnvVars("HELM_CHART_VERSIONS_TEMPLATE"),
			Usage: strings.TrimSpace(`
    Modifies every version that matches a certain condition.
    Template is interpolated with the given matches in the regular expression.

    format(yaml([]struct{ match: RegExp, template: Template(match) }))
    `),
			Required: false,
			Value:    "[]",
		}, &P.HelmChart.VersionsTemplate),

		cli.YAMLFlag(&ucli.StringFlag{
			Category: CATEGORY_HELM_CHART,
			Name:     "helm-chart.sanitize-versions",
			Sources:  cli.EnvVars("HELM_CHART_SANITIZE_VERSIONS"),
			Usage: strings.TrimSpace(`
    Sanitizes the given regex pattern out of version name.
    Template is interpolated with the given matches in the regular expression.

    format(yaml([]struct{ match: RegExp, template: Template(match) }))
    `),
			Required: false,
			Value:    DEFAULT_SANITIZE_VERSIONS,
		}, &P.HelmChart.VersionsSanitize),

		&ucli.StringFlag{
			Category:    CATEGORY_HELM_CHART,
			Name:        "helm-chart.destination",
			Sources:     cli.EnvVars("HELM_CHART_DESTINATION"),
			Usage:       "Destination directory for the packaged helm chart.",
			Required:    false,
			Value:       "./dist/",
			Destination: &P.HelmChart.Destination,
		},

		&ucli.StringFlag{
			Category:    CATEGORY_HELM_CHART,
			Name:        "helm-chart.app-version",
			Sources:     cli.EnvVars("HELM_CHART_APP_VERSION"),
			Usage:       "Application version for the packaged helm chart.",
			Required:    false,
			Value:       "",
			Destination: &P.HelmChart.AppVersion,
		},
	})
