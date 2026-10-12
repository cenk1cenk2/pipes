package build

import (
	"crypto/sha256"
	"encoding/hex"
	json "encoding/json/v2"
	"errors"
	"fmt"
	"os"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/buildah/login"
	buildahmanifest "gitlab.kilic.dev/devops/pipes/buildah/manifest"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

// a package level flag reads its environment only on the first parse, so the pipe is seeded.
func run(runner *tests.TestingCommandRunner, pipe Pipe, tasks ...func(*TaskList) *Task) error {
	GinkgoHelper()

	*P = pipe
	*C = Ctx{}
	login.P.Uri = ""

	return fixtures.Cli(runner, tests.TaskListCli{
		AppName:     "pipe-buildah",
		CommandName: "build",
		TaskLists: []tests.TaskListFactory{
			func(p *Plumber, _ *cli.Command) *TaskList {
				tl := &TaskList{}

				return tl.New(p).
					SetRuntimeDepth(3).
					Set(func(tl *TaskList) Job {
						jobs := []Job{}
						for _, task := range tasks {
							jobs = append(jobs, task(tl).Job())
						}

						return JobSequence(jobs...)
					})
			},
		},
	}).Run()
}

var _ = Describe("Container image tags", func() {
	format := func(uri, name, tag string) string {
		GinkgoHelper()

		P.Image.Name = name
		login.P.Uri = uri

		return tagsCollector().Format(tag)
	}

	It("prefixes the image with the registry the login step authenticated against", func() {
		Expect(format("registry.example.com", "group/image", "v1.0.0")).
			To(Equal("registry.example.com/group/image:v1.0.0"))
	})

	// a pipeline with no registry builds locally, where a prefix names nothing that is pushed.
	It("leaves the image unprefixed without a registry", func() {
		Expect(format("", "group/image", "v1.0.0")).To(Equal("group/image:v1.0.0"))
	})

	It("keeps a registry that already carries a port and a path", func() {
		Expect(format("registry.example.com:5000/mirror", "image", "latest")).
			To(Equal("registry.example.com:5000/mirror/image:latest"))
	})
})

var _ = Describe("Container image tags task", func() {
	It("builds with the tags the user gave and the release handed over", func() {
		GinkgoT().Setenv("SEMANTIC_RELEASE_TAG", "v1.0.0")

		Expect(run(fixtures.Runner(), Pipe{Image: Image{
			Name:                "group/image",
			Tags:                []string{"latest"},
			TagsFromEnvironment: []string{"SEMANTIC_RELEASE_TAG"},
		}}, tags)).To(Succeed())

		Expect(C.Tags).To(Equal([]string{"group/image:latest", "group/image:v1.0.0"}))
	})

	// a release that did not happen leaves the environment source with nothing, and
	// an image without a single tag can not be pushed anywhere.
	It("fails when no source gave a tag", func() {
		Expect(run(fixtures.Runner(), Pipe{Image: Image{
			Name:                "group/image",
			TagsFromEnvironment: []string{"SEMANTIC_RELEASE_TAG_UNSET"},
		}}, tags)).To(MatchError(ContainSubstring("Can not build the image without any tags.")))
	})
})

var _ = Describe("Container manifest dotenv report", func() {
	const envFile = "container-manifest.env"

	BeforeEach(func() {
		GinkgoT().Chdir(GinkgoT().TempDir())
	})

	pipe := func() Pipe {
		return Pipe{
			Image: Image{Name: "group/image", Tags: []string{"v1.0.0-amd64"}, Push: true},
			Manifest: Manifest{
				Target:          "v1.0.0",
				Job:             "build-amd64",
				EnvFile:         envFile,
				PublishedPrefix: "CONTAINER_MANIFEST_PUBLISHED_",
			},
		}
	}

	report := func() string {
		GinkgoHelper()

		content, err := os.ReadFile(envFile)
		Expect(err).NotTo(HaveOccurred())

		return string(content)
	}

	hash := func(tags ...string) string {
		sum := sha256.Sum256([]byte(strings.Join(tags, ",")))

		return strings.ToUpper(hex.EncodeToString(sum[:])[:8])
	}

	It("hands the pushed images over keyed by the job and the tags", func() {
		Expect(run(fixtures.Runner(), pipe(), tags, push, manifest)).To(Succeed())

		Expect(report()).To(Equal(fmt.Sprintf(
			"CONTAINER_MANIFEST_PUBLISHED_BUILD_AMD64_%s=%s\n",
			hash("group/image:v1.0.0-amd64"),
			`[{"target":"group/image:v1.0.0","images":["group/image:v1.0.0-amd64"]}]`,
		)))
	})

	// the manifest job reads the value back into the very same matrix.
	It("hands over a value the manifest job reads back", func() {
		Expect(run(fixtures.Runner(), pipe(), tags, push, manifest)).To(Succeed())

		_, value, _ := strings.Cut(strings.TrimSpace(report()), "=")
		parsed := []buildahmanifest.ManifestMatrix{}

		Expect(json.Unmarshal([]byte(value), &parsed)).To(Succeed())
		Expect(parsed).To(Equal([]buildahmanifest.ManifestMatrix{
			{Target: "group/image:v1.0.0", Images: []string{"group/image:v1.0.0-amd64"}},
		}))
	})

	It("writes the report and keys it as the pipeline asks", func() {
		p := pipe()
		p.Manifest.EnvFile = "published.env"
		p.Manifest.PublishedPrefix = "PUBLISHED_"

		Expect(run(fixtures.Runner(), p, tags, push, manifest)).To(Succeed())

		content, err := os.ReadFile("published.env")
		Expect(err).NotTo(HaveOccurred())
		Expect(string(content)).To(HavePrefix(fmt.Sprintf("PUBLISHED_BUILD_AMD64_%s=", hash("group/image:v1.0.0-amd64"))))
	})

	It("collapses the separators of a matrix job slug into one", func() {
		p := pipe()
		p.Manifest.Job = "container---dockerfile--amd64"

		Expect(run(fixtures.Runner(), p, tags, push, manifest)).To(Succeed())

		Expect(report()).To(HavePrefix(fmt.Sprintf("CONTAINER_MANIFEST_PUBLISHED_CONTAINER_DOCKERFILE_AMD64_%s=", hash("group/image:v1.0.0-amd64"))))
	})

	// a local run has no job name, and the tags alone still keep two builds apart.
	It("keys the images by the tags alone without a job", func() {
		p := pipe()
		p.Manifest.Job = ""

		Expect(run(fixtures.Runner(), p, tags, push, manifest)).To(Succeed())

		Expect(report()).To(HavePrefix(fmt.Sprintf("CONTAINER_MANIFEST_PUBLISHED_%s=", hash("group/image:v1.0.0-amd64"))))
	})

	// the target goes through the same formatting as a tag, so it names the bare tag.
	It("renders the target against the tags", func() {
		p := pipe()
		p.Manifest.Target = `{{ index $ 0 | trimPrefix "group/image:" | trimSuffix "-amd64" }}`

		Expect(run(fixtures.Runner(), p, tags, push, manifest)).To(Succeed())

		Expect(report()).To(ContainSubstring(`"target":"group/image:v1.0.0"`))
	})

	It("writes nothing when the push failed", func() {
		runner := fixtures.Runner(tests.TestingCommandResponse{Name: "buildah", Err: errors.New("push failed")})

		Expect(run(runner, pipe(), tags, push, manifest)).NotTo(Succeed())

		Expect(envFile).NotTo(BeAnExistingFile())
	})

	It("writes nothing when the images are not pushed", func() {
		p := pipe()
		p.Image.Push = false

		Expect(run(fixtures.Runner(), p, tags, push, manifest)).To(Succeed())

		Expect(envFile).NotTo(BeAnExistingFile())
	})

	It("writes nothing without a manifest target", func() {
		p := pipe()
		p.Manifest.Target = ""

		Expect(run(fixtures.Runner(), p, tags, push, manifest)).To(Succeed())

		Expect(envFile).NotTo(BeAnExistingFile())
	})
})
