package shell

import (
	"fmt"

	"github.com/mattn/go-shellwords"
)

// Split breaks a flag value into command arguments the way a shell would, so a quoted
// argument keeps its spaces and an empty value yields no arguments at all.
func Split(value string) ([]string, error) {
	args, err := shellwords.Parse(value)
	if err != nil {
		return nil, fmt.Errorf("Can not split the arguments %q: %w", value, err)
	}

	return args, nil
}
