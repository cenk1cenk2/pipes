package preview

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/pulumi/setup"
	"gitlab.kilic.dev/devops/pipes/pulumi/stack"
)

var _ = Describe("Pulumi report discriminators", func() {
	BeforeEach(func() {
		stack.P.Stack = "main"
		setup.C.Cwd = "."
	})

	It("carries the stack on its own from the working directory", func() {
		Expect(pulumiReportDiscriminators()).To(Equal([]string{"main"}))
	})

	It("carries a project directory alongside the stack", func() {
		setup.C.Cwd = "projects/warehouse"

		Expect(pulumiReportDiscriminators()).To(Equal([]string{"main", "projects/warehouse"}))
	})
})
