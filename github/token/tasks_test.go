package token

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

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

var _ = Describe("GitHub App token", func() {
	var (
		adapter *mockclient.MockApplicationClientAdapter
		pem     string
		file    string
	)

	BeforeEach(func() {
		adapter = mockclient.NewMockApplicationClientAdapter(GinkgoT())
		_, pem = keys.Generate()
		file = filepath.Join(GinkgoT().TempDir(), "github.env")

		*P = Pipe{
			App: client.Config{
				Id:             "123",
				InstallationId: "456",
				PrivateKey:     pem,
				ApiUrl:         "https://api.github.com",
			},
			Token: Token{Variable: "GH_TOKEN", File: file},
		}
		*C = Ctx{Client: adapter}
	})

	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	run := func(tasks ...func(*TaskList) *Task) error {
		GinkgoHelper()

		return fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
			AppName:     "pipe-github",
			CommandName: "token",
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

	Describe("mint", func() {
		It("requests a token for the installation, narrowed down as configured", func() {
			P.Token.Repositories = []string{"listr2"}
			P.Token.Permissions = map[string]string{"contents": "write"}

			adapter.EXPECT().
				CreateInstallationToken(mock.Anything, mock.AnythingOfType("string"), "456", client.TokenRequest{
					Repositories: []string{"listr2"},
					Permissions:  map[string]string{"contents": "write"},
				}).
				Return("ghs_token", nil)

			Expect(run(mint)).To(Succeed())
			Expect(C.Token).To(Equal("ghs_token"))
		})
	})

	Describe("write", func() {
		BeforeEach(func() {
			C.Token = "ghs_token"
		})

		// GitLab keeps the quotes of a dotenv value as part of it, so a quoted token
		// would reach the consumer jobs with the quotes still on.
		It("writes the token unquoted", func() {
			Expect(run(write)).To(Succeed())

			Expect(os.ReadFile(file)).To(BeEquivalentTo("GH_TOKEN=ghs_token\n"))
		})

		It("keeps the file readable by the job alone", func() {
			Expect(run(write)).To(Succeed())

			info, err := os.Stat(file)
			Expect(err).NotTo(HaveOccurred())
			Expect(info.Mode().Perm()).To(Equal(os.FileMode(0o600)))
		})

		It("writes the token under the variable it was given", func() {
			P.Token.Variable = "RELEASE_TOKEN"

			Expect(run(write)).To(Succeed())

			Expect(os.ReadFile(file)).To(BeEquivalentTo("RELEASE_TOKEN=ghs_token\n"))
		})

		// several apps share one dotenv report, each writing its own token into it.
		It("keeps the other variables of an existing file", func() {
			Expect(os.WriteFile(file, []byte("RELEASE_TOKEN=ghs_release\nAPP_ID=12\n"), 0o600)).To(Succeed())

			Expect(run(write)).To(Succeed())

			Expect(os.ReadFile(file)).To(BeEquivalentTo("APP_ID=12\nGH_TOKEN=ghs_token\nRELEASE_TOKEN=ghs_release\n"))
		})

		It("replaces the token of a rerun instead of adding another", func() {
			Expect(os.WriteFile(file, []byte("GH_TOKEN=ghs_stale\n"), 0o600)).To(Succeed())

			Expect(run(write)).To(Succeed())

			Expect(os.ReadFile(file)).To(BeEquivalentTo("GH_TOKEN=ghs_token\n"))
		})

		// semantic-release pushes with GIT_CREDENTIALS verbatim, and an installation token is only accepted as the password of x-access-token.
		It("writes the token as a git credential under the variable it was given", func() {
			P.Token.GitCredentialsVariable = "RELEASE_GIT_CREDENTIALS"

			Expect(run(write)).To(Succeed())

			Expect(os.ReadFile(file)).To(BeEquivalentTo("GH_TOKEN=ghs_token\nRELEASE_GIT_CREDENTIALS=x-access-token:ghs_token\n"))
		})

		// a second token job for another app shares the dotenv report, and only one of them may hand over the git credential.
		It("writes no git credential when its variable is empty", func() {
			P.Token.GitCredentialsVariable = ""

			Expect(run(write)).To(Succeed())

			Expect(os.ReadFile(file)).To(BeEquivalentTo("GH_TOKEN=ghs_token\n"))
		})
	})

	Describe("New", func() {
		var (
			server  *httptest.Server
			plumber *Plumber
			output  *bytes.Buffer
		)

		BeforeEach(func() {
			server = httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.WriteHeader(http.StatusCreated)
					_, _ = w.Write([]byte(`{"token":"ghs_minted"}`))
				}),
			)
			DeferCleanup(server.Close)

			output = &bytes.Buffer{}
		})

		command := func(flags []cli.Flag, environment map[string]string) *tests.TaskListCliFixture {
			GinkgoHelper()

			return fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
				AppName:     "pipe-github",
				CommandName: "token",
				Flags:       flags,
				Environment: environment,
				TaskLists: []tests.TaskListFactory{
					func(p *Plumber, _ *cli.Command) *TaskList {
						plumber = p.SetLoggerOutput(output)

						return New(p)
					},
				},
			})
		}

		It("mints the token into the dotenv file from the environment a job exports", func() {
			key := filepath.Join(GinkgoT().TempDir(), "key.pem")
			Expect(os.WriteFile(key, []byte(pem), 0o600)).To(Succeed())

			fixture := command(Flags, map[string]string{
				"GITHUB_APP_ID":              "123",
				"GITHUB_APP_INSTALLATION_ID": "456",
				"GITHUB_APP_PRIVATE_KEY":     key,
				"GITHUB_API_URL":             server.URL,
				"GITHUB_TOKEN_VARIABLE":      "RELEASE_TOKEN",
				"GITHUB_TOKEN_FILE":          file,
			})

			Expect(fixture.Run()).To(Succeed())
			Expect(fixture.ExitCodes()).To(BeEmpty())

			Expect(os.ReadFile(file)).To(BeEquivalentTo("GIT_CREDENTIALS=x-access-token:ghs_minted\nRELEASE_TOKEN=ghs_minted\n"))

			plumber.Log.Info("minted ghs_minted, pushing with x-access-token:ghs_minted")
			Expect(output.String()).To(ContainSubstring("minted"))
			Expect(output.String()).NotTo(ContainSubstring("ghs_minted"))
		})

		It("refuses a variable name GitLab would not export", func() {
			P.App.ApiUrl = server.URL
			P.Token.Variable = "GH-TOKEN"

			fixture := command(nil, nil)

			Expect(fixture.Run()).To(MatchError("Variable name can only contain letters, digits and underscores: GH-TOKEN"))
			Expect(file).NotTo(BeAnExistingFile())
		})

		It("refuses to run without the app credentials", func() {
			P.App = client.Config{ApiUrl: server.URL}

			fixture := command(nil, nil)

			Expect(fixture.Run()).To(MatchError("GitHub App id, installation id and private key are required to mint a token."))
		})
	})
})

