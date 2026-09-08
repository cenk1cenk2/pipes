package flags

import (
	json "encoding/json/v2"
	"fmt"

	"github.com/urfave/cli/v3"
	"go.yaml.in/yaml/v4"
)

// JSONFlag makes the flag unmarshal its value into dst as part of validation.
// Unknown members are rejected, so a misspelled key fails the flag.
func JSONFlag[T any](flag *cli.StringFlag, dst *T) *cli.StringFlag {
	return unmarshalFlag(flag, dst, func(data []byte, v any) error {
		return json.Unmarshal(data, v, json.RejectUnknownMembers(true))
	})
}

// YAMLFlag is JSONFlag for the flags documented as YAML. JSON parses as YAML, so
// the two only differ in the usage text.
func YAMLFlag[T any](flag *cli.StringFlag, dst *T) *cli.StringFlag {
	return unmarshalFlag(flag, dst, yaml.Unmarshal)
}

func unmarshalFlag[T any](flag *cli.StringFlag, dst *T, unmarshal func([]byte, any) error) *cli.StringFlag {
	flag.ValidateDefaults = true
	flag.Validator = func(v string) error {
		// most of these flags are optional, so an unset one leaves the destination alone.
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
