package run

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/python/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Python run", func() {
	It("runs the arguments it was given inside the installed environment", func() {
		runner := fixtures.Runner()

		*P = Pipe{}
		*setup.C = setup.Ctx{Cwd: "projects/api", Env: map[string]string{}}

		Expect(fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-python",
			CommandName: "run",
			Args:        []string{"pipe-python", "run", "--", "python", "-m", "app"},
			Arguments:   Arguments,
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					return New(p)
				},
			},
		}).Run()).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("uv"))
		Expect(invocation.Args).To(Equal([]string{"run", "--no-sync", "python", "-m", "app"}))
		Expect(invocation.Dir).To(Equal("projects/api"))
	})
})
