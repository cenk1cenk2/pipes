package write

import (
	"os"
	"path/filepath"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	"gitlab.kilic.dev/devops/pipes/select-env/setup"
)

var _ = Describe("Environment file", func() {
	// The flags are not registered with the spec command, since the pipe is seeded
	// directly and a package level flag only reads its environment on first parse.
	// The task reads the selection of the pipe around it off its package level
	// instance, so a spec seeds that the same way it seeds its own.
	run := func(pipe Pipe, selected environment.Ctx) error {
		GinkgoHelper()

		*P = pipe
		*setup.EnvironmentCtx = selected

		return fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
			AppName:     "select-env",
			CommandName: "write",
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *cli.Command) *plumber.TaskList {
					tl := &plumber.TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *plumber.TaskList) plumber.Job {
							return plumber.JobSequence(WriteEnvironmentFile(tl).Job())
						})
				},
			},
		}).Run()
	}

	var file string

	BeforeEach(func() {
		file = filepath.Join(GinkgoT().TempDir(), "env.environment")
	})

	It("writes the variables of the selected environment as a dotenv file", func() {
		Expect(run(
			Pipe{Environment: Environment{File: file}},
			environment.Ctx{
				Environment: "production",
				EnvVars:     map[string]string{"API_URL": "https://api.example.com"},
			},
		)).To(Succeed())

		Expect(os.ReadFile(file)).To(BeEquivalentTo("API_URL=\"https://api.example.com\"\n"))
	})

	// A pipeline sources the file unconditionally, so it has to exist even when
	// the selection resolved to nothing.
	It("writes the file even when the selection carries no variables", func() {
		Expect(run(
			Pipe{Environment: Environment{File: file}},
			environment.Ctx{},
		)).To(Succeed())

		Expect(os.ReadFile(file)).To(BeEquivalentTo("\n"))
	})
})
