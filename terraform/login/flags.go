package login

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

const CATEGORY_LOGIN = "Login"

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	flags.JSONFlag(&cli.StringFlag{
		Category: CATEGORY_LOGIN,
		Name:     "terraform.login.registry.credentials",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_LOGIN_REGISTRY_CREDENTIALS"),
			cli.EnvVar("TF_REGISTRY_CREDENTIALS"),
		),
		Usage:    "Terraform registry credentials. json([]struct { registry: string, token: string })",
		Required: false,
	}, &P.Registry.Credentials),
}
