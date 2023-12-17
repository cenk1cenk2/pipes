package login

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:     "terraform-registry.credentials",
		Usage:    "Terraform registry credentials. json([]struct { registry: string, token: string })",
		Required: false,
		EnvVars:  []string{"TF_REGISTRY_CREDENTIALS"},
	},
}

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList[Pipe]) error {
	if v := tl.CliContext.String("terraform-registry.credentials"); v != "" {
		if err := json.Unmarshal([]byte(v), &tl.Pipe.Registry.Credentials); err != nil {
			return fmt.Errorf("Can not unmarshal registry credentials: %w", err)
		}
	}

	return nil
}
