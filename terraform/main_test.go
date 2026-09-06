package main

import (
	"strings"
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
	RunSpecs(t, "Terraform Pipe Suite")
}

var _ = conformance.Verify(conformance.Pipe{
	Name:        name,
	Description: description,
	New: func(p *plumber.Plumber) *ucli.Command {
		return newCommand(p, "test", options{})
	},
	UncategorizedFlags: []string{
		"terraform.apply.args",
		"terraform.apply.output",
		"terraform.install.args",
		"terraform.install.reconfigure",
		"terraform.install.use-lockfile",
		"terraform.lint.format-check.args",
		"terraform.lint.format-check.enable",
		"terraform.lint.validate.args",
		"terraform.lint.validate.enable",
		"terraform.login.registry.credentials",
		"terraform.plan.args",
		"terraform.plan.output",
		"terraform.plan.pipeline-source",
		"terraform.plan.preview-for-merge-requests",
		"terraform.plan.retry-delay",
		"terraform.plan.retry-tries",
		"terraform.plan.summary.output",
	},
	LegacyEnvAliases: map[string][]string{
		"terraform.apply.args":                            {"TF_APPLY_ARGS", "TERRAFORM_APPLY_ARGS"},
		"terraform.apply.output":                          {"TF_PLAN_CACHE", "TF_APPLY_OUTPUT", "TF_PLAN_OUTPUT", "TERRAFORM_APPLY_OUTPUT"},
		"terraform.ci.api-url":                            {"TF_VAR_CI_API_V4_URL", "CI_API_V4_URL", "TERRAFORM_CI_API_URL"},
		"terraform.ci.project-id":                         {"TF_VAR_CI_PROJECT_ID", "CI_PROJECT_ID", "TERRAFORM_CI_PROJECT_ID"},
		"terraform.cwd":                                   {"TF_ROOT", "TERRAFORM_CWD"},
		"terraform.install.args":                          {"TF_INSTALL_ARGS", "TERRAFORM_INSTALL_ARGS"},
		"terraform.install.reconfigure":                   {"TF_INSTALL_RECONFIGURE", "TERRAFORM_INSTALL_RECONFIGURE"},
		"terraform.install.use-lockfile":                  {"TF_INSTALL_USE_LOCKFILE", "TERRAFORM_INSTALL_USE_LOCKFILE"},
		"terraform.lint.format-check.args":                {"TF_LINT_FMT_CHECK_ARGS", "TERRAFORM_LINT_FORMAT_CHECK_ARGS"},
		"terraform.lint.format-check.enable":              {"TF_LINT_FMT_CHECK_ENABLE", "TERRAFORM_LINT_FORMAT_CHECK_ENABLE"},
		"terraform.lint.validate.args":                    {"TF_LINT_VALIDATE_ARGS", "TERRAFORM_LINT_VALIDATE_ARGS"},
		"terraform.lint.validate.enable":                  {"TF_LINT_VALIDATE_ENABLE", "TERRAFORM_LINT_VALIDATE_ENABLE"},
		"terraform.log-level":                             {"TF_LOG_LEVEL", "TF_LOG", "TERRAFORM_LOG_LEVEL"},
		"terraform.login.registry.credentials":            {"TF_REGISTRY_CREDENTIALS", "TERRAFORM_LOGIN_REGISTRY_CREDENTIALS"},
		"terraform.plan.args":                             {"TF_PLAN_ARGS", "TERRAFORM_PLAN_ARGS"},
		"terraform.plan.output":                           {"TF_PLAN_CACHE", "TF_APPLY_OUTPUT", "TF_PLAN_OUTPUT", "TERRAFORM_PLAN_OUTPUT"},
		"terraform.plan.pipeline-source":                  {"CI_PIPELINE_SOURCE", "TERRAFORM_PLAN_PIPELINE_SOURCE"},
		"terraform.plan.preview-for-merge-requests":       {"TF_PLAN_PREVIEW_FOR_MRS", "TERRAFORM_PLAN_PREVIEW_FOR_MERGE_REQUESTS"},
		"terraform.plan.retry-delay":                      {"TF_PLAN_RETRY_DELAY", "TERRAFORM_PLAN_RETRY_DELAY"},
		"terraform.plan.retry-tries":                      {"TF_PLAN_RETRY_TRIES", "TERRAFORM_PLAN_RETRY_TRIES"},
		"terraform.plan.summary.output":                   {"TERRAFORM_SUMMARY_OUTPUT", "TERRAFORM_PLAN_SUMMARY_OUTPUT"},
		"terraform.publish.module.cwd":                    {"TF_MODULE_CWD", "TF_ROOT", "TERRAFORM_PUBLISH_MODULE_CWD"},
		"terraform.publish.module.name":                   {"TF_MODULE_NAME", "CI_PROJECT_NAME", "TERRAFORM_PUBLISH_MODULE_NAME"},
		"terraform.publish.module.system":                 {"TF_MODULE_SYSTEM", "TERRAFORM_PUBLISH_MODULE_SYSTEM"},
		"terraform.publish.registry.gitlab.api-url":       {"CI_API_V4_URL", "TERRAFORM_PUBLISH_REGISTRY_GITLAB_API_URL"},
		"terraform.publish.registry.gitlab.project-id":    {"CI_PROJECT_ID", "TERRAFORM_PUBLISH_REGISTRY_GITLAB_PROJECT_ID"},
		"terraform.publish.registry.gitlab.token":         {"CI_JOB_TOKEN", "TERRAFORM_PUBLISH_REGISTRY_GITLAB_TOKEN"},
		"terraform.publish.registry.name":                 {"TF_MODULE_REGISTRY", "TERRAFORM_PUBLISH_REGISTRY_NAME"},
		"terraform.state.gitlab-http.http-address":        {"TF_HTTP_ADDRESS", "TF_ADDRESS", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_ADDRESS"},
		"terraform.state.gitlab-http.http-lock-address":   {"TF_HTTP_LOCK_ADDRESS", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_LOCK_ADDRESS"},
		"terraform.state.gitlab-http.http-lock-method":    {"TF_HTTP_LOCK_METHOD", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_LOCK_METHOD"},
		"terraform.state.gitlab-http.http-password":       {"TF_HTTP_PASSWORD", "TF_PASSWORD", "CI_JOB_TOKEN", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_PASSWORD"},
		"terraform.state.gitlab-http.http-retry-wait-min": {"TF_HTTP_RETRY_WAIT_MIN", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_RETRY_WAIT_MIN"},
		"terraform.state.gitlab-http.http-unlock-address": {"TF_HTTP_UNLOCK_ADDRESS", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_UNLOCK_ADDRESS"},
		"terraform.state.gitlab-http.http-unlock-method":  {"TF_HTTP_UNLOCK_METHOD", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_UNLOCK_METHOD"},
		"terraform.state.gitlab-http.http-username":       {"TF_HTTP_USERNAME", "TF_USERNAME", "TERRAFORM_STATE_GITLAB_HTTP_HTTP_USERNAME"},
		"terraform.state.name":                            {"TF_STATE_NAME", "TERRAFORM_STATE_NAME"},
		"terraform.state.strict":                          {"TF_STATE_STRICT", "TERRAFORM_STATE_STRICT"},
		"terraform.state.type":                            {"TF_STATE_TYPE", "TERRAFORM_STATE_TYPE"},
	},
})

// Everything a spec needs is passed as an argument rather than through the
// environment, since the flags are package level and urfave only reads an env
// source on the first parse of a flag instance: driving values in through the
// environment would make the specs depend on the order Ginkgo runs them in.
func run(runner *tests.TestingCommandRunner, opts options, args ...string) error {
	GinkgoHelper()

	// $CI_PIPELINE_SOURCE has no argument to pass it through, so the pipeline this
	// suite runs in is the one the plan reads: on a merge request it plans as a
	// preview and then deletes a plan file the stubbed command never wrote.
	tests.WithoutEnvironment("CI_PIPELINE_SOURCE")

	fixture := fixtures.NewPlumber(func(p *plumber.Plumber) *ucli.Command {
		return newCommand(p, "test", opts)
	})
	fixture.Plumber.SetRuntime(plumber.Runtime{CommandRunner: runner.Runner()})

	return fixture.RunCli(append([]string{name}, args...)...)
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

		Expect(run(runner, options{}, "install")).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"terraform", "terraform"}))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("init -input=false")))
	})

	It("checks the formatting and validates on lint", func() {
		runner := fixtures.Runner()

		Expect(run(runner, options{}, "lint")).To(Succeed())

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
			options{Notes: func(c gitlab.MergeRequestReportConfig) (gitlab.Notes, error) {
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

		Expect(run(runner, options{}, "apply")).To(Succeed())

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
			options{Registry: func(apiUrl, projectId, token string) gitlab.ModuleRegistry {
				dialed = []string{apiUrl, projectId, token}

				return nil
			}},
			"publish",
			"--terraform.publish.module.name", "vpc",
			"--terraform.publish.registry.gitlab.api-url", "https://gitlab.example.test/api/v4",
			"--terraform.publish.registry.gitlab.project-id", "3",
			"--terraform.publish.registry.gitlab.token", "job-token",
		)).To(Succeed())

		Expect(dialed).To(Equal([]string{"https://gitlab.example.test/api/v4", "3", "job-token"}))
		Expect(runner.Invocations()).To(BeEmpty())
	})
})
