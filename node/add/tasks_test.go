package add

import (
	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/node"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Node add", func() {
	seed := func(packageManager string) {
		*setup.NodeCtx = node.Ctx{PackageManager: node.PackageManager{
			Exe:      packageManager,
			Commands: node.PackageManagers[packageManager],
		}}
	}

	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe, packageManager string) error {
		GinkgoHelper()

		*P = pipe
		seed(packageManager)

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-node",
			CommandName: "add",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(add(tl).Job())
						})
				},
			},
		}).Run()
	}

	pipe := func() Pipe {
		return Pipe{Add: Add{Packages: []string{"typescript", "eslint"}, Cwd: "projects/web"}}
	}

	It("adds the packages through the resolved package manager", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), "pnpm")).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("pnpm"))
		Expect(invocation.Args).To(Equal([]string{"add", "typescript", "eslint"}))
		Expect(invocation.Dir).To(Equal("projects/web"))
	})

	// where the global switch goes differs between the package managers, so it has
	// to come out of the resolved commands rather than be spelled here.
	It("installs globally with the switch of the package manager", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Add.Global = true

		Expect(run(runner, p, "yarn")).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"global", "add", "typescript", "eslint"}))
	})

	It("appends the script arguments before the packages", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Add.ScriptArgs = "--registry=registry.example.com"

		Expect(run(runner, p, "pnpm")).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).
			To(Equal([]string{"add", "--registry=registry.example.com", "typescript", "eslint"}))
	})
})