var _ = Describe("marshal", func() {
	It("writes the variables sorted, one per line", func() {
		Expect(marshal(map[string]string{"B": "2", "A": "1"})).To(Equal("A=1\nB=2\n"))
	})

	DescribeTable(
		"rejects a value a dotenv file can not carry unquoted",
		func(value string) {
			_, err := marshal(map[string]string{"GH_TOKEN": value})

			Expect(err).To(MatchError("Value of variable can not carry newlines, quotes or comments in a dotenv file: GH_TOKEN"))
		},
		Entry("newline", "ghs\ntoken"),
		Entry("carriage return", "ghs\rtoken"),
		Entry("double quote", `ghs"token`),
		Entry("single quote", "ghs'token"),
		Entry("backtick", "ghs`token"),
		Entry("comment", "ghs#token"),
	)

	DescribeTable(
		"takes only variable names GitLab exports",
		func(name string, valid bool) {
			_, err := marshal(map[string]string{name: "ghs_token"})

			if valid {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(MatchError(ContainSubstring("Variable name can only contain letters, digits and underscores")))
			}
		},
		Entry("upper case", "GH_TOKEN", true),
		Entry("mixed case with digits", "Release_Token_2", true),
		Entry("dash", "GH-TOKEN", false),
		Entry("space", "GH TOKEN", false),
		Entry("equals", "GH=TOKEN", false),
		Entry("empty", "", false),
	)
})
