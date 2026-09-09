package versions

// Match is a condition and the template that replaces whatever it matched. The
// submatches of the condition are the context the template renders against, so
// "{{ index $ 1 }}" is the first capture group.
type Match struct {
	Match    string `json:"match"    yaml:"match"    validate:"required"`
	Template string `json:"template" yaml:"template" validate:"required"`
}
