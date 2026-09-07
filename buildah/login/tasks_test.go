package login

import (
	"bytes"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Container registry login", func() {
	// The flags are not registered with the spec command, since the pipe is seeded
	// directly and a package level flag only reads its environment on first parse.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe) error {
		GinkgoHelper()

		*P = pipe

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-buildah",
			CommandName: "login",
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *cli.Command) *plumber.TaskList {
					return New(p)
				},
			},
		}).Run()
	}

	disabled := func(pipe Pipe, task func(*plumber.TaskList) *plumber.Task) bool {
		GinkgoHelper()

		*P = pipe

		tl := &plumber.TaskList{}
		tl.New(plumber.NewPlumber(func(_ *plumber.Plumber) *cli.Command {
			return &cli.Command{Name: "test"}
		}))

		return task(tl).IsDisabled()
	}

	credentials := func() Pipe {
		return Pipe{Uri: "registry.example.com", Username: "user", Password: "secret"}
	}

	It("logs in to the registry the pipeline named", func() {
		runner := fixtures.Runner()

		Expect(run(runner, credentials())).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("buildah"))
		Expect(invocation.Args).To(Equal([]string{
			"login", "registry.example.com", "--username", "user", "--password-stdin",
		}))
	})

	// An argument list shows up in a process listing and in the command trace log,
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

	// Half a credential pair is not something to guess at, and a pipeline with
	// neither is relying on an ambient login the pipe must not clobber.
	DescribeTable(
		"logs in only on a complete credential pair",
		func(username, password string, expected bool) {
			pipe := credentials()
			pipe.Username = username
			pipe.Password = password

			Expect(disabled(pipe, ContainerRegistryLogin)).To(Equal(expected))
		},
		Entry("both set", "user", "secret", false),
		Entry("no password", "user", "", true),
		Entry("no username", "", "secret", true),
		Entry("neither", "", "", true),
	)

	// The verify task is the inverse of the login, so exactly one of the two runs
	// under the parent whichever way the pipeline is configured.
	DescribeTable(
		"verifies the ambient login only when there is nothing to log in with",
		func(username, password string, expected bool) {
			pipe := credentials()
			pipe.Username = username
			pipe.Password = password

			Expect(disabled(pipe, ContainerRegistryLoginVerify)).To(Equal(expected))
		},
		Entry("both set", "user", "secret", true),
		Entry("no password", "user", "", false),
		Entry("no username", "", "secret", false),
		Entry("neither", "", "", false),
	)

	// A pipeline relying on an ambient login still has to prove it exists before
	// the build spends time on an image it cannot push.
	It("verifies the ambient login when the pipeline carries no credentials", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{Uri: "registry.example.com"})).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("buildah"))
		Expect(invocation.Args).To(Equal([]string{"login", "registry.example.com"}))
	})

	// The password is the one value that must never reach a log line, and masking
	// it here is what keeps every command the pipe runs after it from having to
	// remember to.
	It("keeps the password out of the log", func() {
		var (
			log    *logrus.Logger
			output = &bytes.Buffer{}
		)

		*P = credentials()

		Expect(fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
			AppName:     "pipe-buildah",
			CommandName: "login",
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *cli.Command) *plumber.TaskList {
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
