package cli_test

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
)

// The steps run commands rather than plain tasks because the command runner is
// what records an order the spec can assert on.
func step(binary string, flags ...ucli.Flag) cli.Step {
	return cli.Step{
		Flags: flags,
		New: func(p *plumber.Plumber) *plumber.TaskList {
			tl := &plumber.TaskList{}

			return tl.New(p).
				SetRuntimeDepth(3).
				Set(func(tl *plumber.TaskList) plumber.Job {
					return plumber.JobSequence(
						tl.CreateTask(binary).
							Set(func(t *plumber.Task) error {
								t.CreateCommand(binary).AddSelfToTheTask()

								return nil
							}).
							ShouldRunAfter(func(t *plumber.Task) error {
								return t.RunCommandJobAsJobSequence()
							}).
							Job(),
					)
				})
		},
	}
}

func run(runner *tests.TestingCommandRunner, command func(*plumber.Plumber) *ucli.Command, args ...string) error {
	GinkgoHelper()

	fixture := fixtures.NewPlumber(command)
	fixture.Plumber.SetRuntime(plumber.Runtime{CommandRunner: runner.Runner()})

	return fixture.RunCli(args...)
}

var _ = Describe("Command", func() {
	It("runs the task lists of every step in the order they were given", func() {
		runner := fixtures.Runner()

		Expect(run(runner, func(p *plumber.Plumber) *ucli.Command {
			return cli.App("pipe", "a pipe", "test", cli.Command(p, "run", "runs", step("first"), step("second")))
		}, "pipe", "run")).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"first", "second"}))
	})

	It("collects the flags and the arguments of every step", func() {
		command := cli.Command(
			nil,
			"run",
			"runs",
			cli.Step{Flags: []ucli.Flag{&ucli.StringFlag{Name: "one"}}},
			cli.Step{
				Flags:     []ucli.Flag{&ucli.StringFlag{Name: "two"}},
				Arguments: []ucli.Argument{&ucli.StringArg{Name: "target"}},
			},
		)

		Expect(command.Flags).To(HaveLen(2))
		Expect(command.Flags[0].Names()).To(ContainElement("one"))
		Expect(command.Flags[1].Names()).To(ContainElement("two"))
		Expect(command.Arguments).To(HaveLen(1))
	})

	// A step reads its flag values while it builds its task list, so a step
	// constructed before the command line was parsed would read the defaults.
	It("builds the task lists only once the command runs", func() {
		built := false

		command := cli.Command(nil, "run", "runs", cli.Step{
			New: func(p *plumber.Plumber) *plumber.TaskList {
				built = true

				return nil
			},
		})

		Expect(command.Action).NotTo(BeNil())
		Expect(built).To(BeFalse())
	})
})

var _ = Describe("Root", func() {
	It("runs the steps without a subcommand", func() {
		runner := fixtures.Runner()

		Expect(run(runner, func(p *plumber.Plumber) *ucli.Command {
			return cli.Root(p, "pipe", "a pipe", "test", step("first"), step("second"))
		}, "pipe")).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"first", "second"}))
	})

	It("carries the version and the description of the pipe", func() {
		command := cli.Root(nil, "pipe", "a pipe", "v1.0.0")

		Expect(command.Name).To(Equal("pipe"))
		Expect(command.Version).To(Equal("v1.0.0"))
		Expect(command.Usage).To(Equal("a pipe"))
		Expect(command.Description).To(Equal("a pipe"))
	})
})

var _ = Describe("App", func() {
	It("dispatches to the command that was named", func() {
		runner := fixtures.Runner()

		Expect(run(runner, func(p *plumber.Plumber) *ucli.Command {
			return cli.App(
				"pipe",
				"a pipe",
				"test",
				cli.Command(p, "build", "builds", step("first")),
				cli.Command(p, "lint", "lints", step("second")),
			)
		}, "pipe", "lint")).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"second"}))
	})
})
