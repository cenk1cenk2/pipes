package git_test

import (
	"io/fs"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/internal/git"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Sync", func() {
	var source, root string

	BeforeEach(func() {
		source = GinkgoT().TempDir()
		root = GinkgoT().TempDir()

		write(filepath.Join(root, "out", "stale.txt"), "stale")
		write(filepath.Join(root, "kept", "other.txt"), "other")
	})

	It("replaces the destination with the source, dropping what the source no longer has", func() {
		write(filepath.Join(source, "fresh.txt"), "fresh")
		write(filepath.Join(source, "nested", "deep.txt"), "deep")

		Expect(git.Sync(fixtures.Task(), root, []git.Path{{Source: source, Destination: "out"}}, false)).To(Succeed())

		Expect(filepath.Join(root, "out", "stale.txt")).NotTo(BeAnExistingFile())
		Expect(read(filepath.Join(root, "out", "fresh.txt"))).To(Equal("fresh"))
		Expect(read(filepath.Join(root, "out", "nested", "deep.txt"))).To(Equal("deep"))
		Expect(read(filepath.Join(root, "kept", "other.txt"))).To(Equal("other"))
	})

	It("leaves out what the excludes match, whole directories included", func() {
		write(filepath.Join(source, "fresh.txt"), "fresh")
		write(filepath.Join(source, "fresh.log"), "log")
		write(filepath.Join(source, "cache", "blob.txt"), "blob")

		paths := []git.Path{{Source: source, Destination: "out", Exclude: []string{"*.log", "cache"}}}
		Expect(git.Sync(fixtures.Task(), root, paths, false)).To(Succeed())

		Expect(filepath.Join(root, "out", "fresh.txt")).To(BeAnExistingFile())
		Expect(filepath.Join(root, "out", "fresh.log")).NotTo(BeAnExistingFile())
		Expect(filepath.Join(root, "out", "cache")).NotTo(BeAnExistingFile())
	})

	// a generate job that produced nothing would otherwise publish the destination
	// deleted, and nothing about the pipeline would look wrong.
	DescribeTable(
		"refuses a source that is missing or empty before touching any destination",
		func(prepare func() string) {
			write(filepath.Join(source, "fresh.txt"), "fresh")

			paths := []git.Path{
				{Source: source, Destination: "kept"},
				{Source: prepare(), Destination: "out"},
			}

			Expect(git.Sync(fixtures.Task(), root, paths, false)).NotTo(Succeed())
			Expect(read(filepath.Join(root, "kept", "other.txt"))).To(Equal("other"))
			Expect(read(filepath.Join(root, "out", "stale.txt"))).To(Equal("stale"))
		},
		Entry("a missing source", func() string {
			return filepath.Join(GinkgoT().TempDir(), "missing")
		}),
		Entry("an empty source", func() string {
			return GinkgoT().TempDir()
		}),
	)

	// an anchored pattern can match the relative path of the source itself, which would
	// skip the whole source and sync it as empty.
	It("never excludes the source itself", func() {
		write(filepath.Join(source, "fresh.txt"), "fresh")

		paths := []git.Path{{Source: source, Destination: "out", Exclude: []string{"?"}}}
		Expect(git.Sync(fixtures.Task(), root, paths, false)).To(Succeed())

		Expect(filepath.Join(root, "out", "fresh.txt")).To(BeAnExistingFile())
	})

	It("syncs a single file onto a file destination", func() {
		write(filepath.Join(source, "values.yaml"), "values")

		paths := []git.Path{{Source: filepath.Join(source, "values.yaml"), Destination: "out/values.yaml"}}
		Expect(git.Sync(fixtures.Task(), root, paths, false)).To(Succeed())

		Expect(read(filepath.Join(root, "out", "values.yaml"))).To(Equal("values"))
		Expect(read(filepath.Join(root, "out", "stale.txt"))).To(Equal("stale"))
	})

	DescribeTable(
		"refuses a source with nothing to sync before touching the destination",
		func(prepare func(), exclude []string) {
			prepare()

			paths := []git.Path{{Source: source, Destination: "out", Exclude: exclude}}
			Expect(git.Sync(fixtures.Task(), root, paths, false)).To(MatchError(ContainSubstring("empty")))
			Expect(read(filepath.Join(root, "out", "stale.txt"))).To(Equal("stale"))
		},
		Entry("only empty directories", func() {
			Expect(os.MkdirAll(filepath.Join(source, "a", "b"), 0o755)).To(Succeed())
		}, nil),
		Entry("only excluded entries", func() {
			write(filepath.Join(source, "fresh.log"), "log")
			write(filepath.Join(source, "cache", "blob.txt"), "blob")
		}, []string{"*.log", "cache"}),
	)

	It("refuses a broken exclude pattern before touching the destination", func() {
		write(filepath.Join(source, "fresh.txt"), "fresh")

		paths := []git.Path{{Source: source, Destination: "out", Exclude: []string{"["}}}
		Expect(git.Sync(fixtures.Task(), root, paths, true)).To(MatchError(ContainSubstring("exclude pattern")))
		Expect(read(filepath.Join(root, "out", "stale.txt"))).To(Equal("stale"))
	})

	// the outer destination is replaced after the inner one was written, or the other
	// way around, and either order loses one of them.
	It("refuses destinations that contain each other", func() {
		write(filepath.Join(source, "fresh.txt"), "fresh")

		paths := []git.Path{
			{Source: source, Destination: "out"},
			{Source: source, Destination: "out/sub"},
		}
		Expect(git.Sync(fixtures.Task(), root, paths, false)).NotTo(Succeed())
		Expect(read(filepath.Join(root, "out", "stale.txt"))).To(Equal("stale"))
	})

	It("refuses a destination behind a symlink leading out of the target", func() {
		outside := GinkgoT().TempDir()
		write(filepath.Join(outside, "out", "precious.txt"), "precious")
		Expect(os.Symlink(outside, filepath.Join(root, "link"))).To(Succeed())
		write(filepath.Join(source, "fresh.txt"), "fresh")

		Expect(git.Sync(fixtures.Task(), root, []git.Path{{Source: source, Destination: "link/out"}}, false)).NotTo(Succeed())
		Expect(read(filepath.Join(outside, "out", "precious.txt"))).To(Equal("precious"))
	})

	Describe("include", func() {
		files := func(root string) []string {
			GinkgoHelper()

			found := []string{}
			Expect(filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}

				relative, err := filepath.Rel(root, path)
				found = append(found, filepath.ToSlash(relative))

				return err
			})).To(Succeed())

			return found
		}

		BeforeEach(func() {
			write(filepath.Join(source, "README.md"), "readme")
			write(filepath.Join(source, "src", "main.go"), "main")
			write(filepath.Join(source, "src", "gen", "client.go"), "client")
			write(filepath.Join(source, "src", "gen", "client_mock.go"), "mock")
			write(filepath.Join(source, "src", "gen", "schema.json"), "schema")
		})

		DescribeTable(
			"syncs what the include and exclude patterns let through",
			func(include []string, exclude []string, expected []string) {
				paths := []git.Path{{Source: source, Destination: "out", Include: include, Exclude: exclude}}
				Expect(git.Sync(fixtures.Task(), root, paths, false)).To(Succeed())

				Expect(files(filepath.Join(root, "out"))).To(ConsistOf(expected))
			},
			Entry("everything without an include", nil, nil, []string{
				"README.md", "src/main.go", "src/gen/client.go", "src/gen/client_mock.go", "src/gen/schema.json",
			}),
			Entry("only the files an include matches", []string{"src/gen/*.go"}, nil, []string{
				"src/gen/client.go", "src/gen/client_mock.go",
			}),
			Entry("the whole tree of a directory an include matches", []string{"src/gen"}, nil, []string{
				"src/gen/client.go", "src/gen/client_mock.go", "src/gen/schema.json",
			}),
			Entry("nothing an exclude matches, even when an include does", []string{"src/gen/*.go"}, []string{"src/gen/*_mock.go"}, []string{
				"src/gen/client.go",
			}),
		)

		// a directory that leads to no include is never walked, which an unreadable one
		// makes visible: entering it would fail the sync.
		It("does not enter a directory no include can match beneath", func() {
			locked := filepath.Join(source, "docs")
			Expect(os.Mkdir(locked, 0o000)).To(Succeed())
			DeferCleanup(os.Chmod, locked, fs.FileMode(0o755))

			paths := []git.Path{{Source: source, Destination: "out", Include: []string{"src/gen/*.go"}}}
			Expect(git.Sync(fixtures.Task(), root, paths, false)).To(Succeed())
			Expect(files(filepath.Join(root, "out"))).To(ConsistOf("src/gen/client.go", "src/gen/client_mock.go"))
		})

		It("refuses a source with nothing the includes match before touching the destination", func() {
			paths := []git.Path{{Source: source, Destination: "out", Include: []string{"*.yaml"}}}

			Expect(git.Sync(fixtures.Task(), root, paths, false)).To(MatchError(ContainSubstring("empty")))
			Expect(read(filepath.Join(root, "out", "stale.txt"))).To(Equal("stale"))
		})

		It("refuses a broken include pattern before touching the destination", func() {
			paths := []git.Path{{Source: source, Destination: "out", Include: []string{"["}}}

			Expect(git.Sync(fixtures.Task(), root, paths, true)).To(MatchError(ContainSubstring("include pattern")))
			Expect(read(filepath.Join(root, "out", "stale.txt"))).To(Equal("stale"))
		})
	})

	It("clears the destination for a missing source when empty sources are allowed", func() {
		paths := []git.Path{{Source: filepath.Join(source, "missing"), Destination: "out"}}

		Expect(git.Sync(fixtures.Task(), root, paths, true)).To(Succeed())
		Expect(filepath.Join(root, "out")).NotTo(BeAnExistingFile())
	})

	DescribeTable(
		"refuses a destination outside the target",
		func(destination string) {
			write(filepath.Join(source, "fresh.txt"), "fresh")

			Expect(git.Sync(fixtures.Task(), root, []git.Path{{Source: source, Destination: destination}}, false)).
				NotTo(Succeed())
			Expect(filepath.Join(root, "out", "stale.txt")).To(BeAnExistingFile())
		},
		Entry("an absolute path", "/tmp/out"),
		Entry("a path climbing out", "../out"),
		Entry("the target itself", "."),
		Entry("the git directory of the target", ".git"),
		Entry("the hooks of the target", ".GIT/hooks"),
	)
})

func write(path string, content string) {
	GinkgoHelper()

	Expect(os.MkdirAll(filepath.Dir(path), 0o755)).To(Succeed())
	Expect(os.WriteFile(path, []byte(content), 0o644)).To(Succeed())
}

func read(path string) string {
	GinkgoHelper()

	content, err := os.ReadFile(path)
	Expect(err).NotTo(HaveOccurred())

	return string(content)
}
