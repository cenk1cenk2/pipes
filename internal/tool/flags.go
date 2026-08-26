package tool

import (
	"fmt"
	"slices"

	ucli "github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

// Flags builds the working directory flag shared by every tool pipe. The legacy
// aliases come first in the chain, so the canonical <ENVPREFIX>_CWD only applies
// where a pipeline has not already set the name it replaced.
func Flags(spec Spec, cfg *Config) []ucli.Flag {
	envs := append(slices.Clone(spec.CwdEnvAliases), spec.EnvPrefix+"_CWD")

	return []ucli.Flag{
		&ucli.StringFlag{
			Category:    spec.Category,
			Name:        spec.FlagPrefix + ".cwd",
			Sources:     cli.EnvVars(envs...),
			Usage:       fmt.Sprintf("Working directory for %s commands.", spec.Name),
			Required:    false,
			Value:       ".",
			Destination: &cfg.Cwd,
		},
	}
}
