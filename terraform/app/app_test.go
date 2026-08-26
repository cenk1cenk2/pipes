package app_test

import (
	"strings"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	ucli "github.com/urfave/cli/v3"
	clientgitlab "gitlab.com/gitlab-org/api/client-go/v2"

	"gitlab.kilic.dev/devops/pipes/internal/gitlab"
	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	mockgitlab "gitlab.kilic.dev/devops/pipes/internal/test/mocks/gitlab"
	"gitlab.kilic.dev/devops/pipes/terraform/app"
)

// Everything a spec needs is passed as an argument rather than through the
// environment, since the flags are package level and urfave only reads an env
// source on the first parse of a flag instance: driving values in through the
// environment would make the specs depend on the order Ginkgo runs them in.
func run(runner *tests.TestingCommandRunner, opts app.Options, args ...string) error {
	GinkgoHelper()

	fixture := fixtures.NewPlumber(func(p *plumber.Plumber) *ucli.Command {
		return app.New(p, "test", opts)
	})
	fixture.Plumber.SetRuntime(plumber.Runtime{CommandRunner: runner.Runner()})

	return fixture.RunCli(append([]string{"pipe-terraform"}, args...)...)
}

// Answers the terraform invocation that starts with the given subcommand. Once a
// runner carries any response every invocation has to match one, so a spec that
// stubs output seeds one of these per command it expects to run.
func responds(subcommand, stdout string) tests.TestingCommandResponse {
	// The stream recorder only keeps whole lines, so output that does not end in a
	// newline never reaches the pipe reading it back.
	if stdout != "" && !strings.HasSuffix(stdout, "\n") {
		stdout += "\n"
	}

	return tests.TestingCommandResponse{
		Match: func(invocation plumber.CommandInvocation) bool {
			return len(invocation.Args) > 0 && invocation.Args[0] == subcommand
		},
		Stdout: stdout,
	}
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
	It("probes the tool and initializes the project on install", func() {
		runner := fixtures.Runner()

		Expect(run(runner, app.Options{}, "install")).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"terraform", "terraform"}))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("init -input=false")))
	})

	It("checks the formatting and validates on lint", func() {
		runner := fixtures.Runner()

		Expect(run(runner, app.Options{}, "lint")).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"terraform", "terraform", "terraform"}))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("fmt -check -diff -recursive")))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("validate")))
	})

	It("plans the project and reports through the notes the options supplied", func() {
		cwd := tests.TempDir()
		// The summary and the merge request report each read the plan back, so the
		// show is answered twice: a response is consumed by the invocation it
		// matches.
		runner := fixtures.Runner(
			responds("version", "Terraform v1.9.8"),
			responds("plan", ""),
			responds("show", `{"format_version":"1.2"}`),
			responds("show", `{"format_version":"1.2"}`),
		)

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
			app.Options{Notes: func(c gitlab.MergeRequestReportConfig) (gitlab.Notes, error) {
				config = c

				return notes, nil
			}},
			"plan",
			"--terraform.cwd", cwd,
			"--gitlab-mr-report.enabled",
			"--gitlab-mr-report.merge-request-iid", "452",
			"--gitlab-mr-report.identifier", "plan",
			"--gitlab-mr-report.token", "glpat-token",
			"--gitlab-mr-report.api-url", "https://gitlab.example.test/api/v4",
			"--gitlab-mr-report.project-id", "3",
		)).To(Succeed())

		Expect(formatted(runner)).To(ContainElement(ContainSubstring("plan -input=false -out=plan")))
		Expect(config.ProjectId).To(Equal("3"))
		Expect(config.Token).To(Equal("glpat-token"))
		Expect(body).To(ContainSubstring("gitlab-pipes:mr-report:plan"))
	})

	It("applies the plan on apply", func() {
		runner := fixtures.Runner()

		Expect(run(runner, app.Options{}, "apply")).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"terraform", "terraform"}))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("apply")))
	})

	// The tags file an earlier job would have written is absent, so the publish
	// finds nothing to upload. What the spec is after is the registry the options
	// supplied being the one the pipe dials.
	It("dials the module registry the options supplied on publish", func() {
		runner := fixtures.Runner()

		var dialed []string

		Expect(run(
			runner,
			app.Options{Registry: func(apiUrl, projectId, token string) gitlab.ModuleRegistry {
				dialed = []string{apiUrl, projectId, token}

				return nil
			}},
			"publish",
			"--terraform-module.name", "vpc",
			"--terraform-module.registry.gitlab.api-url", "https://gitlab.example.test/api/v4",
			"--terraform-module.registry.gitlab.project-id", "3",
			"--terraform-module.registry.gitlab.token", "job-token",
		)).To(Succeed())

		Expect(dialed).To(Equal([]string{"https://gitlab.example.test/api/v4", "3", "job-token"}))
		Expect(runner.Invocations()).To(BeEmpty())
	})
})
