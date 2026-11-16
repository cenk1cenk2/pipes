package patch_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/internal/report/patch"
)

var _ = Describe("Filename", func() {
	It("names a renamed file with spaces by where it ends up", func() {
		Expect(patch.Filename("diff --git a/old name.md b/new name.md\nsimilarity index 100%\n")).To(Equal("new name.md"))
	})
})

var _ = Describe("Previous", func() {
	It("names a renamed file with spaces by where it comes from", func() {
		Expect(patch.Previous("diff --git a/old name.md b/new name.md\nsimilarity index 100%\n")).To(Equal("old name.md"))
	})
})

var _ = Describe("Fence", func() {
	It("stays at three backticks for a diff without a fence of its own", func() {
		Expect(patch.Fence("@@ -1 +1 @@\n-a\n+b\n")).To(Equal("```"))
	})

	It("grows past the longest backtick run the diff carries", func() {
		Expect(patch.Fence("@@ -1,3 +1,3 @@\n+```go\n+````\n")).To(Equal("`````"))
	})
})
