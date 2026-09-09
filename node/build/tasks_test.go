package build

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/node"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Node build", func() {
	seed := func(packageManager string) {
		*setup.NodeCtx = node.Ctx{PackageManager: node.PackageManager{
			Exe:      packageManager,
			Commands: node.PackageManagers[packageManager],
		}}
		*setup.EnvironmentCtx = environment.Ctx{
			Environment: "production",
			EnvVars:     map[string]string{"API_URL": "https://api.example.com"},
		}
	}

	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe, packageManager string) error {
		GinkgoHelper()

		*P = pipe
		seed(packageManager)

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-node",
			CommandName: "build",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(build(tl).Job())
						})
				},
			},
		}).Run()
	}

	pipe := func() Pipe {
		return Pipe{Build: Build{Script: "build", Cwd: "projects/web"}}
	}

	It("runs the build script through the resolved package manager", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), "pnpm")).To(Succeed())

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

		Expect(run(runner, p, "npm")).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"run", "build", "--", "--verbose"}))
	})

	It("templates the script and its arguments against the selected environment", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Build.Script = "build:{{ .Environment }}"
		p.Build.ScriptArgs = "--url {{ index .EnvVars \"API_URL\" }}"

		Expect(run(runner, p, "pnpm")).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).
			To(Equal([]string{"run", "build:production", "--url https://api.example.com"}))
	})

	It("hands the environment variables to the build process", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), "pnpm")).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Env).To(ContainElement("API_URL=https://api.example.com"))
	})
})
