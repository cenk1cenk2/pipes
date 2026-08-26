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
	tagsfile.NewFlags(&P.Chart.VersionFile, "", &P.Chart.VersionFileStrict, false),
	[]ucli.Flag{
		&ucli.StringFlag{
			Category:    CATEGORY_HELM_CHART,
			Name:        "helm.publish.chart.target",
			Sources:     cli.EnvVars("HELM_CHART_TARGET", "HELM_PUBLISH_CHART_TARGET"),
			Usage:       "Helm chart repository target to publish to.",
			Required:    true,
			Destination: &P.Chart.Target,
		},

		&ucli.StringSliceFlag{
			Category:    CATEGORY_HELM_CHART,
			Name:        "helm.publish.chart.versions",
			Sources:     cli.EnvVars("HELM_CHART_VERSIONS", "HELM_PUBLISH_CHART_VERSIONS"),
			Usage:       "Versions for the helm chart to be published.",
			Required:    false,
			Destination: &P.Chart.Versions,
		},

		cli.YAMLFlag(&ucli.StringFlag{
			Category: CATEGORY_HELM_CHART,
			Name:     "helm.publish.chart.versions-template",
			Sources:  cli.EnvVars("HELM_CHART_VERSIONS_TEMPLATE", "HELM_PUBLISH_CHART_VERSIONS_TEMPLATE"),
			Usage: strings.TrimSpace(`
    Modifies every version that matches a certain condition.
    Template is interpolated with the given matches in the regular expression.

    format(yaml([]struct{ match: RegExp, template: Template(match) }))
    `),
			Required: false,
			Value:    "[]",
		}, &P.Chart.VersionsTemplate),

		cli.YAMLFlag(&ucli.StringFlag{
			Category: CATEGORY_HELM_CHART,
			Name:     "helm.publish.chart.versions-sanitize",
			Sources:  cli.EnvVars("HELM_CHART_SANITIZE_VERSIONS", "HELM_PUBLISH_CHART_VERSIONS_SANITIZE"),
			Usage: strings.TrimSpace(`
    Sanitizes the given regex pattern out of version name.
    Template is interpolated with the given matches in the regular expression.

    format(yaml([]struct{ match: RegExp, template: Template(match) }))
    `),
			Required: false,
			Value:    DEFAULT_SANITIZE_VERSIONS,
		}, &P.Chart.VersionsSanitize),

		&ucli.StringFlag{
			Category:    CATEGORY_HELM_CHART,
			Name:        "helm.publish.chart.destination",
			Sources:     cli.EnvVars("HELM_CHART_DESTINATION", "HELM_PUBLISH_CHART_DESTINATION"),
			Usage:       "Destination directory for the packaged helm chart.",
			Required:    false,
			Value:       "./dist/",
			Destination: &P.Chart.Destination,
		},

		&ucli.StringFlag{
			Category:    CATEGORY_HELM_CHART,
			Name:        "helm.publish.chart.app-version",
			Sources:     cli.EnvVars("HELM_CHART_APP_VERSION", "HELM_PUBLISH_CHART_APP_VERSION"),
			Usage:       "Application version for the packaged helm chart.",
			Required:    false,
			Value:       "",
			Destination: &P.Chart.AppVersion,
		},
	})
