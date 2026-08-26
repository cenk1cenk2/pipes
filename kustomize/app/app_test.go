package app_test

import (
	"os"
	"path/filepath"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	"gitlab.kilic.dev/devops/pipes/kustomize/app"
)

// Everything a spec needs is passed as an argument rather than through the
// environment, since the flags are package level and urfave only reads an env
// source on the first parse of a flag instance: driving values in through the
// environment would make the specs depend on the order Ginkgo runs them in.
func run(args ...string) error {
	GinkgoHelper()

	fixture := fixtures.NewPlumber(func(p *plumber.Plumber) *ucli.Command {
		return app.New(p, "test", app.Options{})
	})
	fixture.Plumber.SetRuntime(plumber.Runtime{CommandRunner: fixtures.Runner().Runner()})

	return fixture.RunCli(append([]string{"pipe-kustomize"}, args...)...)
}

// The overlay is rendered through the Kustomize API rather than by running the
// binary, so the spec needs one that actually builds.
func overlay() string {
	GinkgoHelper()

	dir := tests.TempDir()

	Expect(os.WriteFile(
		filepath.Join(dir, "kustomization.yaml"),
		[]byte("apiVersion: kustomize.config.k8s.io/v1beta1\nkind: Kustomization\nresources:\n  - configmap.yaml\n"),
		0o644,
	)).To(Succeed())

	Expect(os.WriteFile(
		filepath.Join(dir, "configmap.yaml"),
		[]byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: example\n"),
		0o644,
	)).To(Succeed())

	return dir
}

var _ = Describe("New", func() {
	It("renders the overlay it resolved on build", func() {
		cwd := overlay()
		output := tests.TempDir("rendered")

		Expect(run(
			"build",
			"--kustomize.cwd", cwd,
			"--kustomize-build.output", output,
		)).To(Succeed())

		rendered, err := os.ReadDir(output)
		Expect(err).NotTo(HaveOccurred())
		Expect(rendered).To(HaveLen(1))

		content, err := os.ReadFile(filepath.Join(output, rendered[0].Name()))
		Expect(err).NotTo(HaveOccurred())
		Expect(string(content)).To(ContainSubstring("kind: ConfigMap"))
	})
})
