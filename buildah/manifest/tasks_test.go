package manifest

import (
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Container manifest", func() {
	var (
		p  *Plumber
		tl *TaskList
	)

	BeforeEach(func() {
		p = NewPlumber(func(_ *Plumber) *cli.Command {
			return &cli.Command{Name: "test"}
		})
		p.SetLoggerOutput(GinkgoWriter)

		tl = &TaskList{}
		tl.New(p)

		*P = Pipe{}
		*C = Ctx{ManifestedImages: map[string][]string{}}
	})

	Describe("user input", func() {
		It("renders the target the images were given for", func() {
			P.Manifest = Manifest{Target: `{{ "registry.example.com/image" }}`, Images: []string{"registry.example.com/image:amd64"}}

			Expect(p.RunJobs(fetchUser(tl).Job())).To(Succeed())

			Expect(C.ManifestedImages).
				To(Equal(map[string][]string{"registry.example.com/image": {"registry.example.com/image:amd64"}}))
		})

		It("takes the matrix without any images", func() {
			P.Manifest = Manifest{Matrix: []ManifestMatrix{{Target: "registry.example.com/image", Images: []string{"registry.example.com/image:amd64"}}}}

			Expect(p.RunJobs(fetchUser(tl).Job())).To(Succeed())

			Expect(C.ManifestedImages).
				To(HaveKeyWithValue("registry.example.com/image", []string{"registry.example.com/image:amd64"}))
		})

		// the fetchers run in parallel under the one task list lock, so a failed
		// render holding it would hang every other fetcher.
		It("releases the lock when the target does not render", func() {
			P.Manifest = Manifest{Target: "{{ .Missing", Images: []string{"registry.example.com/image:amd64"}}

			Expect(p.RunJobs(fetchUser(tl).Job())).NotTo(Succeed())

			Expect(tl.Lock.TryLock()).To(BeTrue())
		})
	})

	Describe("published images", func() {
		// a package level flag reads its environment only on the first parse, so the pipe is seeded.
		run := func(runner *tests.TestingCommandRunner, pipe Pipe) error {
			GinkgoHelper()

			*P = pipe

			return fixtures.Cli(runner, tests.TaskListCli{
				AppName:     "pipe-buildah",
				CommandName: "manifest",
				TaskLists: []tests.TaskListFactory{
					func(p *Plumber, _ *cli.Command) *TaskList {
						return New(p)
					},
				},
			}).Run()
		}

		published := func(manifest Manifest) Pipe {
			manifest.PublishedPrefix = "CONTAINER_MANIFEST_PUBLISHED_"

			return Pipe{Manifest: manifest}
		}

		commands := func(runner *tests.TestingCommandRunner) [][]string {
			args := [][]string{}
			for _, invocation := range runner.Invocations() {
				args = append(args, invocation.Args)
			}

			return args
		}

		It("manifests the images every build job handed over", func() {
			GinkgoT().Setenv("CONTAINER_MANIFEST_PUBLISHED_BUILD_AMD64_0A1B2C3D", `[{"target":"registry.example.com/group/image:v1","images":["registry.example.com/group/image:v1-amd64"]}]`)
			GinkgoT().Setenv("CONTAINER_MANIFEST_PUBLISHED_BUILD_ARM64_4E5F6A7B", `[{"target":"registry.example.com/group/image:v1","images":["registry.example.com/group/image:v1-arm64"]}]`)
			GinkgoT().Setenv("CONTAINER_MANIFEST_OTHER", "not a matrix")

			runner := fixtures.Runner()

			Expect(run(runner, published(Manifest{}))).To(Succeed())

			Expect(commands(runner)).To(Equal([][]string{
				{"manifest", "create", "registry.example.com/group/image:v1"},
				{"manifest", "add", "registry.example.com/group/image:v1", "registry.example.com/group/image:v1-amd64"},
				{"manifest", "add", "registry.example.com/group/image:v1", "registry.example.com/group/image:v1-arm64"},
				{"manifest", "push", "--rm", "registry.example.com/group/image:v1"},
			}))
		})

		It("reads the variables under the prefix the pipeline asks for", func() {
			GinkgoT().Setenv("PUBLISHED_BUILD_AMD64", `[{"target":"registry.example.com/group/image:v1","images":["registry.example.com/group/image:v1-amd64"]}]`)

			runner := fixtures.Runner()

			Expect(run(runner, Pipe{Manifest: Manifest{PublishedPrefix: "PUBLISHED_"}})).To(Succeed())

			Expect(commands(runner)).To(ContainElement([]string{"manifest", "add", "registry.example.com/group/image:v1", "registry.example.com/group/image:v1-amd64"}))
		})

		It("adds an image handed over twice only once", func() {
			GinkgoT().Setenv("CONTAINER_MANIFEST_PUBLISHED_BUILD_AMD64_0A1B2C3D", `[{"target":"registry.example.com/group/image:v1","images":["registry.example.com/group/image:v1-amd64"]}]`)
			GinkgoT().Setenv("CONTAINER_MANIFEST_PUBLISHED_BUILD_AMD64_4E5F6A7B", `[{"target":"registry.example.com/group/image:v1","images":["registry.example.com/group/image:v1-amd64"]}]`)

			runner := fixtures.Runner()

			Expect(run(runner, published(Manifest{
				Target: "registry.example.com/group/image:v1",
				Images: []string{"registry.example.com/group/image:v1-amd64"},
			}))).To(Succeed())

			Expect(commands(runner)).To(Equal([][]string{
				{"manifest", "create", "registry.example.com/group/image:v1"},
				{"manifest", "add", "registry.example.com/group/image:v1", "registry.example.com/group/image:v1-amd64"},
				{"manifest", "push", "--rm", "registry.example.com/group/image:v1"},
			}))
		})

		// the login no longer names the target, so one without its registry would
		// leave buildah to guess where the manifest goes.
		It("refuses a target without the registry host", func() {
			GinkgoT().Setenv("CONTAINER_MANIFEST_PUBLISHED_BUILD_AMD64_0A1B2C3D", `[{"target":"group/image:v1","images":["group/image:v1-amd64"]}]`)

			runner := fixtures.Runner()

			Expect(run(runner, published(Manifest{}))).To(MatchError(ContainSubstring("Image name is not fully qualified")))
			Expect(runner.Invocations()).To(BeEmpty())
		})

		It("reports a variable that does not hold a matrix", func() {
			GinkgoT().Setenv("CONTAINER_MANIFEST_PUBLISHED_BUILD_AMD64_0A1B2C3D", "not a matrix")

			Expect(run(fixtures.Runner(), published(Manifest{}))).
				To(MatchError(ContainSubstring("Can not unmarshal container manifest matrix")))
		})

		// a build job whose component does not upload the report hands nothing over,
		// which would otherwise pass as a manifest with nothing to do.
		It("fails without any image to manifest", func() {
			runner := fixtures.Runner()

			Expect(run(runner, published(Manifest{}))).To(MatchError(ContainSubstring("Can not find any published images to manifest.")))
			Expect(runner.Invocations()).To(BeEmpty())
		})
	})
})
