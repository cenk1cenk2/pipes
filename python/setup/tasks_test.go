package setup

import (
	"path/filepath"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Python environment", func() {
	run := func(cache string) error {
		GinkgoHelper()

		*P = Pipe{Cache: cache}
		*C = Ctx{Env: map[string]string{}}

		return fixtures.Cli(fixtures.Runner(), tests.TaskListCli{
			AppName:     "pipe-python",
			CommandName: "setup",
			TaskLists: []tests.TaskListFactory{
				func(p *Plumber, _ *cli.Command) *TaskList {
					tl := &TaskList{}

					return tl.New(p).
						SetRuntimeDepth(3).
						Set(func(tl *TaskList) Job {
							return JobSequence(env(tl).Job())
						})
				},
			},
		}).Run()
	}

	It("pins the cache of uv to an absolute directory", func() {
		Expect(run(".uv")).To(Succeed())

		cache, err := filepath.Abs(".uv")
		Expect(err).NotTo(HaveOccurred())
		Expect(C.Env).To(Equal(map[string]string{"UV_CACHE_DIR": cache}))
	})

	It("leaves the cache of the environment alone when none is given", func() {
		Expect(run("")).To(Succeed())

		Expect(C.Env).To(BeEmpty())
	})
})
