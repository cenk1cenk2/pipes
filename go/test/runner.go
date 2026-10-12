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
		args := []string{"test"}
		if P.JUnit != "" {
			// go test has no JUnit writer of its own, so gotestsum runs it and writes one.
			args = []string{"--junitfile=" + P.JUnit, "--format=standard-verbose", "--"}
		}

		if P.Coverage != "" {
			args = append(args, "-coverpkg=./...", "-coverprofile="+P.Coverage)
		}

		args = append(args, "-v", "-failfast")

		if P.Args != "" {
			args = append(args, strings.Split(P.Args, " ")...)
		}

		args = append(args, "./...")

		if P.JUnit != "" {
			return "gotestsum", args
		}

		return "go", args
	case RunnerGinkgo:
		fallthrough
	default:
		args := []string{"tool", "ginkgo"}
		if P.JUnit != "" {
			args = append(args, "--junit-report="+P.JUnit)
		}

		if P.Coverage != "" {
			args = append(args, "--cover", "--coverpkg=./...", "--coverprofile="+P.Coverage)
		}

		args = append(args, "-vv", "--silence-skips", "--fail-fast", "-r")

		if P.Args != "" {
			args = append(args, strings.Split(P.Args, " ")...)
		}

		return "go", append(args, "./...")
	}
}
