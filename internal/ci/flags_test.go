package ci_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/ci"
	"gitlab.kilic.dev/devops/pipes/internal/report/terraform"
	"gitlab.kilic.dev/devops/pipes/tests/fixtures"
)

var _ = Describe("NewFlags", func() {
	It("files them all under one category", func() {
		metadata := terraform.Metadata{}

		for _, flag := range ci.NewFlags(ci.Options{Destination: &metadata}) {
			Expect(flag.(cli.CategorizableFlag).GetCategory()).To(Equal(ci.CategoryCI))
		}
	})

	// the flags exist to fill the report metadata, so an unbound one leaves the
	// report with a hole in it that nothing else fills.
	It("binds every flag onto the given metadata", func() {
		metadata := terraform.Metadata{}

		for _, flag := range ci.NewFlags(ci.Options{Destination: &metadata}) {
			//nolint:errcheck
			Expect(flag.(*cli.StringFlag).Destination).NotTo(BeNil(), flag.Names()[0])
		}
	})

	// two pipes each register their own metadata, so one call handing back the flags
	// of another would bind both onto whichever ran last.
	It("builds a fresh set of flags per metadata", func() {
		first, second := terraform.Metadata{}, terraform.Metadata{}

		Expect(fixtures.Flag[*cli.StringFlag](ci.NewFlags(ci.Options{Destination: &first}), "ci.job-name").Destination).
			To(BeIdenticalTo(&first.JobName))
		Expect(fixtures.Flag[*cli.StringFlag](ci.NewFlags(ci.Options{Destination: &second}), "ci.job-name").Destination).
			To(BeIdenticalTo(&second.JobName))
	})
})
