package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// The READMEs are generated from the command tree, so they go stale the moment a
// flag is added or renamed without regenerating them. The specs read the checked
// in file rather than running the generator, since what is published is the file.
var _ = Describe("README", func() {
	for _, p := range pipes {
		Describe(p.Dir, func() {
			var readme string

			BeforeEach(func() {
				contents, err := os.ReadFile(filepath.Join(Root(), p.Dir, "README.md"))
				Expect(err).NotTo(HaveOccurred())

				readme = string(contents)
			})

			It("opens with the name and the description of the pipe", func() {
				lines := strings.Split(readme, "\n")
				Expect(len(lines)).To(BeNumerically(">=", 3))
				Expect(lines[0]).To(Equal("# " + p.Name))
				Expect(lines[2]).To(Equal(p.Description))
			})

			It("documents every environment source of every visible flag", func() {
				for _, f := range p.visibleFlags() {
					for _, env := range f.EnvVars() {
						Expect(readme).To(
							ContainSubstring("$"+env+"`"),
							fmt.Sprintf("%s: %s reads $%s, regenerate the documentation with: task %s:docs", f.Command, f.Name(), env, p.Dir),
						)
					}
				}
			})

			It("documents every subcommand", func() {
				for _, sub := range p.New(nil).Commands {
					Expect(readme).To(ContainSubstring(fmt.Sprintf("### `%s %s`", p.Name, sub.Name)))
				}
			})
		})
	}
})
