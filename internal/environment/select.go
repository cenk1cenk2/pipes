package environment

import (
	"fmt"
	"strings"

	"gitlab.kilic.dev/devops/pipes/internal/git"
)

//revive:disable:line-length-limit

// The conditions a pipe falls back to when the user did not write its own. The
// value is what the generated documentation prints as the flag default, so it
// stays as it is until a step deliberately rewrites the docs.
const DEFAULT_CONDITIONS = `[
    { "match": "^tags/v?\\d+.\\d+.\\d+$", "environment": "production" },
    { "match": "^tags/v?\\d+.\\d+.\\d+-.*\\.\\d+$", "environment": "stage" },
    { "match" :"^heads/main$", "environment": "develop" },
    { "match": "^heads/master$", "environment": "develop" }
]`

// Condition pairs a source control reference pattern with the environment it
// selects.
type Condition struct {
	Match       string `json:"match"       validate:"required"`
	Environment string `json:"environment" validate:"required"`
}

// Select reports the environment of the first condition matching any of the
// references, or an empty string when none of them do.
func Select(conditions []Condition, references []string) (string, error) {
	matches := make([]string, 0, len(conditions))

	for _, condition := range conditions {
		matches = append(matches, condition.Match)
	}

	matched, err := git.MatchAny(matches, references)

	if err != nil {
		return "", fmt.Errorf("Can not process regular expression for environment: %w", err)
	}

	if matched < 0 {
		return "", nil
	}

	return conditions[matched].Environment, nil
}

// Fetch strips the environment prefix off every variable in environ, so that
// STAGE_TOKEN reaches the pipe as TOKEN once stage is the selected environment.
// Variables that do not carry the prefix are kept as they are, which is what
// lets an environment override only the few it cares about, and ENVIRONMENT
// names the selection itself.
func Fetch(environ []string, environment string) map[string]string {
	prefix := strings.ToUpper(environment) + "_"

	vars := make(map[string]string, len(environ)+1)

	for _, v := range environ {
		key, value, found := strings.Cut(v, "=")

		if !found {
			continue
		}

		trimmed, _ := strings.CutPrefix(key, prefix)

		vars[trimmed] = value
	}

	vars["ENVIRONMENT"] = environment

	return vars
}
