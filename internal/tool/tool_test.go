package tool_test

import (
	"regexp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/tool"
)

var terraform = tool.Spec{
	Name:           "terraform",
	Category:       "Project",
	FlagPrefix:     "terraform",
	EnvPrefix:      "TERRAFORM",
	VersionArgs:    []string{"version"},
	VersionPattern: regexp.MustCompile(`Terraform (v\d+\.\d+\.\d+)`),
}

var _ = Describe("NewFlags", func() {
	flag := func(spec tool.Spec, cfg *tool.Config) *cli.StringFlag {
		flags := tool.NewFlags(spec, cfg)
		Expect(flags).To(HaveLen(1))

		return flags[0].(*cli.StringFlag)
	}

	It("names the flag after the pipe and defaults to the current directory", func() {
		f := flag(terraform, &tool.Config{})

		Expect(f.Name).To(Equal("terraform.cwd"))
		Expect(f.Value).To(Equal("."))
		Expect(f.Category).To(Equal("Project"))
	})

	It("reads the environment name the prefix builds", func() {
		f := flag(terraform, &tool.Config{})

		Expect(f.Sources.Chain).To(HaveLen(1))
		Expect(f.Sources.Chain[0].String()).To(ContainSubstring("TERRAFORM_CWD"))
	})

	It("writes into the given destination", func() {
		cfg := &tool.Config{}

		Expect(flag(terraform, cfg).Destination).To(BeIdenticalTo(&cfg.Cwd))
	})
})

var _ = Describe("ParseVersion", func() {
	It("returns the first submatch of the pattern", func() {
		Expect(terraform.ParseVersion("Terraform v1.9.8\non linux_amd64")).To(Equal("v1.9.8"))
	})

	// The probe is only ever logged, and the tool has already proven it runs by
	// answering at all, so an unrecognised banner is reported rather than fatal.
	It("returns the whole output when the pattern does not match", func() {
		Expect(terraform.ParseVersion("something else entirely")).To(Equal("something else entirely"))
	})

	It("returns the whole output when there is no pattern", func() {
		spec := tool.Spec{Name: "go", VersionArgs: []string{"version"}}

		Expect(spec.ParseVersion("go version go1.27.0 linux/amd64\n")).To(Equal("go version go1.27.0 linux/amd64"))
	})

	It("trims the surrounding whitespace", func() {
		spec := tool.Spec{Name: "helm"}

		Expect(spec.ParseVersion("  v3.16.1\n\n")).To(Equal("v3.16.1"))
	})
})

var _ = Describe("NewCtx", func() {
	// The environment map is written to by the pipes that append an environment
	// task, so it has to be usable before any task has run.
	It("hands back a writable environment map", func() {
		ctx := tool.NewCtx()

		Expect(ctx.Env).NotTo(BeNil())

		ctx.Env["GOPATH"] = "/cache"
		Expect(ctx.Env).To(HaveKeyWithValue("GOPATH", "/cache"))
	})
})
