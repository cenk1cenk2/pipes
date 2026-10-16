package login

import (
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

const CategoryLogin = "Login"

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	flags.YAMLFlag(&P.Logins, &cli.StringFlag{
		Category: CategoryLogin,
		Name:     "terraform.login",
		Sources:  cli.NewValueSourceChain(cli.EnvVar("TERRAFORM_LOGIN")),
		Usage:    "Terraform registries to login to. format(yaml([]struct{ uri: string, token: string }))",
		Required: false,
	}),
}
