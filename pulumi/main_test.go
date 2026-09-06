package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	ucli "github.com/urfave/cli/v3"
	clientgitlab "gitlab.com/gitlab-org/api/client-go/v2"

	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/test/conformance"
	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	mockgitlab "gitlab.kilic.dev/devops/pipes/internal/test/mocks/gitlab"
)

func TestPipe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pulumi Pipe Suite")
}

var _ = conformance.Verify(conformance.Pipe{
	Name:        name,
	Description: description,
	New: func(p *plumber.Plumber) *ucli.Command {
		return newCommand(p, "test", options{})
	},
	UncategorizedFlags: []string{
		"pulumi.preview.plan",
		"pulumi.preview.summary.output",
		"pulumi.stack",
		"pulumi.up.plan",
	},
	// Both commands have always read $PULUMI_PLAN, so a job that runs preview and
	// up in turn hands them the same file. The canonical names are per command,
	// which is the way out of that without breaking the jobs that rely on it.
	LegacyEnvAliases: map[string][]string{
		"pulumi.preview.plan":           {"PULUMI_PLAN", "PULUMI_PREVIEW_PLAN"},
		"pulumi.preview.summary.output": {"PULUMI_SUMMARY_OUTPUT", "PULUMI_PREVIEW_SUMMARY_OUTPUT"},
		"pulumi.up.plan":                {"PULUMI_PLAN", "PULUMI_UP_PLAN"},
	},
})

// Everything a spec needs is passed as an argument rather than through the
// environment, since the flags are package level and urfave only reads an env
// source on the first parse of a flag instance: driving values in through the
// environment would make the specs depend on the order Ginkgo runs them in.
func run(runner *tests.TestingCommandRunner, opts options, args ...string) error {
	GinkgoHelper()

	fixture := fixtures.NewPlumber(func(p *plumber.Plumber) *ucli.Command {
		return newCommand(p, "test", opts)
	})
	fixture.Plumber.SetRuntime(plumber.Runtime{CommandRunner: runner.Runner()})

	return fixture.RunCli(append([]string{name}, args...)...)
}

func formatted(runner *tests.TestingCommandRunner) []string {
	GinkgoHelper()

	invocations := runner.Invocations()
	commands := make([]string, len(invocations))

	for i, invocation := range invocations {
		commands[i] = invocation.Formatted
	}

	return commands
}

var _ = Describe("New", func() {
	It("selects the stack and previews it, reporting through the notes the options supplied", func() {
		cwd := tests.TempDir()
		// The preview reads the plan back out of the file pulumi wrote, which the
		// stubbed command never produces.
		Expect(os.WriteFile(
			filepath.Join(cwd, "plan.json"),
			[]byte(`{"manifest":{"version":"3.187.0"},"resourcePlans":{}}`),
			0o644,
		)).To(Succeed())

		runner := fixtures.Runner()

		notes := mockgitlab.NewMockNotes(GinkgoT())
		notes.EXPECT().
			ListMergeRequestNotes(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, &clientgitlab.Response{}, nil)

		var body string
		notes.EXPECT().
			CreateMergeRequestNote(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			RunAndReturn(func(
				_ any,
				_ int64,
				opt *clientgitlab.CreateMergeRequestNoteOptions,
				_ ...clientgitlab.RequestOptionFunc,
			) (*clientgitlab.Note, *clientgitlab.Response, error) {
				body = *opt.Body

				return &clientgitlab.Note{ID: 1}, &clientgitlab.Response{}, nil
			})

		var config gitlab.MergeRequestReportConfig

		Expect(run(
			runner,
			options{Notes: func(c gitlab.MergeRequestReportConfig) (gitlab.Notes, error) {
				config = c

				return notes, nil
			}},
			"preview",
			"--pulumi.cwd", cwd,
			"--pulumi.stack", "production",
			"--gitlab-mr-report.enabled",
			"--gitlab-mr-report.merge-request-iid", "452",
			"--gitlab-mr-report.identifier", "preview",
			"--gitlab-mr-report.token", "glpat-token",
			"--gitlab-mr-report.api-url", "https://gitlab.example.test/api/v4",
			"--gitlab-mr-report.project-id", "3",
		)).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"pulumi", "pulumi", "pulumi"}))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("stack select production")))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("preview --non-interactive --diff --save-plan plan.json")))
		Expect(config.ProjectId).To(Equal("3"))
		Expect(body).To(ContainSubstring("gitlab-pipes:mr-report:preview"))
	})

	It("selects the stack and applies the plan on up", func() {
		runner := fixtures.Runner()

		Expect(run(runner, options{}, "up", "--pulumi.stack", "production")).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"pulumi", "pulumi", "pulumi"}))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("stack select production")))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("up --diff --yes -f --plan plan.json")))
	})
})
