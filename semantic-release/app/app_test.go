package app_test

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	"gitlab.kilic.dev/devops/pipes/semantic-release/app"
)

// Everything a spec needs is passed as an argument rather than through the
// environment, since the flags are package level and urfave only reads an env
// source on the first parse of a flag instance: driving values in through the
// environment would make the specs depend on the order Ginkgo runs them in.
func run(runner *tests.TestingCommandRunner, args ...string) error {
	GinkgoHelper()

	fixture := fixtures.NewPlumber(func(p *plumber.Plumber) *ucli.Command {
		return app.New(p, "test", app.Options{})
	})
	fixture.Plumber.SetRuntime(plumber.Runtime{CommandRunner: runner.Runner()})

	return fixture.RunCli(append([]string{"pipe-semantic-release"}, args...)...)
}

func environmentEnable(flags []ucli.Flag) *ucli.BoolFlag {
	GinkgoHelper()

	for _, flag := range flags {
		if converted, ok := flag.(*ucli.BoolFlag); ok && converted.Name == "environment.enable" {
			return converted
		}
	}

	return nil
}

var _ = Describe("New", func() {
	// The pipe has no subcommand: the steps run off the root command.
	It("releases through semantic-release", func() {
		tests.WithTempWorkingDirectory()

		runner := fixtures.Runner()

		Expect(run(runner)).To(Succeed())

		Expect(runner.InvocationNames()).To(ContainElement("semantic-release"))
	})

	It("releases through multi-semantic-release for a workspace", func() {
		tests.WithTempWorkingDirectory()

		runner := fixtures.Runner()

		Expect(run(runner, "--semantic-release.workspace")).To(Succeed())

		Expect(runner.InvocationNames()).To(ContainElement("multi-semantic-release"))
	})

	// The environment feature is hidden and on by default everywhere it is owned
	// by the pipe. Here it is opt-in, and app.New is what turns it around.
	It("offers the environment selection as an opt-in", func() {
		command := app.New(nil, "test", app.Options{})

		flag := environmentEnable(command.Flags)
		Expect(flag).NotTo(BeNil())
		Expect(flag.Hidden).To(BeFalse())
		Expect(flag.Value).To(BeFalse())
	})
})
