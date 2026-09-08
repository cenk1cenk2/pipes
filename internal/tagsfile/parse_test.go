package tagsfile_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"gitlab.kilic.dev/devops/pipes/internal/tagsfile"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("Parse", func() {
	var (
		log *logrus.Entry
		dir string
	)

	BeforeEach(func() {
		log = fixtures.Log()
		dir = GinkgoT().TempDir()
	})

	write := func(name, content string) string {
		path := filepath.Join(dir, name)
		Expect(os.WriteFile(path, []byte(content), 0600)).To(Succeed())

		return path
	}

	It("reads the tags out of a comma separated file", func() {
		tags, err := tagsfile.Parse(log, write(".tags", "one,two,three"), false)
		Expect(err).NotTo(HaveOccurred())
		Expect(tags).To(Equal([]string{"one", "two", "three"}))
	})

	// the file is usually written by a shell redirect, which leaves the newline in.
	It("strips the line endings the writing job leaves behind", func() {
		tags, err := tagsfile.Parse(log, write(".tags", "one,\ntwo,\r\nthree\n"), false)
		Expect(err).NotTo(HaveOccurred())
		Expect(tags).To(Equal([]string{"one", "two", "three"}))
	})

	It("reads a single tag with no separator", func() {
		tags, err := tagsfile.Parse(log, write(".tags", "only"), false)
		Expect(err).NotTo(HaveOccurred())
		Expect(tags).To(Equal([]string{"only"}))
	})

	// the file is normally produced by an earlier job that may legitimately not have
	// run, so its absence is not on its own a reason to fail the pipeline.
	It("yields nothing and no error when the file is not there", func() {
		tags, err := tagsfile.Parse(log, filepath.Join(dir, "absent"), false)
		Expect(err).NotTo(HaveOccurred())
		Expect(tags).To(BeNil())
	})

	// strict is for the pipes that would otherwise publish something untagged, so
	// there the absence of a configured file is the failure itself.
	It("fails on an absent file under strict", func() {
		tags, err := tagsfile.Parse(log, filepath.Join(dir, "absent"), true)
		Expect(err).To(MatchError(ContainSubstring("Tags file is set but does not exists")))
		Expect(tags).To(BeNil())
	})

	It("yields nothing when no path was configured", func() {
		tags, err := tagsfile.Parse(log, "", false)
		Expect(err).NotTo(HaveOccurred())
		Expect(tags).To(BeNil())
	})

	// strict says the configured file has to be there, and no path is not a
	// configured file.
	It("yields nothing when no path was configured even under strict", func() {
		tags, err := tagsfile.Parse(log, "", true)
		Expect(err).NotTo(HaveOccurred())
		Expect(tags).To(BeNil())
	})

	// strict only ever guards the absence, so a file that is there parses the same
	// way whichever mode the pipe asked for.
	It("parses a present file the same way under strict", func() {
		tags, err := tagsfile.Parse(log, write(".tags", "one,two"), true)
		Expect(err).NotTo(HaveOccurred())
		Expect(tags).To(Equal([]string{"one", "two"}))
	})

	// a path that stats but does not read is a misconfigured pipeline, not a job
	// that has not run yet, so it fails whichever mode it was asked for.
	It("names the file it could not read", func() {
		for _, strict := range []bool{false, true} {
			tags, err := tagsfile.Parse(log, dir, strict)
			Expect(err).To(MatchError(ContainSubstring("Can not read the tags file")))
			Expect(err).To(MatchError(ContainSubstring(dir)))
			Expect(tags).To(BeNil())
		}
	})
})

var _ = Describe("Flags", func() {
	It("registers the strict flag alongside the path", func() {
		var (
			path   string
			strict bool
		)

		flags := tagsfile.NewFlags(tagsfile.Options{Destination: &path, Strict: &strict})

		Expect(flags).To(HaveLen(2))
		Expect(flags[0].Names()).To(Equal([]string{"tags-file"}))
		Expect(flags[1].Names()).To(Equal([]string{"tags-file.strict"}))
	})

	// pipes that always read the file leniently have no strict flag to document.
	It("leaves the strict flag out for a nil destination", func() {
		var path string

		flags := tagsfile.NewFlags(tagsfile.Options{Destination: &path, Value: ".tags"})

		Expect(flags).To(HaveLen(1))
		Expect(flags[0].Names()).To(Equal([]string{"tags-file"}))
	})
})
