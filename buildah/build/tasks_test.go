package build

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/buildah/login"
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
