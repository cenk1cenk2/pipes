package build

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/node"
	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
)

var _ = Describe("Node build", func() {
	// The flags are not registered with the spec command, since the pipe is seeded
	// directly and a package level flag only reads its environment on first parse.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe, deps Deps) error {
		GinkgoHelper()

		*P = pipe

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-node",
			CommandName: "build",
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *ucli.Command) *plumber.TaskList {
					tl := &plumber.TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *plumber.TaskList) plumber.Job {
							return plumber.JobSequence(BuildNodeApplication(tl, deps).Job())
						})
				},
			},
		}).Run()
	}

	deps := func(packageManager string) Deps {
		return Deps{
			Node: &node.Ctx{PackageManager: node.PackageManager{
				Exe:      packageManager,
				Commands: node.PackageManagers[packageManager],
			}},
			Environment: &environment.Ctx{
				Environment: "production",
				EnvVars:     map[string]string{"API_URL": "https://api.example.com"},
			},
		}
	}

	pipe := func() Pipe {
		return Pipe{Build: Build{Script: "build", Cwd: "projects/web"}}
	}

	It("runs the build script through the resolved package manager", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), deps("pnpm"))).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("pnpm"))
		Expect(invocation.Args).To(Equal([]string{"run", "build"}))
		Expect(invocation.Dir).To(Equal("projects/web"))
	})

	// npm is the one package manager that needs the script arguments separated
	// from its own, so the delimiter has to come out of the resolved commands.
	It("separates the script arguments with the delimiter of the package manager", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Build.ScriptArgs = "--verbose"

		Expect(run(runner, p, deps("npm"))).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"run", "build", "--", "--verbose"}))
	})

	It("templates the script and its arguments against the selected environment", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Build.Script = "build:{{ .Environment }}"
		p.Build.ScriptArgs = "--url {{ index .EnvVars \"API_URL\" }}"

		Expect(run(runner, p, deps("pnpm"))).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).
			To(Equal([]string{"run", "build:production", "--url https://api.example.com"}))
	})

	It("hands the environment variables to the build process", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), deps("pnpm"))).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Env).To(ContainElement("API_URL=https://api.example.com"))
	})
})
