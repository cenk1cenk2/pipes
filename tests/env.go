package tests

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// envSource matches one environment variable of a flag row of a generated
// README, which lists them in the order the value source chain resolves them.
var envSource = regexp.MustCompile("`\\$([A-Z0-9_]+)`")

// ReadEnvChains reads the environment source chain of every flag of a pipe that
// answers to more than one name, keyed by the canonical name the chain starts
// with and carrying the fallbacks that follow it in order.
//
// The pipes are package main, so the README is the only place the chains can be
// read from outside the pipe itself. It is generated from the flags and the
// pipeline fails on a stale one, so the order it prints is the order a flag
// resolves its value in.
func ReadEnvChains(dir string) (map[string][]string, error) {
	contents, err := os.ReadFile(filepath.Join(Root(), dir, "README.md"))
	if err != nil {
		return nil, err
	}

	chains := map[string][]string{}

	for _, line := range strings.Split(string(contents), "\n") {
		if !strings.HasPrefix(line, "| `$") {
			continue
		}

		names := []string{}

		// Only the first cell holds the names, and the default of a flag may
		// carry a pipe of its own.
		for _, match := range envSource.FindAllStringSubmatch(strings.Split(line, "|")[1], -1) {
			names = append(names, match[1])
		}

		if len(names) < 2 {
			continue
		}

		chains[names[0]] = names[1:]
	}

	return chains, nil
}
