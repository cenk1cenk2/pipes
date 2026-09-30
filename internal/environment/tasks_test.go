package environment_test

import (
	"bytes"
	"io"

	. "github.com/cenk1cenk2/plumber/v7"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/git"
)

var _ = Describe("SetupTaskList", func() {
	conditions := []environment.Condition{
		{Match: `^tags/`, Environment: "production"},
		{Match: `^heads/main$`, Environment: "develop"},
	}

	// the specs never read the environment of the process running them, which is
	// whatever the machine or the runner holds, secrets included.
	environ := []string{"DEVELOP_TOKEN=develop-token", "PATH=/usr/bin"}

	var log *bytes.Buffer

	run := func(cfg environment.Config) (environment.Ctx, error) {
		GinkgoHelper()

		ctx := environment.Ctx{}
		cfg.Environ = environ
		log = &bytes.Buffer{}

		p := NewPlumber(func(_ *Plumber) *cli.Command {
			return &cli.Command{Name: "test"}
		})
		p.SetLoggerOutput(io.MultiWriter(log, GinkgoWriter))

		return ctx, p.RunJobs(environment.SetupTaskList(p, &cfg, &ctx).Job())
	}

	It("reads the variables of the selected environment into the context", func() {
		ctx, err := run(environment.Config{
			Enable:     true,
			Conditions: conditions,
			Git:        git.Refs{Branch: "main"},
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(ctx.References).To(Equal([]string{"heads/main"}))
		Expect(ctx.Environment).To(Equal("develop"))
		Expect(ctx.EnvVars).To(HaveKeyWithValue("TOKEN", "develop-token"))
		Expect(ctx.EnvVars).To(HaveKeyWithValue("PATH", "/usr/bin"))
	})

	It("names the variables it selected in the log and never their values", func() {
		_, err := run(environment.Config{
			Enable:     true,
			Conditions: conditions,
			Git:        git.Refs{Branch: "main"},
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(log.String()).To(ContainSubstring("develop -> [TOKEN]"))
		Expect(log.String()).NotTo(ContainSubstring("develop-token"))
		Expect(log.String()).NotTo(ContainSubstring("/usr/bin"))
	})

	// the pipes that only inject an environment on request ship the flag off, and a
	// disabled list must leave the context for the rest of the pipe untouched.
	It("resolves nothing at all when the feature is off", func() {
		ctx, err := run(environment.Config{
			Conditions: conditions,
			Git:        git.Refs{Branch: "main"},
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(ctx).To(Equal(environment.Ctx{}))
	})

	// nothing selected is the normal case for a feature branch, so the variables are
	// left alone rather than the pipeline failed.
	It("carries on without a selection", func() {
		ctx, err := run(environment.Config{
			Enable:     true,
			Conditions: conditions,
			Git:        git.Refs{Branch: "feature/one"},
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(ctx.Environment).To(BeEmpty())
		Expect(ctx.EnvVars).To(BeNil())
	})

	// the plumber terminates the process on a task error, so the strict and the
	// fail-on-no-reference paths are asserted on Select and References instead.
})
