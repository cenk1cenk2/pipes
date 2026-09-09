package test

import (
	"strings"
)

type Runner string

const (
	RunnerGinkgo Runner = "ginkgo"
	RunnerGo     Runner = "go"
)

func command() (string, []string) {
	switch Runner(P.Runner) {
	case RunnerGo:
		return "go", arguments([]string{"test"}, []string{"-coverpkg=./...", "-coverprofile=" + P.Coverage}, []string{"-v", "-failfast"})
	case RunnerGinkgo:
		fallthrough
	default:
		return "go", arguments(
			[]string{"tool", "ginkgo"},
			[]string{"--cover", "--coverpkg=./...", "--coverprofile=" + P.Coverage},
			[]string{"-vv", "--silence-skips", "--fail-fast", "-r"},
		)
	}
}

func arguments(prefix, coverage, always []string) []string {
	args := prefix

	if P.Coverage != "" {
		args = append(args, coverage...)
	}

	args = append(args, always...)

	if P.Args != "" {
		args = append(args, strings.Split(P.Args, " ")...)
	}

	return append(args, "./...")
}
