package tests

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manifest", func() {
	var manifest Manifest

	BeforeEach(func() {
		var err error

		manifest, err = ReadManifest()
		Expect(err).NotTo(HaveOccurred())
	})

	It("has one entry per pipe and nothing else", func() {
		names := make([]string, len(manifest.Pipes))
		for i, entry := range manifest.Pipes {
			names[i] = entry.Name
		}

		expected := make([]string, len(pipes))
		for i, p := range pipes {
			expected[i] = p.Dir
		}

		Expect(names).To(Equal(expected))
	})

	for _, p := range pipes {
		Describe(p.Dir, func() {
			var entry ManifestEntry

			BeforeEach(func() {
				for _, candidate := range manifest.Pipes {
					if candidate.Name == p.Dir {
						entry = candidate

						return
					}
				}

				Fail(fmt.Sprintf("%s has no entry in pipes.yaml", p.Dir))
			})

			It("publishes to the image the directory names", func() {
				Expect(entry.Image).To(Equal(p.Image()))
			})

			It("points at a readme that is there", func() {
				Expect(entry.Readme).To(Equal("./" + p.Dir + "/README.md"))
				Expect(filepath.Join(Root(), entry.Readme)).To(BeAnExistingFile())
			})

			It("carries a description", func() {
				Expect(entry.Description).NotTo(BeEmpty())
			})
		})
	}

	// The pipeline publishes the descriptions rather than reading them from the
	// manifest, so the two are checked against each other until one of them can
	// be generated from the other.
	It("agrees with the readme matrix the pipeline publishes", func() {
		matrix, err := ReadReadmeMatrix()
		Expect(err).NotTo(HaveOccurred())

		published := map[string]ReadmeMatrixEntry{}
		for _, entry := range matrix {
			published[entry.Repository] = entry
		}

		for _, entry := range manifest.Pipes {
			Expect(published).To(HaveKey(entry.Image), fmt.Sprintf("%s is not published by the update-docker-hub-readme job", entry.Name))
			Expect(published[entry.Image].File).To(Equal(entry.Readme))
			Expect(published[entry.Image].Description).To(Equal(entry.Description))
		}

		Expect(published).To(HaveLen(len(manifest.Pipes)), "the pipeline publishes an image the manifest does not list")
	})

	It("leaves the excluded modules out", func() {
		for _, dir := range Excluded {
			for _, entry := range manifest.Pipes {
				Expect(entry.Name).NotTo(Equal(dir))
			}

			_, err := os.Stat(filepath.Join(Root(), dir, "app", "app.go"))
			Expect(err).To(HaveOccurred(), fmt.Sprintf("%s builds a command tree, so it is a pipe and not an exclusion", dir))
		}
	})
})
