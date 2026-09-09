package update

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/hub"
	mockhub "gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/test/mocks/hub"
)

var _ = Describe("Docker Hub readme", func() {
	var client *mockhub.MockClientAdapter

	BeforeEach(func() {
		client = mockhub.NewMockClientAdapter(GinkgoT())

		*P = Pipe{DockerHub: DockerHub{
			Username: "user",
			Password: "password",
			Address:  "https://hub.docker.com/v2/repositories",
		}}
		*C = Ctx{Readme: map[string]ParsedReadme{}, Hub: client}
	})

	// a package level flag reads its environment only on the first parse, so the pipe is seeded.
	// a failing task terminates the process through the plumber, so the errors are
	// asserted on verifyReadme instead.
	run := func(tasks ...func(*TaskList) *Task) error {
		GinkgoHelper()

		return fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
			AppName:     "pipe-update-docker-hub-readme",
			CommandName: "update",
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

	// writes a readme file and registers it as the discovered target, since the
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

			Expect(run(login)).To(Succeed())
			Expect(C.Token).To(Equal("jwt-token"))
		})
	})

	Describe("discover", func() {
		It("takes the single repository the pipe was given", func() {
			P.Readme = Readme{Repository: "kilic/pipe", File: "README.md", Description: "a pipe"}

			Expect(run(discover)).To(Succeed())
			Expect(C.Readme).To(Equal(map[string]ParsedReadme{
				"kilic/pipe": {File: "README.md", Description: "a pipe"},
			}))
		})

		It("takes every repository of the matrix", func() {
			P.Readme = Readme{Matrix: []ReadmeMatrixJson{
				{Repository: "kilic/one", File: "one.md"},
				{Repository: "kilic/two", File: "two.md", Description: "the second"},
			}}

			Expect(run(discover)).To(Succeed())
			Expect(C.Readme).To(Equal(map[string]ParsedReadme{
				"kilic/one": {File: "one.md"},
				"kilic/two": {File: "two.md", Description: "the second"},
			}))
		})

		// the matrix is the way to update more than one repository in a job, so a
		// pipeline that sets both should not have to drop the single target.
		It("takes the single repository alongside the matrix", func() {
			P.Readme = Readme{
				Repository: "kilic/pipe",
				File:       "README.md",
				Matrix:     []ReadmeMatrixJson{{Repository: "kilic/one", File: "one.md"}},
			}

			Expect(run(discover)).To(Succeed())
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

			Expect(run(update)).To(Succeed())
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

			Expect(run(update)).To(Succeed())
		})
	})
})

var _ = Describe("verifyReadme", func() {
	BeforeEach(func() {
		*P = Pipe{DockerHub: DockerHub{Address: "https://hub.docker.com/v2/repositories"}}
	})

	readme := ParsedReadme{File: "README.md", Description: "a pipe"}

	It("accepts the readme the repository came back with", func() {
		Expect(verifyReadme(hub.Result{
			StatusCode:      http.StatusOK,
			Description:     "a pipe",
			FullDescription: "# Pipe",
		}, "kilic/pipe", readme, "# Pipe")).To(Succeed())
	})

	It("rejects a full description that is not what was pushed", func() {
		Expect(verifyReadme(hub.Result{
			StatusCode:      http.StatusOK,
			Description:     "a pipe",
			FullDescription: "# Something else",
		}, "kilic/pipe", readme, "# Pipe")).
			To(MatchError("Uploaded README does not match with current repository README file."))
	})

	It("rejects a short description that is not what was pushed", func() {
		Expect(verifyReadme(hub.Result{
			StatusCode:      http.StatusOK,
			Description:     "something else",
			FullDescription: "# Pipe",
		}, "kilic/pipe", readme, "# Pipe")).
			To(MatchError("Uploaded README does not match with current repository README file."))
	})

	// the pipe leaves the short description alone when it was not given one, so
	// whatever the repository already carries is not a mismatch.
	It("ignores the short description the pipe did not push", func() {
		Expect(verifyReadme(hub.Result{
			StatusCode:      http.StatusOK,
			Description:     "whatever was there",
			FullDescription: "# Pipe",
		}, "kilic/pipe", ParsedReadme{File: "README.md"}, "# Pipe")).To(Succeed())
	})

	It("names the repository that does not exist", func() {
		Expect(verifyReadme(hub.Result{StatusCode: http.StatusNotFound}, "kilic/pipe", readme, "# Pipe")).
			To(MatchError("Repository does not exists: https://hub.docker.com/v2/repositories/kilic/pipe"))
	})

	// a repository the user can not edit fails with a status that does not say
	// which of the two went wrong, so the response is what points at the cause.
	It("blames the credentials when the user can not edit the repository", func() {
		Expect(verifyReadme(
			hub.Result{StatusCode: http.StatusForbidden},
			"kilic/pipe",
			readme,
			"# Pipe",
		)).
			To(MatchError("Given user credentials do not have permission to edit repository: https://hub.docker.com/v2/repositories/kilic/pipe"))
	})

	It("reports the status code of a failure the credentials did not cause", func() {
		Expect(verifyReadme(
			hub.Result{StatusCode: http.StatusInternalServerError, CanEdit: true},
			"kilic/pipe",
			readme,
			"# Pipe",
		)).
			To(MatchError("Pushing readme failed with code: 500"))
	})
})
