package markdown_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMarkdown(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Markdown Suite")
}

// the specs assert the colors a document reaches the log with, which the shell that runs
// them would otherwise decide; a spec that wants them gone asks for it on its own.
var _ = BeforeEach(func() {
	GinkgoT().Setenv("NO_COLOR", "")
	GinkgoT().Setenv("CLICOLOR_FORCE", "")
})
