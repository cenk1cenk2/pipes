package release

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
)

var _ = Describe("Semantic release", func() {
	// The flags are not registered with the spec command, since the pipe is seeded
	// directly and a package level flag only reads its environment on first parse.
	run := func(runner *tests.TestingCommandRunner, pipe Pipe, debug bool) error {
		GinkgoHelper()

		*P = pipe
		*C = Ctx{}

		fixture := fixtures.Cli(runner, tests.TaskListCli{
			AppName:     "pipe-semantic-release",
			CommandName: "release",
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *ucli.Command) *plumber.TaskList {
					tl := &plumber.TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *plumber.TaskList) plumber.Job {
							return plumber.JobSequence(RunSemanticRelease(tl).Job())
						})
				},
			},
		})

		// The pipe reads its own debug off the log level, and the fixture leaves it
		// on trace so that a spec can see what the task list did.
		if !debug {
			fixture.Plumber.Log.SetLevel(logrus.InfoLevel)
		}

		return fixture.Run()
	}

	It("runs the release through the single package binary", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{}, false)).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal(SEMANTIC_RELEASE_EXE))
		Expect(invocation.Args).To(BeEmpty())
	})

	It("runs the release through the workspace binary when the pipe is a workspace", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{SemanticRelease: SemanticRelease{Workspace: true}}, false)).
			To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Name).To(Equal(MULTI_SEMANTIC_RELEASE_EXE))
	})

	// semantic-release detects the pipeline it runs in and refuses a dry run on a
	// branch it was not released from, so a dry run has to hide the pipeline.
	Describe("dry run", func() {
		var invocation plumber.CommandInvocation

		BeforeEach(func() {
			runner := fixtures.Runner()

			Expect(run(runner, Pipe{
				SemanticRelease: SemanticRelease{DryRun: true},
				CI:              CI{CommitReference: "feature/branch"},
			}, false)).To(Succeed())

			var ok bool
			invocation, ok = runner.LastInvocation()
			Expect(ok).To(BeTrue())
		})

		It("branches the release off the current reference", func() {
			Expect(invocation.Args).
				To(Equal([]string{"--dry-run", "--no-ci", "--branches", "feature/branch"}))
		})

		It("hides the pipeline from the release library", func() {
			Expect(invocation.Env).To(ContainElements("CI=false", "GITLAB_CI=false"))
		})
	})

	It("hands the debug flag down when the pipe itself runs in debug", func() {
		runner := fixtures.Runner()

		Expect(run(runner, Pipe{}, true)).To(Succeed())

		invocation, ok := runner.LastInvocation()
		Expect(ok).To(BeTrue())
		Expect(invocation.Args).To(Equal([]string{"--debug"}))
	})
})
