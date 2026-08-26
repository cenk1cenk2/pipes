package git_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/internal/git"
)

var _ = Describe("References", func() {
	It("prefixes each ref with the kind the match patterns are written against", func() {
		Expect(git.Refs{Branch: "main", Tag: "v1.0.0"}.References()).
			To(Equal([]string{"tags/v1.0.0", "heads/main"}))
	})

	// A tagged pipeline carries the branch it was tagged on as well, so the order
	// is what decides that a release matches the tag rule and not the branch rule.
	It("puts the tag ahead of the branch", func() {
		Expect(git.Refs{Branch: "main", Tag: "v1.0.0"}.References()[0]).To(Equal("tags/v1.0.0"))
	})

	It("leaves out whichever side is unset", func() {
		Expect(git.Refs{Branch: "main"}.References()).To(Equal([]string{"heads/main"}))
		Expect(git.Refs{Tag: "v1.0.0"}.References()).To(Equal([]string{"tags/v1.0.0"}))
	})

	It("gives back an empty list when the pipe knows neither", func() {
		Expect(git.Refs{}.References()).To(BeEmpty())
	})
})

var _ = Describe("MatchAny", func() {
	references := []string{"tags/v1.0.0", "heads/main"}

	It("reports the first matching pattern", func() {
		matched, err := git.MatchAny([]string{`^heads/`, `^tags/`}, references)
		Expect(err).NotTo(HaveOccurred())
		Expect(matched).To(Equal(0))
	})

	// The pattern order is the user's rule order, so an earlier rule wins even when
	// a later one matches a reference that sorts first.
	It("ranks the patterns above the references", func() {
		matched, err := git.MatchAny([]string{`^heads/main$`, `^tags/v1\.0\.0$`}, references)
		Expect(err).NotTo(HaveOccurred())
		Expect(matched).To(Equal(0))
	})

	It("reports -1 when nothing matches", func() {
		matched, err := git.MatchAny([]string{`^heads/release/`}, references)
		Expect(err).NotTo(HaveOccurred())
		Expect(matched).To(Equal(-1))
	})

	It("reports -1 for no patterns and for no references", func() {
		matched, err := git.MatchAny(nil, references)
		Expect(err).NotTo(HaveOccurred())
		Expect(matched).To(Equal(-1))

		matched, err = git.MatchAny([]string{`^tags/`}, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(matched).To(Equal(-1))
	})

	It("names the pattern it could not compile", func() {
		_, err := git.MatchAny([]string{`^tags/(`}, references)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(`^tags/(`))
	})

	// A bad pattern behind a matching one would otherwise decide whether the pipe
	// fails purely by where the user put it in the list.
	It("stops at the first match before reaching a broken pattern", func() {
		matched, err := git.MatchAny([]string{`^tags/`, `^heads/(`}, references)
		Expect(err).NotTo(HaveOccurred())
		Expect(matched).To(Equal(0))
	})
})

var _ = Describe("Flags", func() {
	It("binds both destinations onto the same refs", func() {
		refs := git.Refs{}
		flags := git.NewFlags(&refs)

		Expect(flags).To(HaveLen(2))
		Expect(flags[0].Names()).To(Equal([]string{"git.branch"}))
		Expect(flags[1].Names()).To(Equal([]string{"git.tag"}))
	})
})
