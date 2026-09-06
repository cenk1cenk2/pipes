package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/test/conformance"
	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
)

func TestPipe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Helm Pipe Suite")
}

var _ = conformance.Verify(conformance.Pipe{
	Name:        name,
	Description: description,
	New: func(p *plumber.Plumber) *ucli.Command {
		return newCommand(p, "test", options{})
	},
	UncategorizedFlags: []string{
		"helm.lint.kubernetes.version",
		"helm.lint.should-template",
	},
	LegacyEnvAliases: map[string][]string{
		"git.branch":                           {"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
		"git.tag":                              {"CI_COMMIT_TAG", "BITBUCKET_TAG"},
		"helm.cwd":                             {"HELM_ROOT", "HELM_CWD"},
		"helm.lint.kubernetes.version":         {"KUBERNETES_VERSION", "HELM_LINT_KUBERNETES_VERSION"},
		"helm.login.registry.password":         {"HELM_REGISTRY_PASSWORD", "HELM_LOGIN_REGISTRY_PASSWORD"},
		"helm.login.registry.uri":              {"HELM_REGISTRY_URI", "HELM_LOGIN_REGISTRY_URI"},
		"helm.login.registry.username":         {"HELM_REGISTRY_USERNAME", "HELM_LOGIN_REGISTRY_USERNAME"},
		"helm.publish.chart.app-version":       {"HELM_CHART_APP_VERSION", "HELM_PUBLISH_CHART_APP_VERSION"},
		"helm.publish.chart.destination":       {"HELM_CHART_DESTINATION", "HELM_PUBLISH_CHART_DESTINATION"},
		"helm.publish.chart.target":            {"HELM_CHART_TARGET", "HELM_PUBLISH_CHART_TARGET"},
		"helm.publish.chart.versions":          {"HELM_CHART_VERSIONS", "HELM_PUBLISH_CHART_VERSIONS"},
		"helm.publish.chart.versions-sanitize": {"HELM_CHART_SANITIZE_VERSIONS", "HELM_PUBLISH_CHART_VERSIONS_SANITIZE"},
		"helm.publish.chart.versions-template": {"HELM_CHART_VERSIONS_TEMPLATE", "HELM_PUBLISH_CHART_VERSIONS_TEMPLATE"},
	},
})

// Everything a spec needs is passed as an argument rather than through the
// environment, since the flags are package level and urfave only reads an env
// source on the first parse of a flag instance: driving values in through the
// environment would make the specs depend on the order Ginkgo runs them in.
func run(runner *tests.TestingCommandRunner, args ...string) error {
	GinkgoHelper()

	fixture := fixtures.NewPlumber(func(p *plumber.Plumber) *ucli.Command {
		return newCommand(p, "test", options{})
	})
	fixture.Plumber.SetRuntime(plumber.Runtime{CommandRunner: runner.Runner()})

	return fixture.RunCli(append([]string{name}, args...)...)
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

// The setup step reads the chart before anything else runs, so every command
// needs one on disk.
func chart() string {
	GinkgoHelper()

	dir := tests.TempDir()
	Expect(os.WriteFile(
		filepath.Join(dir, "Chart.yaml"),
		[]byte("apiVersion: v2\nname: example\nversion: 0.1.0\n"),
		0o644,
	)).To(Succeed())

	return dir
}

var _ = Describe("New", func() {
	It("updates the chart dependencies on install", func() {
		runner := fixtures.Runner()

		Expect(run(runner, "install", "--helm.cwd", chart())).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"helm", "helm"}))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("dependency update")))
	})

	It("lints the chart on lint", func() {
		runner := fixtures.Runner()

		Expect(run(runner, "lint", "--helm.cwd", chart())).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"helm", "helm", "helm"}))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("lint .")))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("template .")))
	})

	It("packages and pushes the chart on publish", func() {
		runner := fixtures.Runner()

		Expect(run(
			runner,
			"publish",
			"--helm.cwd", chart(),
			"--helm.publish.chart.target", "oci://registry.example.test/charts",
			"--helm.publish.chart.versions", "1.0.0",
		)).To(Succeed())

		Expect(formatted(runner)).To(ContainElement(ContainSubstring("package -d ./dist/ . --version 1.0.0")))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("push dist/example-1.0.0.tgz oci://registry.example.test/charts")))
	})
})
