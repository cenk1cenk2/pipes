package login

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/registry"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Container registry login verification", func() {
	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	run := func(runner *tests.TestingCommandRunner, entries []registry.LoginEntry, image string) error {
		GinkgoHelper()

		P.Container = entries

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-buildah",
			CommandName: "build",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					return New(p, image)
				},
			},
		}).Run()
	}

	// the verify is the only login call without credentials.
	verifies := func(runner *tests.TestingCommandRunner) [][]string {
		args := [][]string{}
		for _, invocation := range runner.Invocations() {
			if len(invocation.Args) == 2 {
				args = append(args, invocation.Args)
			}
		}

		return args
	}

	// a pipeline relying on an ambient login still has to prove it exists before
	// the build spends time on an image it cannot push.
	It("verifies the ambient login on the registry of the image", func() {
		runner := fixtures.Runner()

		Expect(run(runner, nil, "registry.example.com/group/image")).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal("buildah"))
		Expect(invocation.Args).To(Equal([]string{"login", "registry.example.com"}))
	})

	It("renders a templated image before reading its registry", func() {
		runner := fixtures.Runner()

		Expect(run(runner, nil, `{{ "registry.example.com/group/image" }}`)).To(Succeed())

		Expect(verifies(runner)).To(Equal([][]string{{"login", "registry.example.com"}}))
	})

	It("logs in with the container login before the verify", func() {
		runner := fixtures.Runner()

		Expect(run(runner, []registry.LoginEntry{{Uri: "docker.io", Username: "user", Password: "secret"}}, "registry.example.com/group/image")).To(Succeed())

		invocations := runner.Invocations()
		Expect(invocations).To(HaveLen(2))
		Expect(invocations[0].Args).To(Equal([]string{"login", "docker.io", "--username", "user", "--password-stdin"}))
		Expect(invocations[1].Args).To(Equal([]string{"login", "registry.example.com"}))
	})

	// a container login entry on the registry the image publishes under already
	// proves the login the verify would look for.
	DescribeTable(
		"skips the verify when the container login covers the image",
		func(uri string, expected [][]string) {
			runner := fixtures.Runner()

			Expect(run(runner, []registry.LoginEntry{{Uri: uri, Username: "user", Password: "secret"}}, "registry.example.com/group/image")).To(Succeed())
			Expect(verifies(runner)).To(Equal(expected))
		},
		Entry("the registry host", "registry.example.com", [][]string{}),
		Entry("the project path", "registry.example.com/group", [][]string{}),
		Entry("another registry", "docker.io", [][]string{{"login", "registry.example.com"}}),
		Entry("a sibling path", "registry.example.com/other", [][]string{{"login", "registry.example.com"}}),
	)

	It("verifies nothing without an image", func() {
		runner := fixtures.Runner()

		Expect(run(runner, nil, "")).To(Succeed())
		Expect(runner.Invocations()).To(BeEmpty())
	})

	It("fails on an image name without the registry host", func() {
		Expect(run(fixtures.Runner(), nil, "group/image")).
			To(MatchError(ContainSubstring("Image name is not fully qualified")))
	})
})
