package client

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/github/test/keys"
)

var _ = Describe("Config", func() {
	DescribeTable(
		"is enabled by any of the app credentials",
		func(cfg Config, expected bool) {
			Expect(cfg.Enabled()).To(Equal(expected))
		},
		Entry("none", Config{ApiUrl: "https://api.github.com"}, false),
		Entry("id", Config{Id: "123"}, true),
		Entry("installation id", Config{InstallationId: "456"}, true),
		Entry("private key", Config{PrivateKey: "key.pem"}, true),
	)

	// a partial set of credentials is a misconfiguration, and falling back to the
	// ready token would hide it until that token expires.
	DescribeTable(
		"takes the credentials only as a complete set",
		func(cfg Config, valid bool) {
			cfg.ApiUrl = "https://api.github.com"

			err := tests.NewPlumber().Plumber.Validate(&cfg)

			if valid {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(HaveOccurred())
			}
		},
		Entry("none", Config{}, true),
		Entry("all", Config{Id: "123", InstallationId: "456", PrivateKey: "key.pem"}, true),
		Entry("no private key", Config{Id: "123", InstallationId: "456"}, false),
		Entry("no installation id", Config{Id: "123", PrivateKey: "key.pem"}, false),
		Entry("id alone", Config{Id: "123"}, false),
	)

	Describe("Mint", func() {
		var (
			pem           string
			authorization string
			plumber       *Plumber
			output        *bytes.Buffer
			subject       ApplicationClientAdapter
		)

		BeforeEach(func() {
			_, pem = keys.Generate()
			authorization = ""

			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					authorization = r.Header.Get("Authorization")

					w.WriteHeader(http.StatusCreated)
					_, _ = w.Write([]byte(`{"token":"ghs_minted"}`))
				}),
			)
			DeferCleanup(server.Close)

			output = &bytes.Buffer{}
			plumber = tests.NewPlumber().Plumber.SetLoggerOutput(output)
			subject = NewApplicationClient(server.URL, "pipe-github")
		})

		mint := func(key string) (string, error) {
			GinkgoHelper()

			cfg := Config{Id: "123", InstallationId: "456", PrivateKey: key}

			return cfg.Mint(context.Background(), plumber, subject, TokenRequest{})
		}

		It("reads the key from the path a GitLab file variable exports", func() {
			file := filepath.Join(GinkgoT().TempDir(), "key.pem")
			Expect(os.WriteFile(file, []byte(pem), 0o600)).To(Succeed())

			Expect(mint(file)).To(Equal("ghs_minted"))
		})

		It("takes the PEM contents in place of a path", func() {
			Expect(mint(pem)).To(Equal("ghs_minted"))
		})

		It("authenticates the request with a JWT", func() {
			_, err := mint(pem)
			Expect(err).NotTo(HaveOccurred())

			Expect(authorization).To(HavePrefix("Bearer "))
			Expect(strings.Split(strings.TrimPrefix(authorization, "Bearer "), ".")).To(HaveLen(3))
		})

		It("fails on a key path that is not there", func() {
			_, err := mint(filepath.Join(GinkgoT().TempDir(), "missing.pem"))

			Expect(err).To(HaveOccurred())
		})

		It("keeps the key, the JWT and the token out of the log", func() {
			_, err := mint(pem)
			Expect(err).NotTo(HaveOccurred())

			jwt := strings.TrimPrefix(authorization, "Bearer ")

			plumber.Log.Info(strings.Join([]string{pem, jwt, "ghs_minted"}, " "))

			// the logger re-indents a multi-line message, so the key is asserted on
			// line by line or the check would pass on the reformatting alone.
			for _, line := range strings.Split(strings.TrimSpace(pem), "\n")[1:] {
				if strings.HasPrefix(line, "-----END") {
					continue
				}

				Expect(output.String()).NotTo(ContainSubstring(line))
			}

			Expect(output.String()).NotTo(ContainSubstring(jwt))
			Expect(output.String()).NotTo(ContainSubstring("ghs_minted"))
		})
	})
})
