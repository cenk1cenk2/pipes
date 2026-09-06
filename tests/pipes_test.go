package tests

import (
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// The conformance specs of a pipe run inside its own package main, so nothing
// here can import them. What this module can still hold every pipe to is that the
// suite is there at all: a pipe whose main package stopped exposing TestPipe
// stopped being checked, and nothing else would say so.
var _ = Describe("Pipes", func() {
	for _, dir := range Pipes() {
		Describe(dir, func() {
			It("exposes its conformance suite as TestPipe", func() {
				// The test binary is only listed, never run, so this costs a compile
				// and nothing else.
				command := exec.Command("go", "test", "-list", "^TestPipe$", ".")
				command.Dir = filepath.Join(Root(), dir)

				out, err := command.CombinedOutput()
				Expect(err).NotTo(HaveOccurred(), string(out))
				Expect(string(out)).To(
					ContainSubstring("TestPipe"),
					"the pipe runs no conformance suite, add a main_test.go that calls conformance.Verify from func TestPipe",
				)
			})
		})
	}
})
