package tests

import (
	"fmt"
	"slices"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pipes", func() {
	It("registers every module that builds a command tree", func() {
		dirs, err := ModuleDirs()
		Expect(err).NotTo(HaveOccurred())

		registered := make([]string, len(pipes))
		for i, p := range pipes {
			registered[i] = p.Dir
		}

		unregistered := []string{}
		for _, dir := range dirs {
			if slices.Contains(Excluded, dir) || slices.Contains(registered, dir) {
				continue
			}

			unregistered = append(unregistered, dir)
		}

		Expect(unregistered).To(BeEmpty(), "add the module to the conformance table, or to Excluded with the reason it is not a pipe")
	})

	It("excludes only directories that are there", func() {
		dirs, err := ModuleDirs()
		Expect(err).NotTo(HaveOccurred())

		for _, dir := range Excluded {
			Expect(dirs).To(ContainElement(dir), "the exclusion outlived the directory it was written for")
		}
	})

	for _, p := range pipes {
		Describe(p.Dir, func() {
			It("names the command after the directory", func() {
				// select-env is the one pipe whose command is not prefixed, since
				// the pipelines that call it were written before the convention.
				expected := "pipe-" + p.Dir
				if p.Dir == "select-env" {
					expected = "select-env"
				}

				Expect(p.Name).To(Equal(expected))
			})

			It("describes itself", func() {
				Expect(p.Description).NotTo(BeEmpty())
				Expect(p.New(nil).Usage).To(Equal(p.Description))
			})

			It("gives every visible flag an environment source", func() {
				for _, f := range p.visibleFlags() {
					Expect(f.EnvVars()).NotTo(BeEmpty(), fmt.Sprintf("%s: %s has no environment source", f.Command, f.Name()))
				}
			})

			It("gives every visible flag a category", func() {
				allowed := uncategorizedFlags[p.Dir]

				var seen []string

				for _, f := range p.visibleFlags() {
					if f.Category() != "" {
						Expect(allowed).NotTo(ContainElement(f.Name()), fmt.Sprintf("%s now has a category, drop it from uncategorizedFlags", f.Name()))

						continue
					}

					Expect(allowed).To(ContainElement(f.Name()), fmt.Sprintf("%s: %s has no category", f.Command, f.Name()))

					if !slices.Contains(seen, f.Name()) {
						seen = append(seen, f.Name())
					}
				}

				slices.Sort(seen)
				Expect(seen).To(Equal(allowed), "uncategorizedFlags lists a flag the tree no longer carries")
			})

			It("keeps the legacy environment names ahead of the canonical one", func() {
				aliased := legacyEnvAliases[p.Dir]

				var seen []string

				for _, f := range p.visibleFlags() {
					envs := f.EnvVars()

					if len(envs) < 2 {
						Expect(aliased).NotTo(HaveKey(f.Name()), fmt.Sprintf("%s lost its legacy environment names", f.Name()))

						continue
					}

					Expect(aliased).To(
						HaveKeyWithValue(f.Name(), envs),
						fmt.Sprintf("%s: %s answers to more than one name, record the chain in legacyEnvAliases", f.Command, f.Name()),
					)

					if !slices.Contains(seen, f.Name()) {
						seen = append(seen, f.Name())
					}
				}

				var recorded []string
				for name := range aliased {
					recorded = append(recorded, name)
				}

				slices.Sort(recorded)
				slices.Sort(seen)

				Expect(seen).To(Equal(recorded), "legacyEnvAliases lists a flag the tree no longer carries")
			})
		})
	}
})
