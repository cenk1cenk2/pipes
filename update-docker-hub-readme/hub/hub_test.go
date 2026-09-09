package hub

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Login", func() {
	var (
		request *http.Request
		body    string
		answer  func(w http.ResponseWriter)
		subject *client
	)

	BeforeEach(func() {
		request = nil
		body = ""
		answer = func(w http.ResponseWriter) {
			_, _ = w.Write([]byte(`{"access_token":"jwt-token"}`))
		}

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

		subject = NewClient("", "pipe-update-docker-hub-readme").(*client)
		subject.loginAddress = server.URL + "/v2/auth/token"
	})

	It("posts the credentials as the account API expects them", func() {
		_, err := subject.Login(context.Background(), "user", "dckr_pat_secret")

		Expect(err).NotTo(HaveOccurred())
		Expect(request.Method).To(Equal(http.MethodPost))
		Expect(request.Header.Get("Content-Type")).To(Equal(JSONRequest))
		Expect(body).To(Equal(`{"identifier":"user","secret":"dckr_pat_secret"}`))
	})

	It("hands back the token", func() {
		Expect(subject.Login(context.Background(), "user", "dckr_pat_secret")).To(Equal("jwt-token"))
	})

	// the service answers a JSON error body on rejected credentials, and decoding it
	// as a token leaves nothing to tell the user what went wrong.
	It("carries the response into the error when the credentials are rejected", func() {
		answer = func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"message":"unauthorized","errinfo":{}}`))
		}

		_, err := subject.Login(context.Background(), "user", "dckr_pat_secret")

		Expect(err).To(MatchError(ContainSubstring(`{"message":"unauthorized","errinfo":{}}`)))
	})

	It("fails when the response is not the expected shape", func() {
		answer = func(w http.ResponseWriter) {
			_, _ = w.Write([]byte("<html>unauthorized</html>"))
		}

		_, err := subject.Login(context.Background(), "user", "dckr_pat_secret")

		Expect(err).To(HaveOccurred())
	})

	It("fails when the response carries no token", func() {
		answer = func(w http.ResponseWriter) {
			_, _ = w.Write([]byte(`{"detail":"Incorrect authentication credentials"}`))
		}

		_, err := subject.Login(context.Background(), "user", "dckr_pat_secret")

		Expect(err).To(MatchError(ContainSubstring("Incorrect authentication credentials")))
	})
})

var _ = Describe("UpdateReadme", func() {
	var (
		request *http.Request
		body    string
		answer  func(w http.ResponseWriter)
		subject ClientAdapter
	)

	BeforeEach(func() {
		request = nil
		body = ""
		answer = func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"can_edit":true,"description":"a pipe","full_description":"# Pipe"}`))
		}

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

		subject = NewClient(server.URL+"/v2/repositories", "pipe-update-docker-hub-readme")
	})

	update := func() (Result, error) {
		GinkgoHelper()

		return subject.UpdateReadme(
			context.Background(),
			"jwt-token",
			"kilic/pipe",
			Readme{Description: "a pipe", Full: "# Pipe"},
		)
	}

	It("patches the repository under the configured address", func() {
		_, err := update()

		Expect(err).NotTo(HaveOccurred())
		Expect(request.Method).To(Equal(http.MethodPatch))
		Expect(request.URL.Path).To(Equal("/v2/repositories/kilic/pipe/"))
	})

	It("sends both descriptions the repository page shows", func() {
		_, err := update()

		Expect(err).NotTo(HaveOccurred())
		Expect(body).To(Equal(`{"description":"a pipe","full_description":"# Pipe"}`))
	})

	It("authenticates with the token the login returned", func() {
		_, err := update()

		Expect(err).NotTo(HaveOccurred())
		Expect(request.Header.Get("Authorization")).To(Equal("JWT jwt-token"))
		Expect(request.Header.Get("Content-Type")).To(Equal(JSONRequest))
		Expect(request.Header.Get("User-Agent")).To(Equal("pipe-update-docker-hub-readme"))
	})

	It("hands back what the repository ended up with", func() {
		result, err := update()

		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(Result{
			CanEdit:         true,
			Description:     "a pipe",
			FullDescription: "# Pipe",
			StatusCode:      http.StatusOK,
		}))
	})

	// the service answers HTML on some failures, and a decode error alone leaves
	// nothing to tell the user what went wrong.
	It("carries the response into the error when it is not the expected shape", func() {
		answer = func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("<html>bad gateway</html>"))
		}

		_, err := update()

		Expect(err).To(MatchError(ContainSubstring("<html>bad gateway</html>")))
	})

	It("keeps the status code of a failure the caller has to decide on", func() {
		answer = func(w http.ResponseWriter) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"can_edit":false}`))
		}

		result, err := update()

		Expect(err).NotTo(HaveOccurred())
		Expect(result.StatusCode).To(Equal(http.StatusForbidden))
		Expect(result.CanEdit).To(BeFalse())
	})
})
