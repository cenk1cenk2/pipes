// Package fixtures wires the plumber test kit the way the pipes need it, so a
// pipe spec has one import to reach for instead of assembling the runtime itself.
package fixtures

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
)

// NewPlumber builds a plumber wired to the Ginkgo writer, with no greeter and no
// global logger left behind for the next spec.
func NewPlumber(constructors ...plumber.PlumberNewFn) *tests.PlumberFixture {
	return tests.NewPlumber(constructors...)
}
