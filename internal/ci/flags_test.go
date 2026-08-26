package ci_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/ci"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
	"gitlab.kilic.dev/devops/pipes/internal/report/iac"
)

var _ = Describe("NewMetadataFlags", func() {
	names := func(flags []ucli.Flag) []string {
		found := []string{}
		for _, flag := range flags {
			found = append(found, flag.Names()...)
		}

		return found
	}

	It("registers every coordinate the report renders", func() {
		metadata := iac.Metadata{}

		Expect(names(ci.NewMetadataFlags(&metadata))).To(Equal([]string{
			"report.job-name",
			"report.job-url",
			"report.pipeline-id",
			"report.pipeline-url",
			"report.commit-sha",
			"report.commit-short-sha",
		}))
	})

	It("files them all under one category", func() {
		metadata := iac.Metadata{}

		for _, flag := range ci.NewMetadataFlags(&metadata) {
			Expect(flag.(ucli.CategorizableFlag).GetCategory()).To(Equal(cli.CATEGORY_CI))
		}
	})

	// The flags exist to fill the report metadata, so each one has to land on its
	// own field of the struct the caller passed rather than a copy of it.
	It("binds each flag onto the given metadata", func() {
		metadata := iac.Metadata{}
		flags := ci.NewMetadataFlags(&metadata)

		//nolint:errcheck
		Expect(flags[0].(*ucli.StringFlag).Destination).To(BeIdenticalTo(&metadata.JobName))
		//nolint:errcheck
		Expect(flags[5].(*ucli.StringFlag).Destination).To(BeIdenticalTo(&metadata.CommitShortSha))
	})
})
