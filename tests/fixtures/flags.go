package fixtures

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"
)

// Flag looks a flag up by any of its names, so a spec asserts on the one it means
// instead of on the position a constructor happened to put it in.
func Flag[T cli.Flag](flags []cli.Flag, name string) T {
	GinkgoHelper()

	for _, flag := range flags {
		for _, candidate := range flag.Names() {
			if candidate != name {
				continue
			}

			typed, ok := flag.(T)
			Expect(ok).To(BeTrue(), "%s is not the flag type the spec asked for", name)

			return typed
		}
	}

	Fail("no flag is registered under " + name)

	var zero T

	return zero
}

// FlagNames renders the names of every flag in order, for the specs that assert on
// where a constructor put one relative to another.
func FlagNames(flags []cli.Flag) []string {
	names := []string{}
	for _, flag := range flags {
		names = append(names, flag.Names()...)
	}

	return names
}
