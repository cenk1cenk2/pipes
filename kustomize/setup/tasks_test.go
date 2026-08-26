package setup

import (
	"github.com/cenk1cenk2/plumber/v6"
	"github.com/cenk1cenk2/plumber/v6/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/test/fixtures"
	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

var _ = Describe("Resolve overlays", func() {
	// The flags are not registered with the spec command, since the pipe is seeded
	// directly and a package level flag only reads its environment on first parse.
	resolve := func(cwd string, paths ...string) []string {
		GinkgoHelper()

		P.Paths = paths
		C.Ctx = tool.NewCtx()
		C.Cwd = cwd
		C.Overlays = nil

		Expect(fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
			AppName:     "pipe-kustomize",
			CommandName: "build",
			TaskLists: []tests.TaskListFactory{
				func(p *plumber.Plumber, _ *ucli.Command) *plumber.TaskList {
					tl := &plumber.TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *plumber.TaskList) plumber.Job {
							return plumber.JobSequence(ResolveOverlays(tl).Job())
						})
				},
			},
		}).Run()).To(Succeed())

		return C.Overlays
	}

	It("builds the working directory itself when no paths were given", func() {
		Expect(resolve("overlays/production")).To(Equal([]string{"overlays/production"}))
	})

	// The working directory flag defaults to ".", but a pipeline that unsets it
	// would otherwise resolve an empty overlay path that Kustomize cannot read.
	It("falls back to the current directory", func() {
		Expect(resolve("")).To(Equal([]string{"."}))
	})

	It("resolves the explicit paths against the working directory", func() {
		Expect(resolve("clusters/prod", "apps/api", "apps/web")).
			To(Equal([]string{"clusters/prod/apps/api", "clusters/prod/apps/web"}))
	})

	// The same overlay reaching the build twice would render it twice and write the
	// output file from two subtasks at once.
	It("sorts the paths and drops the duplicates", func() {
		Expect(resolve(".", "apps/web", "apps/api", "apps/web")).
			To(Equal([]string{"apps/api", "apps/web"}))
	})
})
