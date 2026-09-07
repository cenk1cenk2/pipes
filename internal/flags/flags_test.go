package flags_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

type condition struct {
	Match       string `json:"match"       yaml:"match"`
	Environment string `json:"environment" yaml:"environment"`
}

var _ = Describe("JSONFlag", func() {
	It("unmarshals the value into the destination", func() {
		dst := []condition{}
		flag := flags.JSONFlag(&cli.StringFlag{Name: "conditions"}, &dst)

		Expect(flag.Validator(`[{ "match": "^heads/main$", "environment": "develop" }]`)).To(Succeed())
		Expect(dst).To(Equal([]condition{{Match: "^heads/main$", Environment: "develop"}}))
	})

	// Most of these flags are optional, so an unset one leaves the pipe on its zero
	// value rather than failing before the pipe has a chance to default it.
	It("leaves the destination alone for an empty value", func() {
		dst := []condition{{Match: "kept"}}
		flag := flags.JSONFlag(&cli.StringFlag{Name: "conditions"}, &dst)

		Expect(flag.Validator("")).To(Succeed())
		Expect(dst).To(Equal([]condition{{Match: "kept"}}))
	})

	It("names the flag in the error so the message points at the input", func() {
		dst := []condition{}
		flag := flags.JSONFlag(&cli.StringFlag{Name: "conditions"}, &dst)

		err := flag.Validator("{not json")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("conditions"))
	})

	// The defaults are non-empty literals that would otherwise never be parsed, so
	// a typo in one would only surface once a user overrode something else.
	It("makes the flag validate its own default", func() {
		dst := []condition{}
		flag := flags.JSONFlag(&cli.StringFlag{Name: "conditions"}, &dst)

		Expect(flag.ValidateDefaults).To(BeTrue())
	})
})

var _ = Describe("YAMLFlag", func() {
	It("unmarshals the value into the destination", func() {
		dst := []condition{}
		flag := flags.YAMLFlag(&cli.StringFlag{Name: "sanitize-tags"}, &dst)

		Expect(flag.Validator("- match: \"^tags/\"\n  environment: production\n")).To(Succeed())
		Expect(dst).To(Equal([]condition{{Match: "^tags/", Environment: "production"}}))
	})

	// The defaults are written as JSON but documented as YAML, which only works
	// because YAML is a superset of it.
	It("accepts the JSON the defaults are written in", func() {
		dst := []condition{}
		flag := flags.YAMLFlag(&cli.StringFlag{Name: "sanitize-tags"}, &dst)

		Expect(flag.Validator(`[{ "match": "^tags/", "environment": "production" }]`)).To(Succeed())
		Expect(dst).To(Equal([]condition{{Match: "^tags/", Environment: "production"}}))
	})

	It("names the flag in the error", func() {
		dst := []condition{}
		flag := flags.YAMLFlag(&cli.StringFlag{Name: "sanitize-tags"}, &dst)

		err := flag.Validator("\t- broken")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("sanitize-tags"))
	})
})

var _ = Describe("EnvVars", func() {
	// The chain is read in order, so a flag that answers to more than one name
	// resolves to the first of them a pipeline has set.
	It("builds one source per name, in the order they were given", func() {
		chain := flags.EnvVars("CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH")

		Expect(chain.Chain).To(HaveLen(2))
		Expect(chain.Chain[0].String()).To(ContainSubstring("CI_COMMIT_REF_NAME"))
		Expect(chain.Chain[1].String()).To(ContainSubstring("BITBUCKET_BRANCH"))
	})
})
