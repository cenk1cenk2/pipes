package setup

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Terraform version", func() {
	probe := func(stdout string) string {
		GinkgoHelper()

		*C = Ctx{Env: map[string]string{}}

		runner := fixtures.Runner(tests.TestingCommandResponse{
			Name:   "terraform",
			Args:   []string{"version"},
			Stdout: stdout,
		})

		Expect(fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-terraform",
			CommandName: "setup",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(version(tl).Job())
						})
				},
			},
		}).Run()).To(Succeed())

		return C.Version
	}

	It("takes the version out of the banner", func() {
		Expect(probe("Terraform v1.9.8\non linux_amd64\n")).To(Equal("v1.9.8"))
	})

	// the banner is only ever logged, and terraform has already proven it runs by
	// answering at all, so an unrecognised one is reported rather than fatal.
	It("reports the whole output when the banner does not match", func() {
		Expect(probe("something else entirely\n")).To(Equal("something else entirely"))
	})

	It("trims the surrounding whitespace", func() {
		Expect(probe("  Terraform v1.9.8  \n\n")).To(Equal("v1.9.8"))
	})
})
