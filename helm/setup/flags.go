package setup

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CATEGORY_HELM = "Helm"
)

var Flags = []cli.Flag{

	&cli.StringFlag{
		Category: CATEGORY_HELM,
		Name:     "helm.cwd",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("HELM_CWD"),
		),
		Usage:       "Working directory for Helm commands.",
		Required:    false,
		Value:       ".",
		Destination: &P.Cwd,
	},
}
