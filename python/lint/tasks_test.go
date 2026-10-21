package lint

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/python/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Python lint", func() {
	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe) error {
		GinkgoHelper()

		*P = pipe
		*setup.C = setup.Ctx{Cwd: "projects/api", Env: map[string]string{"UV_CACHE_DIR": "/builds/.uv"}}

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-python",
			CommandName: "lint",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(lint(tl).Job())
						})
				},
			},
		}).Run()
	}

	It("runs every command in order inside the installed environment", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Commands: []string{"ruff format --check .", "ty  check src"}})).To(Succeed())

		invocations := runner.Invocations()
		Expect(invocations).To(HaveLen(2))
		Expect(invocations[0].Args).To(Equal([]string{"run", "--no-sync", "ruff", "format", "--check", "."}))
		Expect(invocations[1].Args).To(Equal([]string{"run", "--no-sync", "ty", "check", "src"}))
		Expect(invocations[1].Dir).To(Equal("projects/api"))
	})
})
