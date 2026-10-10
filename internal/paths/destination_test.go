package paths_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/internal/paths"
)

var _ = Describe("Destination", func() {
	var root string

	BeforeEach(func() {
		root = GinkgoT().TempDir()
		Expect(os.MkdirAll(filepath.Join(root, "out"), 0o755)).To(Succeed())
	})

	DescribeTable(
		"accepts a path inside the target",
		func(path string) {
			Expect(paths.Destination(root, path)).To(Succeed())
		},
		Entry("an existing directory", "out"),
		Entry("a file in an existing directory", "out/values.yaml"),
		Entry("a path whose parents do not exist yet", "new/deeply/nested"),
		Entry("a name that only starts like .git", ".github/workflows"),
	)

	DescribeTable(
		"refuses a path outside the target",
		func(path string) {
			Expect(paths.Destination(root, path)).NotTo(Succeed())
		},
		Entry("an absolute path", "/tmp/out"),
		Entry("a path climbing out", "../out"),
		Entry("the target itself", "."),
		Entry("an empty path", ""),
	)

	// replacing the git directory wipes the repository, and anything written into its
	// hooks runs on the next git call.
	DescribeTable(
		"refuses a path inside a .git directory",
		func(path string) {
			Expect(paths.Destination(root, path)).To(MatchError(ContainSubstring(".git")))
		},
		Entry("the git directory", ".git"),
		Entry("its hooks", ".git/hooks"),
		Entry("a differently cased one", ".GIT/hooks"),
		Entry("a nested one", "out/.git"),
	)

	It("refuses a path whose parent is a symlink leading outside the target", func() {
		outside := GinkgoT().TempDir()
		Expect(os.Symlink(outside, filepath.Join(root, "link"))).To(Succeed())

		Expect(paths.Destination(root, "link/out")).To(MatchError(ContainSubstring("outside")))
	})

	It("accepts a path whose parent is a symlink staying inside the target", func() {
		Expect(os.Symlink(filepath.Join(root, "out"), filepath.Join(root, "link"))).To(Succeed())

		Expect(paths.Destination(root, "link/sub")).To(Succeed())
	})
})

var _ = Describe("Disjoint", func() {
	It("accepts destinations that do not contain each other", func() {
		Expect(paths.Disjoint([]string{"out", "output", "other/out"})).To(Succeed())
	})

	DescribeTable(
		"refuses destinations where one contains another",
		func(destinations []string) {
			Expect(paths.Disjoint(destinations)).NotTo(Succeed())
		},
		Entry("a nested one after its parent", []string{"out", "out/sub"}),
		Entry("a nested one before its parent", []string{"out/sub", "out"}),
		Entry("the same one twice", []string{"out", "./out/"}),
	)
})
