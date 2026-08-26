package cli

import (
	"encoding/json"
	"fmt"

	ucli "github.com/urfave/cli/v3"
	"go.yaml.in/yaml/v4"
)

// JSONFlag makes the flag unmarshal its value into dst as part of validation, so
// the pipe reads a struct where the user wrote a JSON string.
func JSONFlag[T any](flag *ucli.StringFlag, dst *T) *ucli.StringFlag {
	return unmarshalFlag(flag, dst, json.Unmarshal)
}

// YAMLFlag is JSONFlag for the flags documented as YAML. JSON parses as YAML, so
// the two only differ in what the usage text promises.
func YAMLFlag[T any](flag *ucli.StringFlag, dst *T) *ucli.StringFlag {
	return unmarshalFlag(flag, dst, yaml.Unmarshal)
}

func unmarshalFlag[T any](flag *ucli.StringFlag, dst *T, unmarshal func([]byte, any) error) *ucli.StringFlag {
	flag.ValidateDefaults = true
	flag.Validator = func(v string) error {
		// An unset flag leaves the destination at its zero value rather than
		// failing, since most of these are optional.
		if v == "" {
			return nil
		}

		if err := unmarshal([]byte(v), dst); err != nil {
			return fmt.Errorf("Can not unmarshal %s: %w", flag.Name, err)
		}

		return nil
	}

	return flag
}

// EnvVars builds a value source chain out of environment variable names. Legacy
// names come first so a variable already set in a pipeline keeps winning over the
// name that replaced it.
func EnvVars(names ...string) ucli.ValueSourceChain {
	sources := make([]ucli.ValueSource, 0, len(names))

	for _, name := range names {
		sources = append(sources, ucli.EnvVar(name))
	}

	return ucli.NewValueSourceChain(sources...)
}
