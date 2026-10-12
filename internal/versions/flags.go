package versions

import (
	"github.com/urfave/cli/v3"
)

//revive:disable:line-length-limit

const (
	CategoryTagsFromEnvironment = "Tags From Environment"
)

// EnvironmentOptions is what the flag of the environment source is built onto, so
// every pipe that publishes under a version reads it under the same name.
type EnvironmentOptions struct {
	Destination *[]string
}

func NewEnvironmentFlags(opts EnvironmentOptions) []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Category:    CategoryTagsFromEnvironment,
			Name:        "tags-from-environment",
			Sources:     cli.NewValueSourceChain(cli.EnvVar("TAGS_FROM_ENVIRONMENT")),
			Usage:       `Environment variables to read tags from, where a name ending in "*" reads every variable under that prefix. An unset or empty variable is skipped.`,
			Required:    false,
			Destination: opts.Destination,
		},
	}
}
