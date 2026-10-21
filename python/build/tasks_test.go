package build

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/python/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Python build", func() {
	It("builds the distributions into the output directory", func() {
		runner := fixtures.Runner()

		*P = Pipe{Output: "dist", Args: "--wheel"}
		*setup.C = setup.Ctx{Cwd: "projects/api", Env: map[string]string{}}

		Expect(fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-python",
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
		}).Run()).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("uv"))
		Expect(invocation.Args).To(Equal([]string{"build", "--out-dir", "dist", "--wheel"}))
		Expect(invocation.Dir).To(Equal("projects/api"))
	})
})
