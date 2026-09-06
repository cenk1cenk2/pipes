package tool

import (
	"fmt"

	ucli "github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

// Flags builds the working directory flag shared by every tool pipe.
func Flags(spec Spec, cfg *Config) []ucli.Flag {
	return []ucli.Flag{
		&ucli.StringFlag{
			Category:    spec.Category,
			Name:        spec.FlagPrefix + ".cwd",
			Sources:     cli.EnvVars(spec.EnvPrefix + "_CWD"),
			Usage:       fmt.Sprintf("Working directory for %s commands.", spec.Name),
			Required:    false,
			Value:       ".",
			Destination: &cfg.Cwd,
		},
	}
}
