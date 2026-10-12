package cobertura_test

import (
	"encoding/xml"
	"os"
	"path/filepath"

	"github.com/aereal/gobertura"
	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/cenk1cenk2/plumber/v7/tests"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/internal/report/cobertura"
)

var _ = Describe("GoProfileConverter", func() {
	var dir string

	BeforeEach(func() {
		dir = GinkgoT().TempDir()
	})

	write := func(name, content string) string {
		path := filepath.Join(dir, name)
		Expect(os.WriteFile(path, []byte(content), 0o600)).To(Succeed())

		return path
	}

	read := func(path string) gobertura.Coverage {
		GinkgoHelper()

		body, err := os.ReadFile(path)
		Expect(err).NotTo(HaveOccurred())

		var coverage gobertura.Coverage
		Expect(xml.Unmarshal(body, &coverage)).To(Succeed())

		return coverage
	}

	// GitLab joins the source with the file name and matches that against the repository.
	It("sources the report from the module and names its files relative to it", func() {
		write("go.mod", "module example.com/mod\n")
		input := write("coverage.out", "mode: set\nexample.com/mod/pkg/a.go:1.1,2.2 1 1\nexample.com/other/b.go:1.1,2.2 1 1\n")
		output := filepath.Join(dir, "cobertura.xml")

		Expect(cobertura.GoProfileConverter{}.Convert(input, output, dir)).To(Succeed())

		coverage := read(output)
		Expect(coverage.Sources.Sources).To(Equal([]string{dir}))

		filenames := []string{}
		for _, pkg := range coverage.Packages {
			for _, class := range pkg.Classes {
				filenames = append(filenames, class.Filename)
			}
		}
		Expect(filenames).To(ConsistOf("pkg/a.go", "example.com/other/b.go"))
	})

	// a suite under the root module, such as test/e2e, has no go.mod of its own and its
	// package may carry no statements at all.
	It("sources the report from the module the directory is under", func() {
		write("go.mod", "module example.com/mod\n")
		pkg := filepath.Join(dir, "test", "e2e")
		Expect(os.MkdirAll(pkg, 0o755)).To(Succeed())
		input := filepath.Join(pkg, "coverage.out")
		Expect(os.WriteFile(input, []byte("mode: atomic\n"), 0o600)).To(Succeed())
		output := filepath.Join(pkg, "cobertura.xml")

		Expect(cobertura.GoProfileConverter{}.Convert(input, output, pkg)).To(Succeed())

		coverage := read(output)
		Expect(coverage.Sources.Sources).To(Equal([]string{dir}))
		Expect(coverage.Packages).To(BeEmpty())
	})

	It("fails when the directory is not a module", func() {
		input := write("coverage.out", "mode: set\nexample.com/mod/a.go:1.1,2.2 1 1\n")

		Expect(cobertura.GoProfileConverter{}.Convert(input, filepath.Join(dir, "cobertura.xml"), dir)).
			To(MatchError(ContainSubstring("read module of the coverage profile")))
	})
})

var _ = Describe("Source", func() {
	var modules []string

	BeforeEach(func() {
		root := GinkgoT().TempDir()
		modules = []string{filepath.Join(root, "a"), filepath.Join(root, "b")}

		for _, module := range modules {
			Expect(os.MkdirAll(module, 0o755)).To(Succeed())
			Expect(os.WriteFile(filepath.Join(module, "go.mod"), []byte("module example.com/"+filepath.Base(module)+"\n"), 0o600)).To(Succeed())
		}
	})

	convert := func(src cobertura.Source) error {
		GinkgoHelper()

		return src.Convert(NewTaskList(tests.NewPlumber().Plumber).CreateTask("cobertura"))
	}

	It("converts the profile of every directory the tests ran in", func() {
		for _, module := range modules {
			Expect(os.WriteFile(filepath.Join(module, "coverage.out"), []byte("mode: set\nexample.com/"+filepath.Base(module)+"/x.go:1.1,2.2 1 1\n"), 0o600)).To(Succeed())
		}

		Expect(
			convert(cobertura.Source{Adapter: cobertura.GoProfileConverter{}, Input: "coverage.out", Output: "cobertura.xml", Dirs: modules}),
		).To(Succeed())

		for _, module := range modules {
			Expect(filepath.Join(module, "cobertura.xml")).To(BeAnExistingFile())
		}
	})

	// fail-fast stops a workspace at the first failing module, so the rest never write a profile.
	It("skips a directory the tests left no profile in", func() {
		Expect(os.WriteFile(filepath.Join(modules[0], "coverage.out"), []byte("mode: set\nexample.com/a/x.go:1.1,2.2 1 1\n"), 0o600)).To(Succeed())

		Expect(
			convert(cobertura.Source{Adapter: cobertura.GoProfileConverter{}, Input: "coverage.out", Output: "cobertura.xml", Dirs: modules}),
		).To(Succeed())

		Expect(filepath.Join(modules[0], "cobertura.xml")).To(BeAnExistingFile())
		Expect(filepath.Join(modules[1], "cobertura.xml")).NotTo(BeAnExistingFile())
	})

	It("refuses a path every directory would share", func() {
		for _, src := range []cobertura.Source{
			{Input: "coverage.out", Output: filepath.Join(GinkgoT().TempDir(), "cobertura.xml")},
			{Input: "../coverage.out", Output: "cobertura.xml"},
		} {
			src.Adapter = cobertura.GoProfileConverter{}
			src.Dirs = modules

			Expect(convert(src)).To(MatchError(ContainSubstring("otherwise the directories overwrite each other's report")))
		}
	})

	It("hands every directory's own paths to the adapter it was given", func() {
		for _, module := range modules {
			Expect(os.WriteFile(filepath.Join(module, "coverage.out"), []byte("mode: set\n"), 0o600)).To(Succeed())
		}
		adapter := &recordingAdapter{}

		Expect(convert(cobertura.Source{Adapter: adapter, Input: "coverage.out", Output: "cobertura.xml", Dirs: modules})).To(Succeed())

		Expect(adapter.calls).To(Equal([][3]string{
			{filepath.Join(modules[0], "coverage.out"), filepath.Join(modules[0], "cobertura.xml"), modules[0]},
			{filepath.Join(modules[1], "coverage.out"), filepath.Join(modules[1], "cobertura.xml"), modules[1]},
		}))
	})
})

// records what the task asks of it instead of converting anything.
type recordingAdapter struct {
	calls [][3]string
}

func (a *recordingAdapter) Convert(input string, output string, dir string) error {
	a.calls = append(a.calls, [3]string{input, output, dir})

	return nil
}
