package plugins

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Pulumi plugins", func() {
	Describe("Outdated", func() {
		It("keeps the newest version of a plugin", func() {
			outdated, err := Outdated([]byte(`[
				{ "kind": "resource", "name": "vault", "version": "1.2.0" },
				{ "kind": "resource", "name": "vault", "version": "1.10.0" }
			]`))

			Expect(err).NotTo(HaveOccurred())
			Expect(outdated).To(Equal([]Plugin{{Kind: "resource", Name: "vault", Version: "1.2.0"}}))
		})

		It("keeps the only version of a plugin", func() {
			outdated, err := Outdated([]byte(`[{ "kind": "resource", "name": "vault", "version": "1.2.0" }]`))

			Expect(err).NotTo(HaveOccurred())
			Expect(outdated).To(BeEmpty())
		})

		It("treats every kind and name on its own", func() {
			outdated, err := Outdated([]byte(`[
				{ "kind": "resource", "name": "vault", "version": "1.0.0" },
				{ "kind": "resource", "name": "aws", "version": "6.0.0" },
				{ "kind": "language", "name": "vault", "version": "9.0.0" },
				{ "kind": "resource", "name": "vault", "version": "1.1.0" },
				{ "kind": "resource", "name": "aws", "version": "5.0.0" }
			]`))

			Expect(err).NotTo(HaveOccurred())
			Expect(outdated).To(ConsistOf(
				Plugin{Kind: "resource", Name: "vault", Version: "1.0.0"},
				Plugin{Kind: "resource", Name: "aws", Version: "5.0.0"},
			))
		})

		It("keeps the versions it can not compare", func() {
			outdated, err := Outdated([]byte(`[
				{ "kind": "resource", "name": "vault", "version": "1.0.0" },
				{ "kind": "resource", "name": "vault", "version": "2.0.0" },
				{ "kind": "resource", "name": "vault", "version": "not-a-version" },
				{ "kind": "resource", "name": "vault" }
			]`))

			Expect(err).NotTo(HaveOccurred())
			Expect(outdated).To(Equal([]Plugin{{Kind: "resource", Name: "vault", Version: "1.0.0"}}))
		})

		It("fails on a list that is no JSON", func() {
			_, err := Outdated([]byte(`not json`))

			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Prune", func() {
		run := func(runner *tests.TestingCommandRunner, env map[string]string) error {
			GinkgoHelper()

			*setup.C = setup.Ctx{Cwd: "projects/api", Env: env}

			return fixtures.Cli(runner, tests.TaskListCli{
				AppName:     "pipe-pulumi",
				CommandName: "prune",
				TaskLists: []tests.TaskListFactory{
					func(p *Plumber, _ *cli.Command) *TaskList {
						tl := &TaskList{}

						return tl.New(p).
							SetRuntimeDepth(3).
							Set(func(tl *TaskList) Job {
								return JobSequence(Prune(tl).Job())
							})
					},
				},
			}).Run()
		}

		list := tests.TestingCommandResponse{
			Name: "pulumi",
			Args: []string{"plugin", "ls", "--output", "json"},
			Stdout: `[
				{ "kind": "resource", "name": "vault", "version": "1.2.0" },
				{ "kind": "resource", "name": "vault", "version": "1.10.0" },
				{ "kind": "resource", "name": "aws", "version": "6.0.0" }
			]` + "\n",
		}

		It("does not run without a Pulumi home", func() {
			GinkgoT().Setenv("PULUMI_HOME", "")
			runner := fixtures.Runner(list)

			Expect(run(runner, map[string]string{})).To(Succeed())

			Expect(runner.Invocations()).To(BeEmpty())
		})

		It("removes every version but the newest", func() {
			GinkgoT().Setenv("PULUMI_HOME", "")
			runner := fixtures.Runner(list, tests.TestingCommandResponse{Name: "pulumi"})

			Expect(run(runner, map[string]string{"PULUMI_HOME": "/cache"})).To(Succeed())

			invocations := runner.Invocations()
			Expect(invocations).To(HaveLen(2))
			Expect(invocations[0].Args).To(Equal([]string{"plugin", "ls", "--output", "json"}))
			Expect(invocations[1].Args).To(Equal([]string{"plugin", "rm", "resource", "vault", "1.2.0", "--yes"}))
			Expect(invocations[1].Dir).To(Equal("projects/api"))
		})

		It("reads the Pulumi home of the environment too", func() {
			GinkgoT().Setenv("PULUMI_HOME", "/cache")
			runner := fixtures.Runner(list, tests.TestingCommandResponse{Name: "pulumi"})

			Expect(run(runner, map[string]string{})).To(Succeed())

			Expect(runner.Invocations()).To(HaveLen(2))
		})

		It("does not fail the run when the plugins are not readable", func() {
			GinkgoT().Setenv("PULUMI_HOME", "/cache")
			runner := fixtures.Runner(tests.TestingCommandResponse{Name: "pulumi", Stdout: "not json\n"})

			Expect(run(runner, map[string]string{})).To(Succeed())

			Expect(runner.Invocations()).To(HaveLen(1))
		})
	})
})
