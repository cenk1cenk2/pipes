package tool

import (
	"fmt"

	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

// NewFlags builds the working directory flag shared by every tool pipe.
func NewFlags(spec Spec, cfg *Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Category:    spec.Category,
			Name:        spec.FlagPrefix + ".cwd",
			Sources:     flags.EnvVars(spec.EnvPrefix + "_CWD"),
			Usage:       fmt.Sprintf("Working directory for %s commands.", spec.Name),
			Required:    false,
			Value:       ".",
			Destination: &cfg.Cwd,
		},
	}
}
