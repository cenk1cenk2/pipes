package login

import (
	"context"
	"encoding/json"
	"fmt"

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
		Action: func(_ context.Context, _ *cli.Command, v string) error {
			if err := json.Unmarshal([]byte(v), &P.Registry.Credentials); err != nil {
				return fmt.Errorf("Can not unmarshal registry credentials: %w", err)
			}

			return nil
		},
	},
}
