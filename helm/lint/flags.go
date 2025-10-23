package lint

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "kubernetes.version",
		Usage: "Kubernetes version to use for linting charts.",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("KUBERNETES_VERSION"),
		),
		Required:    false,
		Value:       "",
		Destination: &P.Kubernetes.Version,
	},
}
