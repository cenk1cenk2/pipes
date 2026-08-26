package app_test

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/buildah/app"
	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
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

	return fixture.RunCli(append([]string{"pipe-buildah"}, args...)...)
}

func formatted(runner *tests.TestingCommandRunner) []string {
	GinkgoHelper()

	invocations := runner.Invocations()
	commands := make([]string, len(invocations))

	for i, invocation := range invocations {
		commands[i] = invocation.Formatted
	}

	return commands
}

var _ = Describe("New", func() {
	It("probes the tool and authenticates against the registry on login", func() {
		runner := fixtures.Runner()

		Expect(run(
			runner,
			"login",
			"--buildah.login.registry.uri", "registry.example.test",
			"--buildah.login.registry.username", "user",
			"--buildah.login.registry.password", "password",
		)).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"buildah", "buildah"}))
		// The password goes in over stdin, so it is deliberately absent from the
		// arguments the invocation records.
		Expect(formatted(runner)).To(ContainElement(
			ContainSubstring("login registry.example.test --username user --password-stdin"),
		))
	})

	// Without a username and a password the login task disables itself, which is
	// what a pipeline relying on an ambient login does. The registry the tags are
	// prefixed with is the default the login step declares.
	It("builds the image under the tags it was given on build", func() {
		tests.WithTempWorkingDirectory()

		runner := fixtures.Runner()

		Expect(run(
			runner,
			"build",
			"--buildah.build.image.name", "group/image",
			"--buildah.build.image.tags", "v1.0.0",
			"--buildah.build.image.push=false",
		)).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"buildah", "buildah"}))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("build --format oci --pull -t docker.io/group/image:v1.0.0 --file Dockerfile .")))
	})

	It("creates and pushes the manifest of the images it was given on manifest", func() {
		tests.WithTempWorkingDirectory()

		runner := fixtures.Runner()

		Expect(run(
			runner,
			"manifest",
			"--buildah.manifest.target", "group/image:v1.0.0",
			"--buildah.manifest.images", "group/image:v1.0.0-amd64",
		)).To(Succeed())

		Expect(formatted(runner)).To(ContainElement(ContainSubstring("manifest create docker.io/group/image:v1.0.0")))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("manifest add docker.io/group/image:v1.0.0 group/image:v1.0.0-amd64")))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("manifest push --rm docker.io/group/image:v1.0.0")))
	})
})
