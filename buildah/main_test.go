package main

import (
	"testing"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/test/conformance"
	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
)

func TestPipe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Buildah Pipe Suite")
}

var _ = conformance.Verify(conformance.Pipe{
	Name:        name,
	Description: description,
	New: func(p *plumber.Plumber) *ucli.Command {
		return newCommand(p, "test", options{})
	},
	LegacyEnvAliases: map[string][]string{
		"buildah.build.file.context":         {"CONTAINER_FILE_CONTEXT", "BUILDAH_BUILD_FILE_CONTEXT"},
		"buildah.build.file.name":            {"CONTAINER_FILE_NAME", "BUILDAH_BUILD_FILE_NAME"},
		"buildah.build.image.build-args":     {"CONTAINER_IMAGE_BUILD_ARGS", "BUILDAH_BUILD_IMAGE_BUILD_ARGS"},
		"buildah.build.image.cache":          {"CONTAINER_IMAGE_CACHE", "BUILDAH_BUILD_IMAGE_CACHE"},
		"buildah.build.image.format":         {"CONTAINER_IMAGE_FORMAT", "BUILDAH_BUILD_IMAGE_FORMAT"},
		"buildah.build.image.latest-tag":     {"CONTAINER_IMAGE_LATEST_TAG", "BUILDAH_BUILD_IMAGE_LATEST_TAG"},
		"buildah.build.image.name":           {"CONTAINER_IMAGE_NAME", "BUILDAH_BUILD_IMAGE_NAME"},
		"buildah.build.image.platforms":      {"CONTAINER_IMAGE_PLATFORMS", "BUILDAH_BUILD_IMAGE_PLATFORMS"},
		"buildah.build.image.pull":           {"CONTAINER_IMAGE_PULL", "BUILDAH_BUILD_IMAGE_PULL"},
		"buildah.build.image.push":           {"CONTAINER_IMAGE_PUSH", "BUILDAH_BUILD_IMAGE_PUSH"},
		"buildah.build.image.storage-driver": {"CONTAINER_IMAGE_STORAGE_DRIVER", "BUILDAH_STORAGE_DRIVER", "BUILDAH_BUILD_IMAGE_STORAGE_DRIVER"},
		"buildah.build.image.tag-as-latest":  {"CONTAINER_IMAGE_TAGS_AS_LATEST", "BUILDAH_BUILD_IMAGE_TAG_AS_LATEST"},
		"buildah.build.image.tags":           {"CONTAINER_IMAGE_TAGS", "BUILDAH_BUILD_IMAGE_TAGS"},
		"buildah.build.image.tags-sanitize":  {"CONTAINER_IMAGE_SANITIZE_TAGS", "BUILDAH_BUILD_IMAGE_TAGS_SANITIZE"},
		"buildah.build.image.tags-template":  {"CONTAINER_IMAGE_TAGS_TEMPLATE", "BUILDAH_BUILD_IMAGE_TAGS_TEMPLATE"},
		"buildah.build.manifest.file":        {"CONTAINER_MANIFEST_FILE", "BUILDAH_BUILD_MANIFEST_FILE"},
		"buildah.build.manifest.target":      {"CONTAINER_MANIFEST_TARGET", "BUILDAH_BUILD_MANIFEST_TARGET"},
		"buildah.login.registry.password":    {"CONTAINER_REGISTRY_PASSWORD", "BUILDAH_LOGIN_REGISTRY_PASSWORD"},
		"buildah.login.registry.uri":         {"CONTAINER_REGISTRY_URI", "BUILDAH_LOGIN_REGISTRY_URI"},
		"buildah.login.registry.username":    {"CONTAINER_REGISTRY_USERNAME", "BUILDAH_LOGIN_REGISTRY_USERNAME"},
		"buildah.manifest.files":             {"CONTAINER_MANIFEST_FILES", "BUILDAH_MANIFEST_FILES"},
		"buildah.manifest.images":            {"CONTAINER_MANIFEST_IMAGES", "BUILDAH_MANIFEST_IMAGES"},
		"buildah.manifest.matrix":            {"CONTAINER_MANIFEST_MATRIX", "BUILDAH_MANIFEST_MATRIX"},
		"buildah.manifest.target":            {"CONTAINER_MANIFEST_TARGET", "BUILDAH_MANIFEST_TARGET"},
		"git.branch":                         {"CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH"},
		"git.tag":                            {"CI_COMMIT_TAG", "BITBUCKET_TAG"},
	},
})

// Everything a spec needs is passed as an argument rather than through the
// environment, since the flags are package level and urfave only reads an env
// source on the first parse of a flag instance: driving values in through the
// environment would make the specs depend on the order Ginkgo runs them in.
func run(runner *tests.TestingCommandRunner, args ...string) error {
	GinkgoHelper()

	fixture := fixtures.NewPlumber(func(p *plumber.Plumber) *ucli.Command {
		return newCommand(p, "test", options{})
	})
	fixture.Plumber.SetRuntime(plumber.Runtime{CommandRunner: runner.Runner()})

	return fixture.RunCli(append([]string{name}, args...)...)
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
	It("probes the tool and authenticates against the registry on login", func() {
		runner := fixtures.Runner()

		Expect(run(
			runner,
			"login",
			"--buildah.login.registry.uri", "registry.example.test",
			"--buildah.login.registry.username", "user",
			"--buildah.login.registry.password", "password",
		)).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"buildah", "buildah"}))
		// The password goes in over stdin, so it is deliberately absent from the
		// arguments the invocation records.
		Expect(formatted(runner)).To(ContainElement(
			ContainSubstring("login registry.example.test --username user --password-stdin"),
		))
	})

	// Without a username and a password the login task disables itself, which is
	// what a pipeline relying on an ambient login does. The registry the tags are
	// prefixed with is the default the login step declares.
	It("builds the image under the tags it was given on build", func() {
		tests.WithTempWorkingDirectory()

		runner := fixtures.Runner()

		Expect(run(
			runner,
			"build",
			"--buildah.build.image.name", "group/image",
			"--buildah.build.image.tags", "v1.0.0",
			"--buildah.build.image.push=false",
		)).To(Succeed())

		Expect(runner.InvocationNames()).To(Equal([]string{"buildah", "buildah"}))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("build --format oci --pull -t docker.io/group/image:v1.0.0 --file Dockerfile .")))
	})

	It("creates and pushes the manifest of the images it was given on manifest", func() {
		tests.WithTempWorkingDirectory()

		runner := fixtures.Runner()

		Expect(run(
			runner,
			"manifest",
			"--buildah.manifest.target", "group/image:v1.0.0",
			"--buildah.manifest.images", "group/image:v1.0.0-amd64",
		)).To(Succeed())

		Expect(formatted(runner)).To(ContainElement(ContainSubstring("manifest create docker.io/group/image:v1.0.0")))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("manifest add docker.io/group/image:v1.0.0 group/image:v1.0.0-amd64")))
		Expect(formatted(runner)).To(ContainElement(ContainSubstring("manifest push --rm docker.io/group/image:v1.0.0")))
	})
})
