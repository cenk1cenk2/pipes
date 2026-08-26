package environment_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ucli "github.com/urfave/cli/v3"

	"gitlab.kilic.dev/devops/pipes/internal/environment"
	"gitlab.kilic.dev/devops/pipes/internal/git"
)

var _ = Describe("Select", func() {
	conditions := []environment.Condition{
		{Match: `^tags/`, Environment: "production"},
		{Match: `^heads/main$`, Environment: "develop"},
	}

	It("reports the environment of the matching condition", func() {
		Expect(environment.Select(conditions, []string{"heads/main"})).To(Equal("develop"))
	})

	// The conditions are the user's rule order, so an earlier rule wins even when a
	// later one matches a reference that comes first.
	It("ranks the conditions above the references", func() {
		Expect(environment.Select(conditions, []string{"heads/main", "tags/v1.0.0"})).
			To(Equal("production"))
	})

	// Nothing matching is not on its own an error, since it is the strict flag and
	// not this function that decides whether the pipe may carry on without one.
	It("reports an empty environment when nothing matches", func() {
		Expect(environment.Select(conditions, []string{"heads/feature/one"})).To(BeEmpty())
	})

	It("reports an empty environment for no conditions and for no references", func() {
		Expect(environment.Select(nil, []string{"heads/main"})).To(BeEmpty())
		Expect(environment.Select(conditions, nil)).To(BeEmpty())
	})

	It("names the pattern it could not compile", func() {
		_, err := environment.Select([]environment.Condition{{Match: `^tags/(`}}, []string{"tags/v1.0.0"})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(`^tags/(`))
	})

	// The default conditions ship as the flag default, so a typo in one would only
	// surface on the pipeline that happened to need it.
	It("selects through the shipped default conditions", func() {
		defaults := []environment.Condition{}
		Expect(json.Unmarshal([]byte(environment.DEFAULT_CONDITIONS), &defaults)).To(Succeed())

		Expect(environment.Select(defaults, git.Refs{Tag: "v1.0.0"}.References())).To(Equal("production"))
		Expect(environment.Select(defaults, git.Refs{Tag: "v1.0.0-rc.1"}.References())).To(Equal("stage"))
		Expect(environment.Select(defaults, git.Refs{Branch: "main"}.References())).To(Equal("develop"))
		Expect(environment.Select(defaults, git.Refs{Branch: "master"}.References())).To(Equal("develop"))
		Expect(environment.Select(defaults, git.Refs{Branch: "feature/one"}.References())).To(BeEmpty())
	})

	// A tagged pipeline carries the branch it was tagged on as well, and the tag is
	// what decides the environment of a release.
	It("selects the tag environment for a tagged pipeline on main", func() {
		defaults := []environment.Condition{}
		Expect(json.Unmarshal([]byte(environment.DEFAULT_CONDITIONS), &defaults)).To(Succeed())

		Expect(environment.Select(defaults, git.Refs{Branch: "main", Tag: "v1.0.0"}.References())).
			To(Equal("production"))
	})
})

var _ = Describe("Fetch", func() {
	It("strips the environment prefix off the variables that carry it", func() {
		Expect(environment.Fetch([]string{"STAGE_TOKEN=stage-token"}, "stage")).
			To(HaveKeyWithValue("TOKEN", "stage-token"))
	})

	// An environment is expected to override only the few variables it cares about,
	// so everything else has to reach the pipe untouched.
	It("keeps the variables without the prefix as they are", func() {
		Expect(environment.Fetch([]string{"PATH=/usr/bin"}, "stage")).
			To(HaveKeyWithValue("PATH", "/usr/bin"))
	})

	It("lets the prefixed value win over the plain one", func() {
		Expect(environment.Fetch([]string{"TOKEN=plain", "STAGE_TOKEN=stage-token"}, "stage")).
			To(HaveKeyWithValue("TOKEN", "stage-token"))
	})

	It("matches the prefix case insensitively against the selection", func() {
		Expect(environment.Fetch([]string{"PRODUCTION_TOKEN=token"}, "production")).
			To(HaveKeyWithValue("TOKEN", "token"))
	})

	It("names the selection in the variables it hands back", func() {
		Expect(environment.Fetch(nil, "stage")).To(Equal(map[string]string{"ENVIRONMENT": "stage"}))
	})

	It("keeps an empty value", func() {
		Expect(environment.Fetch([]string{"STAGE_TOKEN="}, "stage")).To(HaveKeyWithValue("TOKEN", ""))
	})

	It("keeps everything after the first separator", func() {
		Expect(environment.Fetch([]string{"STAGE_URL=https://host/?a=b"}, "stage")).
			To(HaveKeyWithValue("URL", "https://host/?a=b"))
	})

	It("skips an entry that is not a pair", func() {
		Expect(environment.Fetch([]string{"NOT_A_PAIR"}, "stage")).
			To(Equal(map[string]string{"ENVIRONMENT": "stage"}))
	})
})

var _ = Describe("NewFlags", func() {
	names := func(flags []ucli.Flag) []string {
		found := []string{}
		for _, flag := range flags {
			found = append(found, flag.Names()...)
		}

		return found
	}

	// The git flags come first because the references they carry are what the
	// conditions match against, and the order is what the documentation prints.
	It("registers the git flags ahead of the environment ones", func() {
		cfg := environment.Config{}

		Expect(names(environment.NewFlags(&cfg))).To(Equal([]string{
			"git.branch",
			"git.tag",
			"environment.enable",
			"environment.conditions",
			"environment.fail-on-no-reference",
			"environment.strict",
		}))
	})

	// A pipe reads the selection back off the same instance it registered, so a
	// flag landing on a copy would leave it on the zero value.
	It("binds each flag onto the given configuration", func() {
		cfg := environment.Config{}
		flags := environment.NewFlags(&cfg)

		//nolint:errcheck
		Expect(flags[0].(*ucli.StringFlag).Destination).To(BeIdenticalTo(&cfg.Git.Branch))
		//nolint:errcheck
		Expect(flags[2].(*ucli.BoolFlag).Destination).To(BeIdenticalTo(&cfg.Enable))
		//nolint:errcheck
		Expect(flags[4].(*ucli.BoolFlag).Destination).To(BeIdenticalTo(&cfg.FailOnNoReference))
		//nolint:errcheck
		Expect(flags[5].(*ucli.BoolFlag).Destination).To(BeIdenticalTo(&cfg.Strict))
	})

	// The conditions flag is a JSON string the pipe reads back as a struct, and its
	// own default is the value most pipelines end up running with.
	It("unmarshals the conditions onto the given configuration", func() {
		cfg := environment.Config{}
		flags := environment.NewFlags(&cfg)

		//nolint:errcheck
		Expect(flags[3].(*ucli.StringFlag).Validator(environment.DEFAULT_CONDITIONS)).To(Succeed())
		Expect(cfg.Conditions).To(HaveLen(4))
		Expect(cfg.Conditions[0].Environment).To(Equal("production"))
	})

	// Two pipes each register their own configuration, so one call handing back the
	// flags of another would bind both onto whichever ran last.
	It("builds a fresh set of flags per configuration", func() {
		first, second := environment.Config{}, environment.Config{}

		//nolint:errcheck
		Expect(environment.NewFlags(&first)[2].(*ucli.BoolFlag).Destination).
			NotTo(BeIdenticalTo(environment.NewFlags(&second)[2].(*ucli.BoolFlag).Destination))
	})
})
