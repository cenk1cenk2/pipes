package login

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name: "terraform-registry.credentials",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TF_REGISTRY_CREDENTIALS"),
		),
		Usage:    "Terraform registry credentials. json([]struct { registry: string, token: string })",
		Required: false,
	},
}
