package install

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/node"
	"gitlab.kilic.dev/devops/pipes/node/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Node install", func() {
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
			CommandName: "install",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(install(tl).Job(), prune(tl).Job())
						})
				},
			},
		}).Run()
	}

	pipe := func() Pipe {
		return Pipe{Install: Install{Cwd: "projects/web", Args: "--foo"}}
	}

	It("installs against the lockfile of the resolved package manager", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Install.UseLockFile = true

		Expect(run(runner, p, "pnpm")).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("pnpm"))
		Expect(invocation.Args).To(Equal([]string{"i", "--frozen-lockfile", "--foo"}))
		Expect(invocation.Dir).To(Equal("projects/web"))
	})

	It("installs without the lockfile when the pipeline asked for it", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), "npm")).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"i", "--unsafe-perm", "--foo"}))
	})

	// the cache directory is named after the package manager, and the flag it is
	// passed under differs between them.
	It("points the package manager at a cache named after it", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Install.Cache = true

		Expect(run(runner, p, "yarn")).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"install", "--foo", "--prefer-offline", "--cache-folder", ".yarn"}))
	})

	It("prunes the pnpm store it installed into", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Install.Cache = true

		Expect(run(runner, p, "pnpm")).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"pnpm", "pnpm"}))

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"store", "prune", "--store-dir", ".pnpm"}))
		Expect(invocation.Dir).To(Equal("projects/web"))
	})

	It("leaves the pnpm store alone when nothing is cached", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe(), "pnpm")).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"pnpm"}))
	})

	// only pnpm prunes what nothing references, the other caches are left as they are.
	DescribeTable("leaves the cache of the other package managers alone",
		func(packageManager string) {
			runner := fixtures.Runner()

			p := pipe()
			p.Install.Cache = true

			Expect(run(runner, p, packageManager)).To(Succeed())

			Expect(runner.InvocationNames()).To(Equal([]string{packageManager}))
		},
		Entry("npm", "npm"),
		Entry("yarn", "yarn"),
	)
})
