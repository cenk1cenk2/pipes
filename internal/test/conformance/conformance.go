// Package conformance is the contract every pipe answers to. A pipe hands Verify
// the tree it builds and gets the whole contract back as specs of its own suite,
// so a convention that only lived in a review comment fails a build instead.
//
// The specs live here rather than in a module that imports every pipe, which is
// what lets a pipe keep its command tree to itself.
package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cenk1cenk2/plumber/v6"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"
)

// Pipe is everything a pipe has to state about itself for the contract to be
// checked: the tree it builds, and the tables recording what it is still allowed
// to be missing.
type Pipe struct {
	Name        string
	Description string
	// New builds the tree. A nil plumber is enough for the shape, since a step
	// only reaches for one once the command runs.
	New func(p *plumber.Plumber) *ucli.Command
	// Unprefixed drops the pipe- prefix the command name is otherwise held to.
	// Only select-env sets it, since the pipelines that call it were written
	// before the convention.
	Unprefixed bool
	// UncategorizedFlags are the visible flags that predate the categories the
	// generated documentation files flags under. The list only ever shrinks: a new
	// flag without a category fails, and so does an entry that has since been
	// given one.
	UncategorizedFlags []string
	// LegacyEnvAliases is the exact, ordered environment source chain of every
	// visible flag that answers to more than one name. The names a pipeline
	// already sets are kept forever and listed ahead of the canonical one, so the
	// order is the precedence and dropping or reordering an entry silently changes
	// which value a running pipeline picks up.
	//
	// The table is closed in both directions: a flag listed here has to keep this
	// chain, and a flag that grows a second source has to be added here.
	LegacyEnvAliases map[string][]string
}

// Verify registers the conformance specs of the pipe. It answers with a bool so
// it can be assigned at package level the way a Ginkgo container is.
func Verify(pipe Pipe) bool {
	// The directory is resolved while the tree is built rather than while a spec
	// runs, since the behaviour specs of a pipe move into a temporary working
	// directory and would otherwise change what the contract reads.
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dir := filepath.Base(wd)
	readmePath := filepath.Join(wd, "README.md")

	return Describe("Conformance", func() {
		It("names the command after the directory", func() {
			expected := dir
			if !pipe.Unprefixed {
				expected = "pipe-" + dir
			}

			Expect(pipe.Name).To(Equal(expected))
		})

		It("describes itself", func() {
			Expect(pipe.Description).NotTo(BeEmpty())
			Expect(pipe.New(nil).Usage).To(Equal(pipe.Description))
		})

		It("gives every visible flag an environment source", func() {
			for _, f := range pipe.visibleFlags() {
				Expect(f.EnvVars()).NotTo(BeEmpty(), fmt.Sprintf("%s: %s has no environment source", f.Command, f.Name()))
			}
		})

		It("gives every visible flag a category", func() {
			allowed := pipe.UncategorizedFlags

			var seen []string

			for _, f := range pipe.visibleFlags() {
				if f.Category() != "" {
					Expect(allowed).NotTo(ContainElement(f.Name()), fmt.Sprintf("%s now has a category, drop it from UncategorizedFlags", f.Name()))

					continue
				}

				Expect(allowed).To(ContainElement(f.Name()), fmt.Sprintf("%s: %s has no category", f.Command, f.Name()))

				if !slices.Contains(seen, f.Name()) {
					seen = append(seen, f.Name())
				}
			}

			slices.Sort(seen)
			Expect(seen).To(Equal(allowed), "UncategorizedFlags lists a flag the tree no longer carries")
		})

		It("keeps the legacy environment names ahead of the canonical one", func() {
			aliased := pipe.LegacyEnvAliases

			var seen []string

			for _, f := range pipe.visibleFlags() {
				envs := f.EnvVars()

				if len(envs) < 2 {
					Expect(aliased).NotTo(HaveKey(f.Name()), fmt.Sprintf("%s lost its legacy environment names", f.Name()))

					continue
				}

				Expect(aliased).To(
					HaveKeyWithValue(f.Name(), envs),
					fmt.Sprintf("%s: %s answers to more than one name, record the chain in LegacyEnvAliases", f.Command, f.Name()),
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

			Expect(seen).To(Equal(recorded), "LegacyEnvAliases lists a flag the tree no longer carries")
		})

		// The README is generated from the command tree, so it goes stale the moment
		// a flag is added or renamed without regenerating it. The specs read the
		// checked in file rather than running the generator, since what is published
		// is the file.
		Describe("README", func() {
			var readme string

			BeforeEach(func() {
				contents, err := os.ReadFile(readmePath)
				Expect(err).NotTo(HaveOccurred())

				readme = string(contents)
			})

			It("opens with the name and the description of the pipe", func() {
				lines := strings.Split(readme, "\n")
				Expect(len(lines)).To(BeNumerically(">=", 3))
				Expect(lines[0]).To(Equal("# " + pipe.Name))
				Expect(lines[2]).To(Equal(pipe.Description))
			})

			It("documents every environment source of every visible flag", func() {
				for _, f := range pipe.visibleFlags() {
					for _, env := range f.EnvVars() {
						Expect(readme).To(
							ContainSubstring("$"+env+"`"),
							fmt.Sprintf("%s: %s reads $%s, regenerate the documentation with: task %s:docs", f.Command, f.Name(), env, dir),
						)
					}
				}
			})

			It("documents every subcommand", func() {
				for _, sub := range pipe.New(nil).Commands {
					Expect(readme).To(ContainSubstring(fmt.Sprintf("### `%s %s`", pipe.Name, sub.Name)))
				}
			})
		})
	})
}

// flagRef is a flag together with the command path it was reached through, so a
// failure names the subcommand a reader has to open.
type flagRef struct {
	Command string
	Flag    ucli.Flag
}

func (f flagRef) Name() string {
	return f.Flag.Names()[0]
}

func (f flagRef) Category() string {
	if c, ok := f.Flag.(ucli.CategorizableFlag); ok {
		return c.GetCategory()
	}

	return ""
}

func (f flagRef) EnvVars() []string {
	if d, ok := f.Flag.(ucli.DocGenerationFlag); ok {
		return d.GetEnvVars()
	}

	return nil
}

func (f flagRef) Visible() bool {
	if v, ok := f.Flag.(ucli.VisibleFlag); ok {
		return v.IsVisible()
	}

	return true
}

// visibleFlags walks the whole tree and returns every flag a user could be shown.
// A flag shared between subcommands is the same instance, so it comes back once
// per command it is reachable from.
func (p Pipe) visibleFlags() []flagRef {
	refs := []flagRef{}

	var walk func(path string, c *ucli.Command)

	walk = func(path string, c *ucli.Command) {
		for _, f := range c.Flags {
			ref := flagRef{Command: path, Flag: f}

			if ref.Visible() {
				refs = append(refs, ref)
			}
		}

		for _, sub := range c.Commands {
			walk(path+" "+sub.Name, sub)
		}
	}

	walk(p.Name, p.New(nil))

	return refs
}
