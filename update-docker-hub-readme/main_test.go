package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/test/conformance"
	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	"gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/hub"
	mockhub "gitlab.kilic.dev/devops/pipes/update-docker-hub-readme/test/mocks/hub"
)

func TestPipe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Update Docker Hub Readme Pipe Suite")
}

var _ = conformance.Verify(conformance.Pipe{
	Name:        name,
	Description: description,
	New: func(p *plumber.Plumber) *ucli.Command {
		return newCommand(p, "test", options{})
	},
	LegacyEnvAliases: map[string][]string{
		"docker-hub.password":           {"DOCKER_PASSWORD", "DOCKER_HUB_PASSWORD"},
		"docker-hub.readme.description": {"README_SHORT_DESCRIPTION", "DOCKER_HUB_README_DESCRIPTION"},
		"docker-hub.readme.file":        {"README_FILE", "DOCKER_HUB_README_FILE"},
		"docker-hub.readme.matrix":      {"README_MATRIX", "DOCKER_HUB_README_MATRIX"},
		"docker-hub.readme.repository":  {"DOCKER_IMAGE_NAME", "CONTAINER_IMAGE_NAME", "README_REPOSITORY", "DOCKER_HUB_README_REPOSITORY"},
		"docker-hub.username":           {"DOCKER_USERNAME", "DOCKER_HUB_USERNAME"},
	},
})

// Everything a spec needs is passed as an argument rather than through the
// environment, since the flags are package level and urfave only reads an env
// source on the first parse of a flag instance: driving values in through the
// environment would make the specs depend on the order Ginkgo runs them in.
func run(opts options, args ...string) error {
	GinkgoHelper()

	fixture := fixtures.NewPlumber(func(p *plumber.Plumber) *ucli.Command {
		return newCommand(p, "test", opts)
	})
	fixture.Plumber.SetRuntime(plumber.Runtime{CommandRunner: fixtures.Runner().Runner()})

	return fixture.RunCli(append([]string{name}, args...)...)
}

var _ = Describe("New", func() {
	It("pushes the readme through the client the options supplied", func() {
		readme := filepath.Join(tests.TempDir(), "README.md")
		Expect(os.WriteFile(readme, []byte("# the readme\n"), 0o644)).To(Succeed())

		client := mockhub.NewMockClient(GinkgoT())
		client.EXPECT().
			Login(mock.Anything, "user", "password").
			Return("jwt-token", nil)

		var pushed hub.Readme
		client.EXPECT().
			UpdateReadme(mock.Anything, "jwt-token", "cenk1cenk2/pipe-go", mock.Anything).
			RunAndReturn(func(_ context.Context, _, _ string, r hub.Readme) (hub.Result, error) {
				pushed = r

				return hub.Result{
					CanEdit:         true,
					Description:     r.Description,
					FullDescription: r.Full,
					StatusCode:      http.StatusOK,
				}, nil
			})

		var dialed []string

		Expect(run(
			options{Hub: func(address, userAgent string) hub.Client {
				dialed = []string{address, userAgent}

				return client
			}},
			"--docker-hub.username", "user",
			"--docker-hub.password", "password",
			"--docker-hub.address", "https://hub.example.test/v2/repositories",
			"--docker-hub.readme.repository", "cenk1cenk2/pipe-go",
			"--docker-hub.readme.file", readme,
			"--docker-hub.readme.description", "Golang operations for pipelines.",
		)).To(Succeed())

		// The user agent is the name of the pipe, which is what makes the dial worth
		// asserting on here rather than in the update package.
		Expect(dialed).To(Equal([]string{"https://hub.example.test/v2/repositories", "pipe-update-docker-hub-readme"}))
		Expect(pushed).To(Equal(hub.Readme{
			Description: "Golang operations for pipelines.",
			Full:        "# the readme\n",
		}))
	})
})
