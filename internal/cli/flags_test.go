package cli_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

type condition struct {
	Match       string `json:"match"       yaml:"match"`
	Environment string `json:"environment" yaml:"environment"`
}

var _ = Describe("JSONFlag", func() {
	It("unmarshals the value into the destination", func() {
		dst := []condition{}
		flag := cli.JSONFlag(&ucli.StringFlag{Name: "conditions"}, &dst)

		Expect(flag.Validator(`[{ "match": "^heads/main$", "environment": "develop" }]`)).To(Succeed())
		Expect(dst).To(Equal([]condition{{Match: "^heads/main$", Environment: "develop"}}))
	})

	// Most of these flags are optional, so an unset one leaves the pipe on its zero
	// value rather than failing before the pipe has a chance to default it.
	It("leaves the destination alone for an empty value", func() {
		dst := []condition{{Match: "kept"}}
		flag := cli.JSONFlag(&ucli.StringFlag{Name: "conditions"}, &dst)

		Expect(flag.Validator("")).To(Succeed())
		Expect(dst).To(Equal([]condition{{Match: "kept"}}))
	})

	It("names the flag in the error so the message points at the input", func() {
		dst := []condition{}
		flag := cli.JSONFlag(&ucli.StringFlag{Name: "conditions"}, &dst)

		err := flag.Validator("{not json")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("conditions"))
	})

	// The defaults are non-empty literals that would otherwise never be parsed, so
	// a typo in one would only surface once a user overrode something else.
	It("makes the flag validate its own default", func() {
		dst := []condition{}
		flag := cli.JSONFlag(&ucli.StringFlag{Name: "conditions"}, &dst)

		Expect(flag.ValidateDefaults).To(BeTrue())
	})

	It("hands back the same flag it was given", func() {
		dst := []condition{}
		flag := &ucli.StringFlag{Name: "conditions"}

		Expect(cli.JSONFlag(flag, &dst)).To(BeIdenticalTo(flag))
	})
})

var _ = Describe("YAMLFlag", func() {
	It("unmarshals the value into the destination", func() {
		dst := []condition{}
		flag := cli.YAMLFlag(&ucli.StringFlag{Name: "sanitize-tags"}, &dst)

		Expect(flag.Validator("- match: \"^tags/\"\n  environment: production\n")).To(Succeed())
		Expect(dst).To(Equal([]condition{{Match: "^tags/", Environment: "production"}}))
	})

	// The defaults are written as JSON but documented as YAML, which only works
	// because YAML is a superset of it.
	It("accepts the JSON the defaults are written in", func() {
		dst := []condition{}
		flag := cli.YAMLFlag(&ucli.StringFlag{Name: "sanitize-tags"}, &dst)

		Expect(flag.Validator(`[{ "match": "^tags/", "environment": "production" }]`)).To(Succeed())
		Expect(dst).To(Equal([]condition{{Match: "^tags/", Environment: "production"}}))
	})

	It("names the flag in the error", func() {
		dst := []condition{}
		flag := cli.YAMLFlag(&ucli.StringFlag{Name: "sanitize-tags"}, &dst)

		err := flag.Validator("\t- broken")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("sanitize-tags"))
	})
})

var _ = Describe("EnvVars", func() {
	It("builds one source per name", func() {
		Expect(cli.EnvVars("CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH").Chain).To(HaveLen(2))
	})

	// The chain is read in order, so the name a pipeline already sets has to come
	// first for it to keep winning over the name that replaced it.
	It("keeps the given order", func() {
		chain := cli.EnvVars("CI_COMMIT_REF_NAME", "BITBUCKET_BRANCH")

		Expect(chain.Chain[0].String()).To(ContainSubstring("CI_COMMIT_REF_NAME"))
		Expect(chain.Chain[1].String()).To(ContainSubstring("BITBUCKET_BRANCH"))
	})

	It("builds an empty chain for no names", func() {
		Expect(cli.EnvVars().Chain).To(BeEmpty())
	})
})
