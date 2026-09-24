package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	var (
		request *http.Request
		body    string
		answer  func(w http.ResponseWriter)
		subject ApplicationClientAdapter
	)

	BeforeEach(func() {
		request = nil
		body = ""

		server := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				read, err := io.ReadAll(r.Body)
				Expect(err).NotTo(HaveOccurred())

				request = r
				body = string(read)

				answer(w)
			}),
		)
		DeferCleanup(server.Close)

		subject = NewApplicationClient(server.URL+"/", "pipe-github")
	})

	Describe("CreateInstallationToken", func() {
		BeforeEach(func() {
			answer = func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"token":"ghs_token","expires_at":"2026-09-24T12:00:00Z"}`))
			}
		})

		create := func(request TokenRequest) (string, error) {
			GinkgoHelper()

			return subject.CreateInstallationToken(context.Background(), "app-jwt", "456", request)
		}

		It("posts to the access tokens of the installation", func() {
			_, err := create(TokenRequest{})

			Expect(err).NotTo(HaveOccurred())
			Expect(request.Method).To(Equal(http.MethodPost))
			Expect(request.URL.Path).To(Equal("/app/installations/456/access_tokens"))
		})

		It("authenticates as the app with the headers GitHub expects", func() {
			_, err := create(TokenRequest{})

			Expect(err).NotTo(HaveOccurred())
			Expect(request.Header.Get("Authorization")).To(Equal("Bearer app-jwt"))
			Expect(request.Header.Get("Accept")).To(Equal("application/vnd.github+json"))
			Expect(request.Header.Get("X-GitHub-Api-Version")).To(Equal("2022-11-28"))
			Expect(request.Header.Get("Content-Type")).To(Equal(JSONRequest))
			Expect(request.Header.Get("User-Agent")).To(Equal("pipe-github"))
		})

		It("narrows the token down to the repositories and permissions it was given", func() {
			_, err := create(TokenRequest{
				Repositories: []string{"listr2"},
				Permissions:  map[string]string{"statuses": "write"},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(Equal(`{"repositories":["listr2"],"permissions":{"statuses":"write"}}`))
		})

		// an empty list would narrow the token down to no repository at all, while
		// leaving the field out keeps everything the installation was granted.
		It("leaves out what it was not asked to narrow down", func() {
			_, err := create(TokenRequest{})

			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(Equal(`{}`))
		})

		It("hands back the token", func() {
			Expect(create(TokenRequest{})).To(Equal("ghs_token"))
		})

		It("surfaces only the message of a rejected request", func() {
			answer = func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"message":"A JSON web token could not be decoded","documentation_url":"https://docs.github.com"}`))
			}

			_, err := create(TokenRequest{})

			Expect(err).To(MatchError("Can not create installation token: GitHub responded with code: 401 > A JSON web token could not be decoded"))
			Expect(err.Error()).NotTo(ContainSubstring("app-jwt"))
		})

		It("fails when the response carries no token", func() {
			answer = func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{}`))
			}

			_, err := create(TokenRequest{})

			Expect(err).To(MatchError("Can not create installation token, response carries no token."))
		})
	})

	Describe("CreateCommitStatus", func() {
		BeforeEach(func() {
			answer = func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"id":1,"state":"success"}`))
			}
		})

		create := func(status CommitStatus) error {
			GinkgoHelper()

			return subject.CreateCommitStatus(context.Background(), "ghs_token", "burningforge/listr2", "abc123", status)
		}

		It("posts to the statuses of the commit", func() {
			Expect(create(CommitStatus{State: "success", Context: "Gitlab CI"})).To(Succeed())

			Expect(request.Method).To(Equal(http.MethodPost))
			Expect(request.URL.Path).To(Equal("/repos/burningforge/listr2/statuses/abc123"))
		})

		It("authenticates with the token", func() {
			Expect(create(CommitStatus{State: "success", Context: "Gitlab CI"})).To(Succeed())

			Expect(request.Header.Get("Authorization")).To(Equal("Bearer ghs_token"))
			Expect(request.Header.Get("Accept")).To(Equal("application/vnd.github+json"))
			Expect(request.Header.Get("X-GitHub-Api-Version")).To(Equal("2022-11-28"))
		})

		It("sends the status as GitHub names its fields", func() {
			Expect(create(CommitStatus{
				State:       "pending",
				TargetUrl:   "https://gitlab.kilic.dev/pipelines/1",
				Description: "running",
				Context:     "Gitlab CI",
			})).To(Succeed())

			Expect(body).To(Equal(
				`{"state":"pending","target_url":"https://gitlab.kilic.dev/pipelines/1","description":"running","context":"Gitlab CI"}`,
			))
		})

		It("surfaces only the message of a rejected status", func() {
			answer = func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"message":"Not Found"}`))
			}

			Expect(create(CommitStatus{State: "success", Context: "Gitlab CI"})).
				To(MatchError("Can not create commit status: GitHub responded with code: 404 > Not Found"))
		})

		DescribeTable(
			"tells a commit GitHub does not have apart from other rejections",
			func(code int, message string, expected bool) {
				answer = func(w http.ResponseWriter) {
					w.WriteHeader(code)
					_, _ = w.Write([]byte(`{"message":"` + message + `"}`))
				}

				var response *ResponseError
				Expect(errors.As(create(CommitStatus{State: "success", Context: "Gitlab CI"}), &response)).To(BeTrue())

				Expect(response.CommitNotFound()).To(Equal(expected))
			},
			Entry("missing commit", http.StatusUnprocessableEntity, "No commit found for SHA: abc123", true),
			Entry("other validation", http.StatusUnprocessableEntity, "Validation Failed", false),
			Entry("missing repository", http.StatusNotFound, "Not Found", false),
			Entry("forbidden", http.StatusForbidden, "Resource not accessible by integration", false),
		)
	})
})
