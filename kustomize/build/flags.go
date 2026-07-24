package build

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/common/gitlab"
)

//revive:disable:line-length-limit

const (
	CATEGORY_KUSTOMIZE_BUILD = "Kustomize Build"
)

var Flags = CombineFlags(
	[]cli.Flag{

		&cli.BoolFlag{
			Category: CATEGORY_KUSTOMIZE_BUILD,
			Name:     "kustomize-build.enable-helm",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("KUSTOMIZE_ENABLE_HELM"),
			),
			Usage:       "Enable the Helm chart inflation generator while building overlays.",
			Required:    false,
			Value:       true,
			Destination: &P.EnableHelm,
		},

		&cli.StringFlag{
			Category: CATEGORY_KUSTOMIZE_BUILD,
			Name:     "kustomize-build.helm-command",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("KUSTOMIZE_HELM_COMMAND"),
			),
			Usage:       "Helm binary to use for the Helm chart inflation generator.",
			Required:    false,
			Value:       "helm",
			Destination: &P.HelmCommand,
		},

		&cli.StringFlag{
			Category: CATEGORY_KUSTOMIZE_BUILD,
			Name:     "kustomize-build.load-restrictor",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("KUSTOMIZE_LOAD_RESTRICTOR"),
			),
			Usage:       "Load restrictor for Kustomize file access. format(enum(\"rootOnly\", \"none\"))",
			Required:    false,
			Value:       "none",
			Destination: &P.LoadRestrictor,
		},

		&cli.StringFlag{
			Category: CATEGORY_KUSTOMIZE_BUILD,
			Name:     "kustomize-build.kube-version",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("KUSTOMIZE_KUBE_VERSION"),
			),
			Usage:       "Kubernetes version passed to the Helm chart inflation generator.",
			Required:    false,
			Value:       "",
			Destination: &P.KubeVersion,
		},

		&cli.StringFlag{
			Category: CATEGORY_KUSTOMIZE_BUILD,
			Name:     "kustomize-build.output",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("KUSTOMIZE_BUILD_OUTPUT"),
			),
			Usage:       "Output directory to write the rendered manifests per overlay. Leave empty to skip writing.",
			Required:    false,
			Value:       "",
			Destination: &P.Output,
		},

		&cli.BoolFlag{
			Category: CATEGORY_KUSTOMIZE_BUILD,
			Name:     "kustomize-build.fail-fast",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("KUSTOMIZE_BUILD_FAIL_FAST"),
			),
			Usage:       "Fail on the first overlay that can not be built instead of collecting all results.",
			Required:    false,
			Value:       false,
			Destination: &P.FailFast,
		},
	},
	gitlab.NewMergeRequestReportFlags(&P.MergeRequestReport),
)
