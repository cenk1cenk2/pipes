package login

import (
	"bytes"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Helm registry login", func() {
	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe) error {
		GinkgoHelper()

		*P = pipe

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-helm",
			CommandName: "login",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					return New(p)
				},
			},
		}).Run()
	}

	credentials := func() Pipe {
		return Pipe{Uri: "registry.example.com", Username: "user", Password: "secret"}
	}

	It("logs in to the chart repository the pipeline named", func() {
		runner := fixtures.Runner()

		Expect(run(runner, credentials())).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("helm"))
		Expect(invocation.Args).To(Equal([]string{
			"registry", "login", "registry.example.com", "--username", "user", "--password-stdin",
		}))
	})

	// an argument list shows up in a process listing and in the command trace log,
	// which is exactly what the password has to stay out of.
	It("hands the password over stdin rather than in the arguments", func() {
		runner := fixtures.Runner()

		Expect(run(runner, credentials())).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).NotTo(ContainElement("secret"))

		stdin, err := tests.ReadInvocationStdin(invocation)
		Expect(err).NotTo(HaveOccurred())
		Expect(stdin).To(Equal("secret"))
	})

	// half a credential pair is not something to guess at, and a pipeline with
	// neither is relying on an ambient login the pipe must not clobber. Helm has
	// no verify counterpart, so nothing runs at all.
	DescribeTable(
		"logs in only on a complete credential pair",
		func(username, password string, invocations int) {
			pipe := credentials()
			pipe.Username = username
			pipe.Password = password

			runner := fixtures.Runner()

			Expect(run(runner, pipe)).To(Succeed())
			Expect(runner.Invocations()).To(HaveLen(invocations))
		},
		Entry("both set", "user", "secret", 1),
		Entry("no password", "user", "", 0),
		Entry("no username", "", "secret", 0),
		Entry("neither", "", "", 0),
	)

	// the password is the one value that must never reach a log line, and masking
	// it here is what keeps every command the pipe runs after it from having to
	// remember to.
	It("keeps the password out of the log", func() {
		var (
			log    *logrus.Logger
			output = &bytes.Buffer{}
		)

		*P = credentials()

		Expect(fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
			AppName:     "pipe-helm",
			CommandName: "login",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					log = p.Log
					log.SetOutput(output)

					return New(p)
				},
			},
		}).Run()).To(Succeed())

		log.Infof("logging in with %s", P.Password)
		Expect(output.String()).NotTo(ContainSubstring("secret"))
	})
})
