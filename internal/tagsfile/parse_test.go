package tagsfile_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"gitlab.kilic.dev/devops/pipes/internal/tagsfile"
)

var _ = Describe("Parse", func() {
	var (
		log *logrus.Entry
		dir string
	)

	BeforeEach(func() {
		logger := logrus.New()
		logger.SetOutput(GinkgoWriter)
		log = logrus.NewEntry(logger)
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

	// The file is usually written by a shell redirect, which leaves the newline in.
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

	// The file is normally produced by an earlier job that may legitimately not have
	// run, so its absence is not on its own a reason to fail the pipeline.
	It("yields nothing and no error when the file is not there", func() {
		tags, err := tagsfile.Parse(log, filepath.Join(dir, "absent"), false)
		Expect(err).NotTo(HaveOccurred())
		Expect(tags).To(BeNil())
	})

	It("yields nothing for an absent file even under strict", func() {
		tags, err := tagsfile.Parse(log, filepath.Join(dir, "absent"), true)
		Expect(err).NotTo(HaveOccurred())
		Expect(tags).To(BeNil())
	})

	It("yields nothing when no path was configured", func() {
		tags, err := tagsfile.Parse(log, "", false)
		Expect(err).NotTo(HaveOccurred())
		Expect(tags).To(BeNil())
	})
})

var _ = Describe("Flags", func() {
	It("registers the strict flag alongside the path", func() {
		var (
			path   string
			strict bool
		)

		flags := tagsfile.NewFlags(&path, "", &strict, false)

		Expect(flags).To(HaveLen(2))
		Expect(flags[0].Names()).To(Equal([]string{"tags-file"}))
		Expect(flags[1].Names()).To(Equal([]string{"tags-file.strict"}))
	})

	// Pipes that always read the file leniently have no strict flag to document.
	It("leaves the strict flag out for a nil destination", func() {
		var path string

		flags := tagsfile.NewFlags(&path, ".tags", nil, false)

		Expect(flags).To(HaveLen(1))
		Expect(flags[0].Names()).To(Equal([]string{"tags-file"}))
	})
})
