package update

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/hub"
	mockhub "gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/test/mocks/hub"
)

var _ = Describe("Docker Hub readme", func() {
	var client *mockhub.MockClient

	BeforeEach(func() {
		client = mockhub.NewMockClient(GinkgoT())

		*P = Pipe{DockerHub: DockerHub{
			Username: "user",
			Password: "password",
			Address:  "https://hub.docker.com/v2/repositories",
		}}
		*C = Ctx{Readme: map[string]ParsedReadme{}, Hub: client}
	})

	// The flags are not registered with the spec command, since the pipe is seeded
	// directly and a package level flag only reads its environment on first parse.
	//
	// A failing task terminates the process through the plumber, so only what the
	// pipe does on its way through is asserted here and the errors it answers with
	// are asserted on VerifyReadme instead.
	run := func(tasks ...func(*plumber.TaskList) *plumber.Task) error {
		GinkgoHelper()

		return fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
			AppName:     "pipe-update-docker-hub-readme",
			CommandName: "update",
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *ucli.Command) *plumber.TaskList {
					tl := &plumber.TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *plumber.TaskList) plumber.Job {
							jobs := make([]plumber.Job, 0, len(tasks))

							for _, task := range tasks {
								jobs = append(jobs, task(tl).Job())
							}

							return plumber.JobSequence(jobs...)
						})
				},
			},
		}).Run()
	}

	// Writes a readme file and registers it as the discovered target, since the
	// update reads the file off disk rather than out of the pipe.
	target := func(repository, content, description string) {
		GinkgoHelper()

		file := filepath.Join(GinkgoT().TempDir(), "README.md")
		Expect(os.WriteFile(file, []byte(content), 0600)).To(Succeed())

		C.Readme[repository] = ParsedReadme{File: file, Description: description}
	}

	Describe("login", func() {
		It("keeps the token the credentials bought", func() {
			client.EXPECT().Login(mock.Anything, "user", "password").Return("jwt-token", nil)

			Expect(run(LoginToDockerHubRegistry)).To(Succeed())
			Expect(C.Token).To(Equal("jwt-token"))
		})
	})

	Describe("discover", func() {
		It("takes the single repository the pipe was given", func() {
			P.Readme = Readme{Repository: "kilic/pipe", File: "README.md", Description: "a pipe"}

			Expect(run(DiscoverJobs)).To(Succeed())
			Expect(C.Readme).To(Equal(map[string]ParsedReadme{
				"kilic/pipe": {File: "README.md", Description: "a pipe"},
			}))
		})

		It("takes every repository of the matrix", func() {
			P.Readme = Readme{Matrix: []ReadmeMatrixJson{
				{Repository: "kilic/one", File: "one.md"},
				{Repository: "kilic/two", File: "two.md", Description: "the second"},
			}}

			Expect(run(DiscoverJobs)).To(Succeed())
			Expect(C.Readme).To(Equal(map[string]ParsedReadme{
				"kilic/one": {File: "one.md"},
				"kilic/two": {File: "two.md", Description: "the second"},
			}))
		})

		// The matrix is the way to update more than one repository in a job, so a
		// pipeline that sets both should not have to drop the single target.
		It("takes the single repository alongside the matrix", func() {
			P.Readme = Readme{
				Repository: "kilic/pipe",
				File:       "README.md",
				Matrix:     []ReadmeMatrixJson{{Repository: "kilic/one", File: "one.md"}},
			}

			Expect(run(DiscoverJobs)).To(Succeed())
			Expect(C.Readme).To(HaveLen(2))
		})
	})

	Describe("update", func() {
		BeforeEach(func() {
			C.Token = "jwt-token"
		})

		It("pushes the file contents as the full description", func() {
			target("kilic/pipe", "# Pipe", "a pipe")

			client.EXPECT().
				UpdateReadme(mock.Anything, "jwt-token", "kilic/pipe", hub.Readme{
					Description: "a pipe",
					Full:        "# Pipe",
				}).
				Return(hub.Result{
					StatusCode:      http.StatusOK,
					Description:     "a pipe",
					FullDescription: "# Pipe",
				}, nil)

			Expect(run(UpdateDockerReadme)).To(Succeed())
		})

		It("updates every discovered repository", func() {
			target("kilic/one", "# One", "")
			target("kilic/two", "# Two", "")

			client.EXPECT().
				UpdateReadme(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				RunAndReturn(func(_ context.Context, _, _ string, readme hub.Readme) (hub.Result, error) {
					return hub.Result{
						StatusCode:      http.StatusOK,
						FullDescription: readme.Full,
					}, nil
				}).
				Twice()

			Expect(run(UpdateDockerReadme)).To(Succeed())
		})
	})
})

var _ = Describe("VerifyReadme", func() {
	BeforeEach(func() {
		*P = Pipe{DockerHub: DockerHub{Address: "https://hub.docker.com/v2/repositories"}}
	})

	readme := ParsedReadme{File: "README.md", Description: "a pipe"}

	It("accepts the readme the repository came back with", func() {
		Expect(VerifyReadme(hub.Result{
			StatusCode:      http.StatusOK,
			Description:     "a pipe",
			FullDescription: "# Pipe",
		}, "kilic/pipe", readme, "# Pipe")).To(Succeed())
	})

	It("rejects a full description that is not what was pushed", func() {
		Expect(VerifyReadme(hub.Result{
			StatusCode:      http.StatusOK,
			Description:     "a pipe",
			FullDescription: "# Something else",
		}, "kilic/pipe", readme, "# Pipe")).
			To(MatchError("Uploaded README does not match with current repository README file."))
	})

	It("rejects a short description that is not what was pushed", func() {
		Expect(VerifyReadme(hub.Result{
			StatusCode:      http.StatusOK,
			Description:     "something else",
			FullDescription: "# Pipe",
		}, "kilic/pipe", readme, "# Pipe")).
			To(MatchError("Uploaded README does not match with current repository README file."))
	})

	// The pipe leaves the short description alone when it was not given one, so
	// whatever the repository already carries is not a mismatch.
	It("ignores the short description the pipe did not push", func() {
		Expect(VerifyReadme(hub.Result{
			StatusCode:      http.StatusOK,
			Description:     "whatever was there",
			FullDescription: "# Pipe",
		}, "kilic/pipe", ParsedReadme{File: "README.md"}, "# Pipe")).To(Succeed())
	})

	It("names the repository that does not exist", func() {
		Expect(VerifyReadme(hub.Result{StatusCode: http.StatusNotFound}, "kilic/pipe", readme, "# Pipe")).
			To(MatchError("Repository does not exists: https://hub.docker.com/v2/repositories/kilic/pipe"))
	})

	// A repository the user can not edit fails with a status that does not say
	// which of the two went wrong, so the response is what points at the cause.
	It("blames the credentials when the user can not edit the repository", func() {
		Expect(VerifyReadme(
			hub.Result{StatusCode: http.StatusForbidden},
			"kilic/pipe",
			readme,
			"# Pipe",
		)).
			To(MatchError("Given user credentials do not have permission to edit repository: https://hub.docker.com/v2/repositories/kilic/pipe"))
	})

	It("reports the status code of a failure the credentials did not cause", func() {
		Expect(VerifyReadme(
			hub.Result{StatusCode: http.StatusInternalServerError, CanEdit: true},
			"kilic/pipe",
			readme,
			"# Pipe",
		)).
			To(MatchError("Pushing readme failed with code: 500"))
	})
})
