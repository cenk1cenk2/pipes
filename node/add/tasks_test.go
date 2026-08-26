package pipe

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

var _ = Describe("Node add", func() {
	// The flags are not registered with the spec command, since the pipe is seeded
	// directly and a package level flag only reads its environment on first parse.
	run := func(runner *tests.TestingCommandRunner, p Pipe, deps Deps) error {
		GinkgoHelper()

		*P = p

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-node",
			CommandName: "add",
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *ucli.Command) *plumber.TaskList {
					tl := &plumber.TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *plumber.TaskList) plumber.Job {
							return plumber.JobSequence(AddNodeModules(tl, deps).Job())
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
				EnvVars:     map[string]string{"REGISTRY": "registry.example.com"},
			},
		}
	}

	pipe := func() Pipe {
		return Pipe{NodeAdd: NodeAdd{Packages: []string{"typescript", "eslint"}, Cwd: "projects/web"}}
	}

	It("adds the packages through the resolved package manager", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), deps("pnpm"))).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("pnpm"))
		Expect(invocation.Args).To(Equal([]string{"add", "typescript", "eslint"}))
		Expect(invocation.Dir).To(Equal("projects/web"))
	})

	// Where the global switch goes differs between the package managers, so it has
	// to come out of the resolved commands rather than be spelled here.
	It("installs globally with the switch of the package manager", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.NodeAdd.Global = true

		Expect(run(runner, p, deps("yarn"))).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"global", "add", "typescript", "eslint"}))
	})

	It("templates the script arguments against the selected environment", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.NodeAdd.ScriptArgs = "--registry={{ index .EnvVars \"REGISTRY\" }}"

		Expect(run(runner, p, deps("pnpm"))).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).
			To(Equal([]string{"add", "--registry=registry.example.com", "typescript", "eslint"}))
	})
})
