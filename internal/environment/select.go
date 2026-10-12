package environment

import (
	"fmt"
	"slices"
	"strings"

	"gitlab.kilic.dev/devops/pipes/internal/git"
)

//revive:disable:line-length-limit

// The conditions a pipe falls back to; the value is printed verbatim as the flag default in the generated documentation.
const DefaultConditions = `[
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

// Fetch returns only the variables in environ that carry the environment prefix, with
// the prefix stripped, so STAGE_TOKEN reaches the pipe as TOKEN once stage is selected.
// ENVIRONMENT names the selection itself.
func Fetch(environ []string, environment string) map[string]string {
	prefix := prefix(environment)

	vars := map[string]string{}

	for _, v := range environ {
		key, value, found := strings.Cut(v, "=")

		if !found {
			continue
		}

		if trimmed, ok := strings.CutPrefix(key, prefix); ok {
			vars[trimmed] = value
		}
	}

	vars["ENVIRONMENT"] = environment

	return vars
}

// Selected names the variables the environment brings in, as the pipe sees them once
// their prefix is stripped.
func Selected(environ []string, environment string) []string {
	prefix := prefix(environment)

	names := []string{}
	for _, v := range environ {
		key, _, found := strings.Cut(v, "=")
		if !found {
			continue
		}

		if name, ok := strings.CutPrefix(key, prefix); ok {
			names = append(names, name)
		}
	}
	slices.Sort(names)

	return names
}

func prefix(environment string) string {
	return strings.ToUpper(environment) + "_"
}
