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
		if opts.Command == nil {
			opts.Command = []string{"helm", "registry", "login"}
		}

		return registry.NewLogin(opts)
	}

	parse := func(value string) ([]registry.LoginEntry, error) {
		GinkgoHelper()

		entries := []registry.LoginEntry{}
		flags := registry.NewLoginFlags(registry.LoginFlagsOptions{
			Category:    "Container Registry",
			Alias:       "PIPE_CONTAINER_LOGIN",
			Destination: &entries,
		})

		return entries, flags[0].(*cli.StringFlag).Validator(value)
	}

	run := func(runner *tests.TestingCommandRunner, l *registry.Login) error {
		GinkgoHelper()

		return fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe",
			CommandName: "login",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						Set(func(tl *TaskList) Job {
							return JobSequence(l.Job(tl))
						})
				},
			},
		}).Run()
	}

	Describe("parsing", func() {
		It("reads a list of entries", func() {
			Expect(parse("- uri: ghcr.io\n  username: user\n  password: secret\n")).
				To(Equal([]registry.LoginEntry{{Uri: "ghcr.io", Username: "user", Password: "secret"}}))
		})

		// the template default composes the list as JSON, which only works because
		// YAML is a superset of it.
		It("reads a JSON list of entries", func() {
			Expect(parse(`[{ "uri": "docker.io", "username": "user", "password": "secret" }]`)).
				To(Equal([]registry.LoginEntry{{Uri: "docker.io", Username: "user", Password: "secret"}}))
		})

		DescribeTable(
			"rejects anything but a list",
			func(value string) {
				_, err := parse(value)

				Expect(err).To(MatchError(ContainSubstring("container.login")))
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
			opts.Entries = []registry.LoginEntry{
				{Uri: "docker.io", Username: "user", Password: "secret"},
				{Uri: "ghcr.io", Username: "other", Password: "token"},
			}

			return newLogin(opts)
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

						tl := &TaskList{}

						return tl.New(p).
							Set(func(tl *TaskList) Job {
								return JobSequence(l.Job(tl))
							})
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
			l := newLogin(registry.LoginOptions{Entries: []registry.LoginEntry{{Username: "user", Password: "secret"}}})
			runner := fixtures.Runner()

			Expect(run(runner, l)).NotTo(Succeed())
			Expect(runner.Invocations()).To(BeEmpty())
		})

		// a bad credential fails at the registry, so the entry goes to the login
		// command as written.
		DescribeTable(
			"hands an entry to the login command as written",
			func(entry registry.LoginEntry) {
				l := newLogin(registry.LoginOptions{Entries: []registry.LoginEntry{entry}})
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

	// buildah publishes under an image the login has to stand in for.
	DescribeTable(
		"covers an image on its host or a path above it",
		func(name string, uri string, expected bool) {
			image, err := registry.ParseImage(name)
			Expect(err).NotTo(HaveOccurred())

			l := newLogin(registry.LoginOptions{Entries: []registry.LoginEntry{{Uri: uri, Username: "user", Password: "secret"}}})

			Expect(l.Covers(image)).To(Equal(expected))
		},
		Entry("the same host", "registry.example.com/image", "registry.example.com", true),
		Entry("the host of a path", "registry.example.com/group/project", "registry.example.com", true),
		Entry("a parent path", "registry.example.com/group/project", "registry.example.com/group", true),
		Entry("another host", "registry.example.com/image", "docker.io", false),
		Entry("a host sharing the prefix", "registry.example.com.evil/image", "registry.example.com", false),
		Entry("a deeper path than the image", "registry.example.com/group", "registry.example.com/group/project", false),
	)
})

var _ = Describe("Container image name", func() {
	DescribeTable(
		"splits a fully qualified image name at its registry host",
		func(name string, expected registry.Image) {
			image, err := registry.ParseImage(name)
			Expect(err).NotTo(HaveOccurred())
			Expect(*image).To(Equal(expected))
			Expect(image.String()).To(Equal(name))
		},
		Entry("a registry host", "docker.io/cenk1cenk2/pipe-buildah", registry.Image{Host: "docker.io", Path: "cenk1cenk2/pipe-buildah"}),
		Entry("a registry with a port", "registry.example.com:5000/image", registry.Image{Host: "registry.example.com:5000", Path: "image"}),
		Entry("a nested project path", "registry.example.com/group/project/image:v1", registry.Image{Host: "registry.example.com", Path: "group/project/image:v1"}),
		Entry("localhost", "localhost/image", registry.Image{Host: "localhost", Path: "image"}),
	)

	// the login no longer names the images, so a name without its registry would
	// leave the container tooling to guess one from its short-name configuration.
	DescribeTable(
		"fails on an image name without the registry host",
		func(name string) {
			_, err := registry.ParseImage(name)

			Expect(err).To(MatchError(ContainSubstring("Image name is not fully qualified")))
		},
		Entry("a bare name", "image"),
		Entry("a docker hub shorthand", "cenk1cenk2/pipe-buildah"),
		Entry("a project path", "group/project/image"),
	)
})
