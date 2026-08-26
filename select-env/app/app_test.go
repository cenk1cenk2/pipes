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
	"gitlab.kilic.dev/devops/pipes/select-env/app"
)

// The reference and the output file are passed as arguments rather than through
// the environment, since the flags are package level and urfave only reads an
// env source on the first parse of a flag instance. The variables the pipe
// selects between are not flags, so those still go in through the environment.
func run(args ...string) error {
	GinkgoHelper()

	fixture := fixtures.NewPlumber(func(p *plumber.Plumber) *ucli.Command {
		return app.New(p, "test", app.Options{})
	})
	fixture.Plumber.SetRuntime(plumber.Runtime{CommandRunner: fixtures.Runner().Runner()})

	return fixture.RunCli(append([]string{"select-env"}, args...)...)
}

var _ = Describe("New", func() {
	It("writes the variables of the environment the reference selected", func() {
		file := filepath.Join(tests.TempDir(), "env.environment")

		tests.WithEnvironment(map[string]string{"DEVELOP_DATABASE_URL": "postgres://develop"})

		Expect(run("--git.branch", "main", "--environment.file", file)).To(Succeed())

		// The file is dotenv, so the values are written quoted, and the prefix that
		// named the environment is stripped off the variable it selected.
		content, err := os.ReadFile(file)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(content)).To(ContainSubstring(`DATABASE_URL="postgres://develop"`))
		Expect(string(content)).To(ContainSubstring(`ENVIRONMENT="develop"`))
	})
})
