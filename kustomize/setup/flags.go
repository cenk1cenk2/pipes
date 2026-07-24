package setup

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_KUSTOMIZE = "Kustomize"
)

var Flags = []cli.Flag{

	&cli.StringFlag{
		Category: CATEGORY_KUSTOMIZE,
		Name:     "kustomize.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("KUSTOMIZE_ROOT"),
		),
		Usage:       "Working directory for Kustomize commands.",
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},

	&cli.StringSliceFlag{
		Category: CATEGORY_KUSTOMIZE,
		Name:     "kustomize.paths",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("KUSTOMIZE_PATHS"),
		),
		Usage:       "Explicit overlay paths to build relative to the working directory. When set, overlay discovery is skipped.",
		Required:    false,
		Destination: &P.Paths,
	},

	&cli.StringSliceFlag{
		Category: CATEGORY_KUSTOMIZE,
		Name:     "kustomize.discovery-pattern",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("KUSTOMIZE_DISCOVERY_PATTERN"),
		),
		Usage:    "Glob patterns to discover Kustomize overlays under the working directory. format(glob)",
		Required: false,
		Value: []string{
			"**/kustomization.yaml",
			"**/kustomization.yml",
			"**/Kustomization",
		},
		Destination: &P.DiscoveryPattern,
	},
}
