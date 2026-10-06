package test

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/helm/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Helm unittest", func() {
	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe) error {
		GinkgoHelper()

		*setup.C = setup.Ctx{Cwd: "charts/app", Env: map[string]string{}}
		*P = pipe

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-helm",
			CommandName: "test",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					return New(p)
				},
			},
		}).Run()
	}

	pipe := func() Pipe {
		return Pipe{Output: Output{Type: "JUnit"}}
	}

	It("runs the default test suites of the chart in the directory the setup step resolved", func() {
		runner := fixtures.Runner()

		Expect(run(runner, pipe())).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("helm"))
		Expect(invocation.Args).To(Equal([]string{"unittest", "."}))
		Expect(invocation.Dir).To(Equal("charts/app"))
	})

	It("passes every test suite glob, values file and strict parsing before the chart", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Files = []string{"tests/*_test.yaml", "tests/**/*_test.yaml"}
		p.Values = []string{"tests/values/*.yaml"}
		p.Strict = true

		Expect(run(runner, p)).To(Succeed())

		invocation, _ := runner.LastInvocation()
		Expect(invocation.Args).To(Equal([]string{
			"unittest",
			"--file", "tests/*_test.yaml",
			"--file", "tests/**/*_test.yaml",
			"--values", "tests/values/*.yaml",
			"--strict",
			".",
		}))
	})

	// the type names the format of a report, so without a file there is nothing it applies to.
	DescribeTable(
		"writes a report only when a file is given",
		func(file string, args []string) {
			runner := fixtures.Runner()

			p := pipe()
			p.Output.File = file

			Expect(run(runner, p)).To(Succeed())

			invocation, _ := runner.LastInvocation()
			Expect(invocation.Args).To(Equal(args))
		},
		Entry("without a file", "", []string{"unittest", "."}),
		Entry(
			"with a file",
			"helm-unittest.xml",
			[]string{"unittest", "--output-file", "helm-unittest.xml", "--output-type", "JUnit", "."},
		),
	)

	It("rejects a report format helm-unittest does not write", func() {
		runner := fixtures.Runner()

		p := pipe()
		p.Output.Type = "TAP"

		Expect(run(runner, p)).NotTo(Succeed())
		Expect(runner.Invocations()).To(BeEmpty())
	})
})
