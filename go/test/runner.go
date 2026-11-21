package test

import (
	"gitlab.kilic.dev/devops/pipes/internal/shell"
)

type Runner string

const (
	RunnerGinkgo Runner = "ginkgo"
	RunnerGo     Runner = "go"
)

func command() (string, []string, error) {
	extra, err := shell.Split(P.Args)
	if err != nil {
		return "", nil, err
	}

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

		args = append(args, extra...)

		args = append(args, "./...")

		if P.JUnit != "" {
			return "gotestsum", args, nil
		}

		return "go", args, nil
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

		args = append(args, extra...)

		return "go", append(args, "./..."), nil
	}
}
