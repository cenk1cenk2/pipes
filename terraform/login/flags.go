package login

import (
	json "encoding/json/v2"
	"fmt"

	"github.com/urfave/cli/v3"
)

const CATEGORY_LOGIN = "Login"

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category: CATEGORY_LOGIN,
		Name:     "terraform.login.registry.credentials",
		Sources: cli.NewValueSourceChain(
			cli.EnvVar("TERRAFORM_LOGIN_REGISTRY_CREDENTIALS"),
			cli.EnvVar("TF_REGISTRY_CREDENTIALS"),
		),
		Usage:            "Terraform registry credentials. json([]struct { registry: string, token: string })",
		Required:         false,
		ValidateDefaults: false,
		Validator: func(v string) error {
			if v == "" {
				return nil
			}

			if err := json.Unmarshal([]byte(v), &P.Registry.Credentials, json.RejectUnknownMembers(true)); err != nil {
				return fmt.Errorf("Can not unmarshal registry credentials: %w", err)
			}

			return nil
		},
	},
}
