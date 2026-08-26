package tool

import (
	"regexp"
	"strings"
)

// Spec describes the command line tool a pipe drives. Every field is fixed for a
// given pipe, so it is declared next to that pipe's flags rather than configured.
type Spec struct {
	// Name is the binary the pipe runs, and the name the documentation calls it by.
	Name string
	// Category is the heading the generated documentation files the flag under. It
	// stays whatever the pipe already printed, since the docs are a published page.
	Category   string
	FlagPrefix string
	EnvPrefix  string
	// CwdEnvAliases are the environment variable names the working directory
	// answered to before the pipes agreed on one. They are kept forever and listed
	// ahead of the canonical name, so a pipeline that already sets one keeps
	// winning over the name that replaced it.
	CwdEnvAliases []string
	VersionArgs   []string
	// VersionPattern narrows the probe output down to its first submatch. Without
	// one the whole output is reported.
	VersionPattern *regexp.Regexp
}

// Config is the working directory flag destination every tool pipe shares.
type Config struct {
	Cwd string `validate:"omitempty,dir"`
}

// Ctx is what Setup resolves. A pipe holds on to the same instance it handed to
// Setup, since the values only land once that task list has run.
type Ctx struct {
	Cwd     string
	Version string
	Env     map[string]string
}

func NewCtx() *Ctx {
	return &Ctx{Env: map[string]string{}}
}

// ParseVersion pulls the version out of the probe output. Output that does not
// match is reported whole rather than treated as an error, since the probe only
// ever feeds a log line and the tool has already proven it runs.
func (s Spec) ParseVersion(output string) string {
	output = strings.TrimSpace(output)

	if s.VersionPattern == nil {
		return output
	}

	matches := s.VersionPattern.FindStringSubmatch(output)

	if len(matches) < 2 {
		return output
	}

	return matches[1]
}
