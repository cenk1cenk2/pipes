package cli_test

import (
	"bytes"

	"github.com/cenk1cenk2/plumber/v6"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

// The registry MarkSecret writes to is package level, matching the package level
// flag variables the pipes declare, so the specs share one plumber and assert on
// what its logger actually prints rather than on registry contents.
var _ = Describe("MarkSecret", func() {
	var (
		p      *plumber.Plumber
		output *bytes.Buffer
	)

	BeforeEach(func() {
		p = plumber.NewPlumber(func(_ *plumber.Plumber) *ucli.Command {
			return &ucli.Command{Name: "test"}
		})

		output = &bytes.Buffer{}
		p.Log.SetOutput(output)
	})

	It("hands back the destination it was given", func() {
		var token string

		Expect(cli.MarkSecret(&token)).To(BeIdenticalTo(&token))
	})

	It("keeps a marked value out of the log once the pipe is validated", func() {
		pipe := struct{ Token string }{}
		cli.MarkSecret(&pipe.Token)
		pipe.Token = "glpat-not-in-the-log"

		Expect(cli.Validated(p, &pipe)).To(Succeed())

		p.Log.Infof("authenticating with %s", pipe.Token)
		Expect(output.String()).NotTo(ContainSubstring("glpat-not-in-the-log"))
		Expect(output.String()).To(ContainSubstring("[REDACTED]"))
	})

	// An unset optional secret would otherwise register the empty string, which
	// masks every message rather than none of them.
	It("does not mask on an empty value", func() {
		pipe := struct{ Token string }{}
		cli.MarkSecret(&pipe.Token)

		Expect(cli.Validated(p, &pipe)).To(Succeed())

		p.Log.Infoln("nothing sensitive here")
		Expect(output.String()).NotTo(ContainSubstring("[REDACTED]"))
	})
})

var _ = Describe("Validated", func() {
	var p *plumber.Plumber

	BeforeEach(func() {
		p = plumber.NewPlumber(func(_ *plumber.Plumber) *ucli.Command {
			return &ucli.Command{Name: "test"}
		})
		p.Log.SetOutput(GinkgoWriter)
	})

	It("applies the struct defaults", func() {
		pipe := struct {
			Format string `default:"oci"`
		}{}

		Expect(cli.Validated(p, &pipe)).To(Succeed())
		Expect(pipe.Format).To(Equal("oci"))
	})

	It("fails on a pipe that does not validate", func() {
		pipe := struct {
			Format string `validate:"oneof=oci docker"`
		}{Format: "tarball"}

		Expect(cli.Validated(p, &pipe)).To(HaveOccurred())
	})

	It("accepts a pipe that validates", func() {
		pipe := struct {
			Format string `validate:"oneof=oci docker"`
		}{Format: "oci"}

		Expect(cli.Validated(p, &pipe)).To(Succeed())
	})
})
