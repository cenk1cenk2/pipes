package publish

import (
	"strings"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
	"gitlab.kilic.dev/devops/pipes/internal/tagsfile"

	. "github.com/cenk1cenk2/plumber/v6"
)

//revive:disable:line-length-limit

const (
	CategoryHelmChart = "Helm Chart"

	DefaultSanitizeVersions = `[
    { "match": "([^/]*)/(.*)", "template": "{{ index $ 1 | upper }}_{{ index $ 2 }}" }
]`
)

var Flags = CombineFlags(
	tagsfile.NewFlags(tagsfile.Options{Destination: &P.Chart.VersionFile, Strict: &P.Chart.VersionFileStrict}),
	[]cli.Flag{
		&cli.StringFlag{
			Category: CategoryHelmChart,
			Name:     "helm.publish.chart.target",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("HELM_PUBLISH_CHART_TARGET"),
				cli.EnvVar("HELM_CHART_TARGET"),
			),
			Usage:       "Helm chart repository target to publish to.",
			Required:    true,
			Destination: &P.Chart.Target,
		},

		&cli.StringSliceFlag{
			Category: CategoryHelmChart,
			Name:     "helm.publish.chart.versions",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("HELM_PUBLISH_CHART_VERSIONS"),
				cli.EnvVar("HELM_CHART_VERSIONS"),
			),
			Usage:       "Versions for the helm chart to be published.",
			Required:    false,
			Destination: &P.Chart.Versions,
		},

		flags.YAMLFlag(&P.Chart.VersionsTemplate, &cli.StringFlag{
			Category: CategoryHelmChart,
			Name:     "helm.publish.chart.versions-template",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("HELM_PUBLISH_CHART_VERSIONS_TEMPLATE"),
				cli.EnvVar("HELM_CHART_VERSIONS_TEMPLATE"),
			),
			Usage: strings.TrimSpace(`
Modifies every version that matches a certain condition.
Template is interpolated with the given matches in the regular expression.

format(yaml([]struct{ match: RegExp, template: Template(match) }))
`),
			Required: false,
			Value:    "[]",
		}),

		flags.YAMLFlag(&P.Chart.VersionsSanitize, &cli.StringFlag{
			Category: CategoryHelmChart,
			Name:     "helm.publish.chart.versions-sanitize",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("HELM_PUBLISH_CHART_VERSIONS_SANITIZE"),
				cli.EnvVar("HELM_CHART_SANITIZE_VERSIONS"),
			),
			Usage: strings.TrimSpace(`
Sanitizes the given regex pattern out of version name.
Template is interpolated with the given matches in the regular expression.

format(yaml([]struct{ match: RegExp, template: Template(match) }))
`),
			Required: false,
			Value:    DefaultSanitizeVersions,
		}),

		&cli.StringFlag{
			Category: CategoryHelmChart,
			Name:     "helm.publish.chart.destination",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("HELM_PUBLISH_CHART_DESTINATION"),
				cli.EnvVar("HELM_CHART_DESTINATION"),
			),
			Usage:       "Destination directory for the packaged helm chart.",
			Required:    false,
			Value:       "./dist/",
			Destination: &P.Chart.Destination,
		},

		&cli.StringFlag{
			Category: CategoryHelmChart,
			Name:     "helm.publish.chart.app-version",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("HELM_PUBLISH_CHART_APP_VERSION"),
				cli.EnvVar("HELM_CHART_APP_VERSION"),
			),
			Usage:       "Application version for the packaged helm chart.",
			Required:    false,
			Value:       "",
			Destination: &P.Chart.AppVersion,
		},
	})
