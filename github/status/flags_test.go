package status

import (
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Flags", func() {
	chains := map[string][]string{
		"status.token":      {"GITHUB_STATUS_TOKEN", "GH_TOKEN"},
		"status.project":    {"GITHUB_STATUS_PROJECT"},
		"status.state":      {"GITHUB_STATUS_STATE", "GITHUB_STATUS_REPORT"},
		"status.sha":        {"GITHUB_STATUS_SHA", "CI_COMMIT_SHA"},
		"status.target-url": {"GITHUB_STATUS_TARGET_URL", "CI_PIPELINE_URL"},
	}

	// the lookup is asserted on the sources, since a package level flag reads its
	// environment only on the first parse.
	lookup := func(name string, environment map[string]string) string {
		GinkgoHelper()

		tests.WithoutEnvironment(chains[name]...)
		tests.WithEnvironment(environment)

		sources := fixtures.Flag[*cli.StringFlag](Flags, name).Sources
		value, found := sources.Lookup()
		Expect(found).To(BeTrue())

		return value
	}

	It("documents every name of a chain", func() {
		for name, chain := range chains {
			Expect(fixtures.Flag[*cli.StringFlag](Flags, name).GetEnvVars()).To(Equal(chain), name)
		}
	})

	// the names the gh-status reference exports have to keep working for the
	// pipelines that still set them.
	DescribeTable(
		"reads the names the previous status reference used",
		func(name, variable string) {
			Expect(lookup(name, map[string]string{variable: "legacy"})).To(Equal("legacy"))
		},
		Entry("token", "status.token", "GH_TOKEN"),
		Entry("project", "status.project", "GITHUB_STATUS_PROJECT"),
		Entry("state", "status.state", "GITHUB_STATUS_REPORT"),
		Entry("sha", "status.sha", "CI_COMMIT_SHA"),
		Entry("target url", "status.target-url", "CI_PIPELINE_URL"),
	)

	DescribeTable(
		"prefers the name of the pipe over the one it falls back to",
		func(name string) {
			chain := chains[name]

			Expect(lookup(name, map[string]string{chain[0]: "preferred", chain[1]: "fallback"})).To(Equal("preferred"))
		},
		Entry("token", "status.token"),
		Entry("state", "status.state"),
		Entry("sha", "status.sha"),
		Entry("target url", "status.target-url"),
	)
})
