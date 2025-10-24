package publish

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/common/flags"
	"gopkg.in/yaml.v2"

	. "github.com/cenk1cenk2/plumber/v6"
)

//revive:disable:line-length-limit

const (
	CATEGORY_HELM_CHART = "Helm Chart"
)

var Flags = CombineFlags(flags.NewGitFlags(
	flags.GitFlagsSetup{
		GitBranchDestination: &P.Git.Branch,
		GitTagDestination:    &P.Git.Tag,
	},
), flags.NewTagsFileFlags(
	flags.TagsFileFlagsSetup{
		TagsFileDestination: &P.HelmChart.VersionFile,
		TagsFileRequired:    false,
	},
), flags.NewTagsFileStrictFlags(
	flags.TagsFileStrictFlagsSetup{
		TagsFileStrictDestination: &P.HelmChart.VersionFileStrict,
		TagsFileStrictRequired:    false,
	},
), []cli.Flag{
	&cli.StringSliceFlag{
		Category: CATEGORY_HELM_CHART,
		Name:     "helm-chart.versions",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_CHART_VERSIONS"),
		),
		Usage:       "Versions for the helm chart to be published.",
		Required:    true,
		Destination: &P.HelmChart.Versions,
	},

	&cli.StringFlag{
		Category: CATEGORY_HELM_CHART,
		Name:     "helm-chart.versions-template",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_CHART_VERSIONS_TEMPLATE"),
		),
		Usage: strings.TrimSpace(`
    Modifies every version that matches a certain condition.
    Template is interpolated with the given matches in the regular expression.

    format(yaml([]struct{ match: RegExp, template: Template(match) }))
    `),
		Required:         false,
		Value:            "[]",
		ValidateDefaults: true,
		Validator: func(v string) error {
			if v == "" {
				return nil
			}

			if err := yaml.Unmarshal([]byte(v), &P.HelmChart.VersionsTemplate); err != nil {
				return fmt.Errorf("Cannot unmarshal helm chart templating version conditions: %w", err)
			}

			return nil
		},
	},

	&cli.StringFlag{
		Category: CATEGORY_HELM_CHART,
		Name:     "helm-chart.sanitize-versions",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_CHART_SANITIZE_VERSIONS"),
		),
		Usage: strings.TrimSpace(`
    Sanitizes the given regex pattern out of version name.
    Template is interpolated with the given matches in the regular expression.

    format(yaml([]struct{ match: RegExp, template: Template(match) }))
    `),
		Required:         false,
		Value:            flags.FLAG_DEFAULT_DOCKER_IMAGE_SANITIZE_TAGS,
		ValidateDefaults: true,
		Validator: func(v string) error {
			if v == "" {
				return nil
			}

			if err := yaml.Unmarshal([]byte(v), &P.HelmChart.VersionsSanitize); err != nil {
				return fmt.Errorf("Cannot unmarshal helm chart sanitizing versions conditions: %w", err)
			}

			return nil
		},
	},

	&cli.StringFlag{
		Category: CATEGORY_HELM_CHART,
		Name:     "helm-chart.destination",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_CHART_DESTINATION"),
		),
		Usage:       "Destination directory for the packaged helm chart.",
		Required:    false,
		Value:       ".",
		Destination: &P.HelmChart.Destination,
	},

	&cli.StringFlag{
		Category: CATEGORY_HELM_CHART,
		Name:     "helm-chart.app-version",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_CHART_APP_VERSION"),
		),
		Usage:       "Application version for the packaged helm chart.",
		Required:    false,
		Value:       "",
		Destination: &P.HelmChart.AppVersion,
	},
})
