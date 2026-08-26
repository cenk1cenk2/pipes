package setup

import (
	. "github.com/cenk1cenk2/plumber/v6"
	ucli "github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/git"
)

//revive:disable:line-length-limit

const (
	CATEGORY_ENVIRONMENT = "Environment"

	DEFAULT_ENVIRONMENT_CONDITIONS = `[
    { "match": "^tags/v?\\d+.\\d+.\\d+$", "environment": "production" },
    { "match": "^tags/v?\\d+.\\d+.\\d+-.*\\.\\d+$", "environment": "stage" },
    { "match" :"^heads/main$", "environment": "develop" },
    { "match": "^heads/master$", "environment": "develop" }
]`
)

var Flags = CombineFlags(
	git.NewFlags(&P.Git),
	[]ucli.Flag{
		&ucli.BoolFlag{
			Category:    CATEGORY_ENVIRONMENT,
			Name:        "environment.enable",
			Sources:     cli.EnvVars("ENVIRONMENT_ENABLE"),
			Usage:       "Enable environment injection.",
			Required:    false,
			Hidden:      true,
			Value:       true,
			Destination: &P.Environment.Enable,
		},

		// CATEGORY_ENVIRONMENT

		cli.JSONFlag(&ucli.StringFlag{
			Category: CATEGORY_ENVIRONMENT,
			Name:     "environment.conditions",
			Sources:  cli.EnvVars("ENVIRONMENT_CONDITIONS"),
			Usage: `Regex pattern to select an environment.
      Use either "heads/" for narrowing the search to branches or "tags/" for narrowing the search to tags.
      json([]struct{ match: RegExp, environment: string })`,
			Required: false,
			Value:    DEFAULT_ENVIRONMENT_CONDITIONS,
		}, &P.Environment.Conditions),

		&ucli.BoolFlag{
			Category:    CATEGORY_ENVIRONMENT,
			Name:        "environment.fail-on-no-reference",
			Sources:     cli.EnvVars("ENVIRONMENT_FAIL_ON_NO_REFERENCE"),
			Usage:       "Fail on missing environment references.",
			Required:    false,
			Value:       true,
			Destination: &P.Environment.FailOnNoReference,
		},

		&ucli.BoolFlag{
			Category:    CATEGORY_ENVIRONMENT,
			Name:        "environment.strict",
			Sources:     cli.EnvVars("ENVIRONMENT_STRICT"),
			Usage:       "Fail on no environment selected.",
			Required:    false,
			Value:       true,
			Destination: &P.Environment.Strict,
		},
	})
