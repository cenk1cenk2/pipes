package lint

import (
	"github.com/urfave/cli/v3"
)

const CATEGORY_HELM_LINT = "Helm Lint"

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category:    CATEGORY_HELM_LINT,
		Name:        "helm.lint.kubernetes.version",
		Usage:       "Kubernetes version to use for linting charts.",
		Sources:     cli.NewValueSourceChain(cli.EnvVar("HELM_LINT_KUBERNETES_VERSION")),
		Required:    false,
		Value:       "",
		Destination: &P.Kubernetes.Version,
	},

	&cli.BoolFlag{
		Category: CATEGORY_HELM_LINT,
		Name:     "helm.lint.should-template",
		Usage:    "If set to true, the lint command will also template the chart.",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_LINT_SHOULD_TEMPLATE"),
		),
		Required:    false,
		Value:       true,
		Destination: &P.ShouldTemplate,
	},
}
