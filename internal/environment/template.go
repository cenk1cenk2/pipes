package environment

// What the flags rendered against Template promise in the generated
// documentation.
const HELP_FORMAT_TEMPLATE = "Template(struct { Environment: string, EnvVars: map[string]string })"

// Template is the context the flags that accept one are rendered against.
type Template struct {
	Environment string
	EnvVars     map[string]string
}
