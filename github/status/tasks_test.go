package status

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/github/client"
	"gitlab.kilic.dev/devops/pipes/github/test/keys"
	mockclient "gitlab.kilic.dev/devops/pipes/github/test/mocks/client"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("GitHub commit status", func() {
	var (
		adapter *mockclient.MockApplicationClientAdapter
		pem     string
	)

	BeforeEach(func() {
		adapter = mockclient.NewMockApplicationClientAdapter(GinkgoT())
		_, pem = keys.Generate()

		*P = Pipe{
			App: client.Config{ApiUrl: "https://api.github.com"},
			Status: Status{
				Token:   "ghs_given",
				Project: "burningforge/listr2",
				State:   "success",
				Sha:     "abc123",
				Context: "Gitlab CI",
			},
		}
		*C = Ctx{Client: adapter, Repository: "listr2"}
	})

	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	run := func(tasks ...func(*TaskList) *Task) error {
		GinkgoHelper()

		return fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
			AppName:     "pipe-github",
			CommandName: "status",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							jobs := make([]Job, 0, len(tasks))

							for _, task := range tasks {
								jobs = append(jobs, task(tl).Job())
							}

							return JobSequence(jobs...)
						})
				},
			},
		}).Run()
	}

	Describe("authenticate", func() {
		It("takes the token it was given without the app credentials", func() {
			Expect(run(authenticate)).To(Succeed())

			Expect(C.Token).To(Equal("ghs_given"))
		})

		It("mints a token of its own, narrowed to the statuses of the repository", func() {
			P.App = client.Config{Id: "123", InstallationId: "456", PrivateKey: pem, ApiUrl: "https://api.github.com"}

			adapter.EXPECT().
				CreateInstallationToken(mock.Anything, mock.AnythingOfType("string"), "456", client.TokenRequest{
					Repositories: []string{"listr2"},
					Permissions:  map[string]string{"statuses": "write"},
				}).
				Return("ghs_minted", nil)

			Expect(run(authenticate)).To(Succeed())

			Expect(C.Token).To(Equal("ghs_minted"))
		})
	})

	Describe("post", func() {
		It("posts the status for the commit with the resolved token", func() {
			C.Token = "ghs_token"
			P.Status.TargetUrl = "https://gitlab.kilic.dev/pipelines/1"
			P.Status.Description = "passed"

			adapter.EXPECT().
				CreateCommitStatus(mock.Anything, "ghs_token", "burningforge/listr2", "abc123", client.CommitStatus{
					State:       "success",
					TargetUrl:   "https://gitlab.kilic.dev/pipelines/1",
					Description: "passed",
					Context:     "Gitlab CI",
				}).
				Return(nil)

			Expect(run(post)).To(Succeed())
		})
	})

	Describe("New", func() {
		type received struct {
			path          string
			authorization string
			body          string
		}

		var (
			requests []received
			plumber  *Plumber
			output   *bytes.Buffer
		)

		BeforeEach(func() {
			requests = []received{}
			output = &bytes.Buffer{}

			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					body, err := io.ReadAll(r.Body)
					Expect(err).NotTo(HaveOccurred())

					requests = append(requests, received{
						path:          r.URL.Path,
						authorization: r.Header.Get("Authorization"),
						body:          string(body),
					})

					w.WriteHeader(http.StatusCreated)
					_, _ = w.Write([]byte(`{"token":"ghs_minted"}`))
				}),
			)
			DeferCleanup(server.Close)

			P.App.ApiUrl = server.URL
		})

		command := func(flags []cli.Flag, environment map[string]string) *tests.TaskListCliFixture {
			GinkgoHelper()

			return fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
				AppName:     "pipe-github",
				CommandName: "status",
				Flags:       flags,
				Environment: environment,
				WithoutEnvironment: []string{
					"GITHUB_APP_ID",
					"GITHUB_APP_INSTALLATION_ID",
					"GITHUB_APP_PRIVATE_KEY",
					"GITHUB_API_URL",
					"GITHUB_STATUS_TOKEN",
					"GH_TOKEN",
					"GITHUB_STATUS_PROJECT",
					"GITHUB_STATUS_STATE",
					"GITHUB_STATUS_REPORT",
					"GITHUB_STATUS_SHA",
					"CI_COMMIT_SHA",
					"GITHUB_STATUS_TARGET_URL",
					"CI_PIPELINE_URL",
					"GITHUB_STATUS_CONTEXT",
					"GITHUB_STATUS_DESCRIPTION",
				},
				TaskLists: []tests.TaskListFactory{
					func(p *Plumber, _ *cli.Command) *TaskList {
						plumber = p.SetLoggerOutput(output)

						return New(p)
					},
				},
			})
		}

		// the pipelines still export what the curl based reference read, blanking
		// the new names where a job inherits them.
		It("posts the status from the environment the previous reference used", func() {
			fixture := command(Flags, map[string]string{
				"GITHUB_API_URL":           P.App.ApiUrl,
				"GITHUB_STATUS_TOKEN":      "",
				"GH_TOKEN":                 "ghs_legacy",
				"GITHUB_STATUS_PROJECT":    "burningforge/listr2",
				"GITHUB_STATUS_STATE":      "",
				"GITHUB_STATUS_REPORT":     "pending",
				"GITHUB_STATUS_SHA":        "",
				"CI_COMMIT_SHA":            "def456",
				"GITHUB_STATUS_TARGET_URL": "",
				"CI_PIPELINE_URL":          "https://gitlab.kilic.dev/pipelines/2",
			})

			Expect(fixture.Run()).To(Succeed())

			Expect(requests).To(Equal([]received{{
				path:          "/repos/burningforge/listr2/statuses/def456",
				authorization: "Bearer ghs_legacy",
				body:          `{"state":"pending","target_url":"https://gitlab.kilic.dev/pipelines/2","context":"Gitlab CI"}`,
			}}))

			plumber.Log.Info("posted with ghs_legacy")
			Expect(output.String()).To(ContainSubstring("posted with"))
			Expect(output.String()).NotTo(ContainSubstring("ghs_legacy"))
		})

		It("mints its own token in place and posts with it", func() {
			P.App = client.Config{Id: "123", InstallationId: "456", PrivateKey: pem, ApiUrl: P.App.ApiUrl}
			P.Status.Token = ""

			Expect(command(nil, nil).Run()).To(Succeed())

			Expect(requests).To(HaveLen(2))
			Expect(requests[0].path).To(Equal("/app/installations/456/access_tokens"))
			Expect(requests[0].body).To(Equal(`{"repositories":["listr2"],"permissions":{"statuses":"write"}}`))
			Expect(requests[1].path).To(Equal("/repos/burningforge/listr2/statuses/abc123"))
			Expect(requests[1].authorization).To(Equal("Bearer ghs_minted"))

			plumber.Log.Info("posted with ghs_minted")
			Expect(output.String()).NotTo(ContainSubstring("ghs_minted"))
		})

		It("rejects a state GitHub does not know", func() {
			P.Status.State = "done"

			Expect(command(nil, nil).Run()).To(MatchError("Validation failed."))
			Expect(requests).To(BeEmpty())
		})

		It("refuses to run without a token or the app credentials", func() {
			P.Status.Token = ""

			Expect(command(nil, nil).Run()).
				To(MatchError("Either a GitHub token or the GitHub App credentials are required to post a status."))
			Expect(requests).To(BeEmpty())
		})

		DescribeTable(
			"refuses a project that is not owner/repository",
			func(project string) {
				P.Status.Project = project

				Expect(command(nil, nil).Run()).
					To(MatchError("Project has to be in the form of owner/repository: " + project))
				Expect(requests).To(BeEmpty())
			},
			Entry("no owner", "listr2"),
			Entry("empty owner", "/listr2"),
			Entry("empty repository", "burningforge/"),
			Entry("nested", "burningforge/listr2/extra"),
		)
	})
})
