package utils

import (
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

var vars map[string]string

func ApplyEnvironmentTemplate(raw string) (string, error) {
	if len(raw) == 0 {
		return raw, nil
	}

	initEnvironmentVariables()

	result, err := InlineTemplate(raw, vars)

	if err != nil {
		return "", err
	}

	return result, nil
}

func ApplyEnvironmentTemplates(raw []string) ([]string, error) {
	if len(raw) == 0 {
		return raw, nil
	}

	initEnvironmentVariables()

	templated := []string{}

	for _, v := range raw {
		result, err := InlineTemplate(v, vars)

		if err != nil {
			return nil, err
		}

		templated = append(templated, result)
	}

	return templated, nil
}

func initEnvironmentVariables() {
	if vars == nil {
		vars = ParseEnvironmentVariablesToMap()
	}
}
