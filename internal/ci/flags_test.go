package ci_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/ci"
	"gitlab.kilic.dev/devops/pipes/internal/report/iac"
)

var _ = Describe("NewFlags", func() {
	names := func(flags []cli.Flag) []string {
		found := []string{}
		for _, flag := range flags {
			found = append(found, flag.Names()...)
		}

		return found
	}

	It("registers every coordinate the report renders", func() {
		metadata := iac.Metadata{}

		Expect(names(ci.NewFlags(ci.Options{Destination: &metadata}))).To(Equal([]string{
			"ci.job-name",
			"ci.job-url",
			"ci.pipeline-id",
			"ci.pipeline-url",
			"ci.commit-sha",
			"ci.commit-short-sha",
		}))
	})

	It("files them all under one category", func() {
		metadata := iac.Metadata{}

		for _, flag := range ci.NewFlags(ci.Options{Destination: &metadata}) {
			Expect(flag.(cli.CategorizableFlag).GetCategory()).To(Equal(ci.CATEGORY_CI))
		}
	})

	// The flags exist to fill the report metadata, so each one has to land on its
	// own field of the struct the caller passed rather than a copy of it.
	It("binds each flag onto the given metadata", func() {
		metadata := iac.Metadata{}
		flags := ci.NewFlags(ci.Options{Destination: &metadata})

		//nolint:errcheck
		Expect(flags[0].(*cli.StringFlag).Destination).To(BeIdenticalTo(&metadata.JobName))
		//nolint:errcheck
		Expect(flags[5].(*cli.StringFlag).Destination).To(BeIdenticalTo(&metadata.CommitShortSha))
	})
})
