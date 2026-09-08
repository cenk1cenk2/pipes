package lint

import (
	"github.com/urfave/cli/v3"
)

const CATEGORY_HELM_LINT = "Helm Lint"

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_HELM_LINT,
		Name:     "helm.lint.kubernetes.version",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_LINT_KUBERNETES_VERSION"),
			cli.EnvVar("KUBERNETES_VERSION"),
		),
		Usage:       "Kubernetes version to use for linting charts.",
		Required:    false,
		Value:       "",
		Destination: &P.Kubernetes.Version,
	},

	&cli.BoolFlag{
		Category: CATEGORY_HELM_LINT,
		Name:     "helm.lint.should-template",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_LINT_SHOULD_TEMPLATE"),
		),
		Usage:       "Template the chart while linting.",
		Required:    false,
		Value:       true,
		Destination: &P.ShouldTemplate,
	},
}
