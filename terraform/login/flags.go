package login

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

const CategoryLogin = "Login"

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	flags.JSONFlag(&P.Registry.Credentials, &cli.StringFlag{
		Category: CategoryLogin,
		Name:     "terraform.login.registry.credentials",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_LOGIN_REGISTRY_CREDENTIALS"),
			cli.EnvVar("TF_REGISTRY_CREDENTIALS"),
		),
		Usage:    "Terraform registry credentials. format(json([]struct{ registry: string, token: string }))",
		Required: false,
	}),
}
