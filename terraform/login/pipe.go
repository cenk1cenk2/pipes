package login

import (
	"encoding/json"
	"fmt"

	. "github.com/cenk1cenk2/plumber/v6"
)

type (
	Registry struct {
		Credentials []TerraformRegistryCredentialsJson
	}

	Pipe struct {
		Registry
	}
)

var TL = TaskList{}

var P = &Pipe{}

func New(p *Plumber) *TaskList {
	return TL.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(tl *TaskList) error {
			if v := p.Cli.String("terraform.registry.credentials"); v != "" {
				if err := json.Unmarshal([]byte(v), &P.Registry.Credentials); err != nil {
					return fmt.Errorf("Can not unmarshal registry credentials: %w", err)
				}
			}

			if err := p.Validate(P); err != nil {
				return err
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				GenerateTerraformRegistryCredentialsEnvVars(tl).Job(),
			)
		})
}
