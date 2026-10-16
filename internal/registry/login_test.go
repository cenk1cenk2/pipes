package registry_test

import (
	"bytes"
	"fmt"
	"os"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/registry"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Container registry login", func() {
	newLogin := func(opts registry.LoginOptions) *registry.Login {
		opts.Category = "Container Registry"
		opts.Alias = "PIPE_CONTAINER_LOGIN"
		if opts.Command == nil {
			opts.Command = []string{"helm", "registry", "login"}
		}

		return registry.NewLogin(opts)
	}

	parse := func(l *registry.Login, value string) error {
		GinkgoHelper()

		return l.Flags[0].(*cli.StringFlag).Validator(value)
	}

	run := func(runner *tests.TestingCommandRunner, l *registry.Login) error {
		GinkgoHelper()

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe",
			CommandName: "login",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					return l.New(p)
				},
			},
		}).Run()
	}

	Describe("parsing", func() {
		It("reads a list of entries", func() {
			l := newLogin(registry.LoginOptions{})

			Expect(parse(l, "- uri: ghcr.io\n  username: user\n  password: secret\n")).To(Succeed())
			Expect(l.Entries).To(Equal([]registry.LoginEntry{{Uri: "ghcr.io", Username: "user", Password: "secret"}}))
		})

		// the template default composes the list as JSON, which only works because
		// YAML is a superset of it.
		It("reads a JSON list of entries", func() {
			l := newLogin(registry.LoginOptions{})

			Expect(parse(l, `[{ "uri": "docker.io", "username": "user", "password": "secret" }]`)).To(Succeed())
			Expect(l.Entries).To(Equal([]registry.LoginEntry{{Uri: "docker.io", Username: "user", Password: "secret"}}))
		})

		DescribeTable(
			"rejects anything but a list",
			func(value string) {
				Expect(parse(newLogin(registry.LoginOptions{}), value)).To(MatchError(ContainSubstring("container.login")))
			},
			Entry("a single entry", "uri: ghcr.io\nusername: user\npassword: secret\n"),
			Entry("a JSON object", `{ "username": "user", "password": "secret" }`),
			Entry("a scalar", "docker.io"),
		)
	})

	Describe("logging in", func() {
		// the variable is process wide, so every spec starts without it and leaves
		// whatever the environment had before behind.
		BeforeEach(func() {
			previous, ok := os.LookupEnv("HELM_REGISTRY_CONFIG")
			Expect(os.Unsetenv("HELM_REGISTRY_CONFIG")).To(Succeed())

			DeferCleanup(func() {
				if ok {
					Expect(os.Setenv("HELM_REGISTRY_CONFIG", previous)).To(Succeed())

					return
				}

				Expect(os.Unsetenv("HELM_REGISTRY_CONFIG")).To(Succeed())
			})
		})

		credentials := func(opts registry.LoginOptions) *registry.Login {
			l := newLogin(opts)
			l.Entries = []registry.LoginEntry{
				{Uri: "docker.io", Username: "user", Password: "secret"},
				{Uri: "ghcr.io", Username: "other", Password: "token"},
			}

			return l
		}

		It("logs in to every entry with the pipe's command", func() {
			runner := fixtures.Runner()

			Expect(run(runner, credentials(registry.LoginOptions{Command: []string{"buildah", "login"}}))).To(Succeed())

			invocations := runner.Invocations()
			Expect(invocations).To(HaveLen(2))
			Expect(invocations[0].Name).To(Equal("buildah"))
			Expect(invocations[0].Args).To(Equal([]string{"login", "docker.io", "--username", "user", "--password-stdin"}))
			Expect(invocations[1].Args).To(Equal([]string{"login", "ghcr.io", "--username", "other", "--password-stdin"}))
		})

		// an argument list shows up in a process listing and in the command trace log,
		// which is exactly what the password has to stay out of.
		It("hands every password over stdin rather than in the arguments", func() {
			runner := fixtures.Runner()

			Expect(run(runner, credentials(registry.LoginOptions{}))).To(Succeed())

			for i, password := range []string{"secret", "token"} {
				invocation := runner.Invocations()[i]
				Expect(invocation.Args).NotTo(ContainElement(password))

				stdin, err := tests.ReadInvocationStdin(invocation)
				Expect(err).NotTo(HaveOccurred())
				Expect(stdin).To(Equal(password))
			}
		})

		It("keeps every password out of the log", func() {
			var (
				plumber *Plumber
				output  = &bytes.Buffer{}
				l       = credentials(registry.LoginOptions{})
			)

			Expect(fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
				AppName:     "pipe",
				CommandName: "login",
				TaskLists: []tests.TaskListFactory{
					func(p *Plumber, _ *cli.Command) *TaskList {
						plumber = p.SetLoggerOutput(output)

						return l.New(p)
					},
				},
			}).Run()).To(Succeed())

			plumber.Log.Info(fmt.Sprintf("logging in with %s and %s", "secret", "token"))
			Expect(output.String()).NotTo(ContainSubstring("secret"))
			Expect(output.String()).NotTo(ContainSubstring("token"))
		})

		It("runs nothing without an entry", func() {
			runner := fixtures.Runner()

			Expect(run(runner, newLogin(registry.LoginOptions{}))).To(Succeed())
			Expect(runner.Invocations()).To(BeEmpty())
		})

		// the registry is spelled out on every entry, there is no default to fall
		// back to.
		It("fails on an entry without a uri", func() {
			l := newLogin(registry.LoginOptions{})
			l.Entries = []registry.LoginEntry{{Username: "user", Password: "secret"}}
			runner := fixtures.Runner()

			Expect(run(runner, l)).NotTo(Succeed())
			Expect(runner.Invocations()).To(BeEmpty())
		})

		// a bad credential fails at the registry, so the entry goes to the login
		// command as written.
		DescribeTable(
			"hands an entry to the login command as written",
			func(entry registry.LoginEntry) {
				l := newLogin(registry.LoginOptions{})
				l.Entries = []registry.LoginEntry{entry}
				runner := fixtures.Runner()

				Expect(run(runner, l)).To(Succeed())

				invocation, ok := runner.LastInvocation()
				Expect(ok).To(BeTrue())
				Expect(invocation.Args).To(Equal([]string{"registry", "login", entry.Uri, "--username", entry.Username, "--password-stdin"}))
			},
			Entry("no password", registry.LoginEntry{Uri: "docker.io", Username: "user"}),
			Entry("an unexpanded variable reference", registry.LoginEntry{Uri: "docker.io", Username: "$DOCKERHUB_USERNAME", Password: "$DOCKERHUB_PASSWORD"}),
		)

		// the chart inflation runs helm with a temporary config home, so the login
		// is only visible to it through a registry config pinned for the process.
		It("pins the helm registry config it was given", func() {
			Expect(run(fixtures.Runner(), credentials(registry.LoginOptions{HelmRegistryConfig: "/tmp/pipe/config.json"}))).To(Succeed())

			Expect(os.Getenv("HELM_REGISTRY_CONFIG")).To(Equal("/tmp/pipe/config.json"))
		})

		It("keeps a helm registry config the pipeline already set", func() {
			Expect(os.Setenv("HELM_REGISTRY_CONFIG", "/custom/config.json")).To(Succeed())

			Expect(run(fixtures.Runner(), credentials(registry.LoginOptions{HelmRegistryConfig: "/tmp/pipe/config.json"}))).To(Succeed())

			Expect(os.Getenv("HELM_REGISTRY_CONFIG")).To(Equal("/custom/config.json"))
		})

		It("pins nothing without an entry to log in with", func() {
			Expect(run(fixtures.Runner(), newLogin(registry.LoginOptions{HelmRegistryConfig: "/tmp/pipe/config.json"}))).To(Succeed())

			_, ok := os.LookupEnv("HELM_REGISTRY_CONFIG")
			Expect(ok).To(BeFalse())
		})
	})

	// buildah publishes under a target the login has to stand in for.
	DescribeTable(
		"covers a target on its host or a path above it",
		func(target string, uri string, expected bool) {
			l := newLogin(registry.LoginOptions{})
			l.Entries = []registry.LoginEntry{{Uri: uri, Username: "user", Password: "secret"}}

			Expect(l.Covers(target)).To(Equal(expected))
		},
		Entry("the same host", "registry.example.com", "registry.example.com", true),
		Entry("the host of a path", "registry.example.com/group/project", "registry.example.com", true),
		Entry("a parent path", "registry.example.com/group/project", "registry.example.com/group", true),
		Entry("another host", "registry.example.com", "docker.io", false),
		Entry("a host sharing the prefix", "registry.example.com.evil", "registry.example.com", false),
		Entry("a deeper path than the target", "registry.example.com", "registry.example.com/group", false),
	)
})
