package install

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

var _ = Describe("Node install", func() {
	// The flags are not registered with the spec command, since the pipe is seeded
	// directly and a package level flag only reads its environment on first parse.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe, deps Deps) error {
		GinkgoHelper()

		*P = pipe

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-node",
			CommandName: "install",
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *ucli.Command) *plumber.TaskList {
					tl := &plumber.TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *plumber.TaskList) plumber.Job {
							return plumber.JobSequence(InstallNodeDependencies(tl, deps).Job())
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
				EnvVars: map[string]string{"NODE_AUTH_TOKEN": "npm-token"},
			},
		}
	}

	pipe := func() Pipe {
		return Pipe{Install: Install{Cwd: "projects/web", Args: "--foo"}}
	}

	It("installs against the lockfile of the resolved package manager", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Install.UseLockFile = true

		Expect(run(runner, p, deps("pnpm"))).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("pnpm"))
		Expect(invocation.Args).To(Equal([]string{"i", "--frozen-lockfile", "--foo"}))
		Expect(invocation.Dir).To(Equal("projects/web"))
	})

	It("installs without the lockfile when the pipeline asked for it", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), deps("npm"))).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"i", "--unsafe-perm", "--foo"}))
	})

	// The cache directory is named after the package manager, and the flag it is
	// passed under differs between them.
	It("points the package manager at a cache named after it", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Install.Cache = true

		Expect(run(runner, p, deps("yarn"))).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"install", "--foo", "--prefer-offline", "--cache-folder", ".yarn"}))
	})

	It("hands the environment variables to the package manager", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), deps("pnpm"))).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Env).To(ContainElement("NODE_AUTH_TOKEN=npm-token"))
	})
})
