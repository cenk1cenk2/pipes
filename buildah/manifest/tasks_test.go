package manifest

import (
	. "github.com/cenk1cenk2/plumber/v7"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"
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
})
