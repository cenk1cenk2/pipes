package build

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/ci"
	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
)

//revive:disable:line-length-limit

type LoadRestrictor string

const (
	CategoryKustomizeBuild = "Kustomize Build"

	LoadRestrictorRootOnly LoadRestrictor = "rootOnly"
	LoadRestrictorNone     LoadRestrictor = "none"
)

var Flags = CombineFlags(
	[]cli.Flag{

		&cli.BoolFlag{
			Category: CategoryKustomizeBuild,
			Name:     "kustomize.build.enable-helm",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("KUSTOMIZE_BUILD_ENABLE_HELM"),
				cli.EnvVar("KUSTOMIZE_ENABLE_HELM"),
			),
			Usage:       "Enable the Helm chart inflation generator while building overlays.",
			Required:    false,
			Value:       true,
			Destination: &P.EnableHelm,
		},

		&cli.StringFlag{
			Category: CategoryKustomizeBuild,
			Name:     "kustomize.build.helm-command",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("KUSTOMIZE_BUILD_HELM_COMMAND"),
				cli.EnvVar("KUSTOMIZE_HELM_COMMAND"),
			),
			Usage:       "Helm binary to use for the Helm chart inflation generator.",
			Required:    false,
			Value:       "helm",
			Destination: &P.HelmCommand,
		},

		&cli.StringFlag{
			Category: CategoryKustomizeBuild,
			Name:     "kustomize.build.load-restrictor",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("KUSTOMIZE_BUILD_LOAD_RESTRICTOR"),
				cli.EnvVar("KUSTOMIZE_LOAD_RESTRICTOR"),
			),
			Usage:       `Restricts which files Kustomize may load, where lifting the restriction matches the behaviour of ArgoCD. format(enum("rootOnly", "none"))`,
			Required:    false,
			Value:       string(LoadRestrictorNone),
			Destination: &P.LoadRestrictor,
		},

		&cli.StringFlag{
			Category: CategoryKustomizeBuild,
			Name:     "kustomize.build.kube-version",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("KUSTOMIZE_BUILD_KUBE_VERSION"),
				cli.EnvVar("KUSTOMIZE_KUBE_VERSION"),
				cli.EnvVar("KUBERNETES_VERSION"),
			),
			Usage:       "Kubernetes version passed to the Helm chart inflation generator.",
			Required:    false,
			Value:       "",
			Destination: &P.KubeVersion,
		},

		&cli.StringFlag{
			Category: CategoryKustomizeBuild,
			Name:     "kustomize.build.output",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("KUSTOMIZE_BUILD_OUTPUT"),
			),
			Usage:       "Output directory to write the rendered manifests per overlay. Leave empty to skip writing.",
			Required:    false,
			Value:       "",
			Destination: &P.Output,
		},

		&cli.StringFlag{
			Category: CategoryKustomizeBuild,
			Name:     "kustomize.build.diff-base",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("KUSTOMIZE_BUILD_DIFF_BASE"),
			),
			Usage:       "Commit the overlays are rendered on again for the merge request report, which diffs the rendered manifests against it. Defaults to the target branch of a merged results pipeline, else the merge base of the merge request.",
			Required:    false,
			Value:       "",
			Destination: &P.DiffBase,
		},

		// GitLab sets both in every merge request pipeline, the target branch one empty
		// unless it is a merged results pipeline, and a source chain stops at the first
		// variable that is set, so they can not share one flag.
		&cli.StringFlag{
			Category:    CategoryKustomizeBuild,
			Name:        "kustomize.build.target-branch-sha",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("CI_MERGE_REQUEST_TARGET_BRANCH_SHA")),
			Usage:       "Head of the target branch a merged results pipeline merged into, which the report diffs against.",
			Required:    false,
			Value:       "",
			Destination: &P.TargetBranchSha,
		},

		&cli.StringFlag{
			Category:    CategoryKustomizeBuild,
			Name:        "kustomize.build.diff-base-sha",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("CI_MERGE_REQUEST_DIFF_BASE_SHA")),
			Usage:       "Merge base of the merge request, which the report diffs against in any other merge request pipeline.",
			Required:    false,
			Value:       "",
			Destination: &P.DiffBaseSha,
		},
	},
	gitlab.NewFlags(gitlab.Options{Destination: &P.MergeRequestReport}),
	ci.NewFlags(ci.Options{Destination: &P.ReportMetadata}),
	terraform.NewFlags(terraform.Options{Destination: &P.ReportLog}),
)
