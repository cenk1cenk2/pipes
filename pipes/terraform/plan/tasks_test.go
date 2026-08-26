package plan

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.kilic.dev/devops/pipes/pipes/terraform/setup"
	"gitlab.kilic.dev/devops/pipes/pipes/terraform/state"
)

var _ = Describe("Terraform report discriminators", func() {
	BeforeEach(func() {
		state.P.State.Name = "default"
		setup.C.Cwd = "."
	})

	It("holds nothing back when neither value has moved off its default", func() {
		Expect(terraformReportDiscriminators()).To(BeEmpty())
	})

	It("carries a named state", func() {
		state.P.State.Name = "production"

		Expect(terraformReportDiscriminators()).To(Equal([]string{"production"}))
	})

	It("carries a root module other than the working directory", func() {
		setup.C.Cwd = ".deploy/sun"

		Expect(terraformReportDiscriminators()).To(Equal([]string{".deploy/sun"}))
	})

	It("carries both, state first", func() {
		state.P.State.Name = "production"
		setup.C.Cwd = ".deploy/sun"

		Expect(terraformReportDiscriminators()).To(Equal([]string{"production", ".deploy/sun"}))
	})

	It("reports the state name only once it has been named", func() {
		Expect(terraformStateName()).To(BeEmpty())

		state.P.State.Name = "production"
		Expect(terraformStateName()).To(Equal("production"))
	})
})
