package install

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/python/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Python install", func() {
	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe, env map[string]string) error {
		GinkgoHelper()

		*P = pipe
		*setup.C = setup.Ctx{Cwd: "projects/api", Env: env}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-python",
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

	It("installs against a lockfile that has to be up to date", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Locked: true, Args: "--all-groups"}, map[string]string{})).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("uv"))
		Expect(invocation.Args).To(Equal([]string{"sync", "--locked", "--all-groups"}))
		Expect(invocation.Dir).To(Equal("projects/api"))
	})

	It("installs the lockfile as it is when the pipeline asked for it", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{}, map[string]string{})).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"sync", "--frozen"}))
	})

	It("prunes the cache it installed from", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Locked: true}, map[string]string{"UV_CACHE_DIR": "/builds/.uv"})).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"uv", "uv"}))

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"cache", "prune", "--ci"}))
		Expect(invocation.Dir).To(Equal("projects/api"))
	})

	It("leaves the cache of the environment alone", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Locked: true}, map[string]string{})).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"uv"}))
	})
})
