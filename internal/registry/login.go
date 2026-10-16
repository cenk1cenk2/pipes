package registry

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/flags"
)

//revive:disable:line-length-limit

type (
	// LoginEntry is one container registry the pipe authenticates against.
	LoginEntry struct {
		Uri      string `yaml:"uri"      validate:"required"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}

	// LoginOptions is what a pipe builds its container registry login onto.
	LoginOptions struct {
		Category string
		// Alias is the pipe specific variable, read before the shared one.
		Alias string
		// Command is the login command of the pipe, called with the uri and the credentials.
		Command []string
		// HelmRegistryConfig pins HELM_REGISTRY_CONFIG for the process when it is unset.
		HelmRegistryConfig string
	}

	// Login is the container registry login of a pipe, its flags and the entries they fill.
	Login struct {
		Flags   []cli.Flag
		Entries []LoginEntry

		options LoginOptions
	}
)

func NewLogin(opts LoginOptions) *Login {
	l := &Login{options: opts}

	l.Flags = []cli.Flag{
		flags.YAMLFlag(&l.Entries, &cli.StringFlag{
			Category: opts.Category,
			Name:     "container.login",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar(opts.Alias),
				cli.EnvVar("CONTAINER_LOGIN"),
			),
			Usage:    "Container registries to login to. format(yaml([]struct{ uri: string, username: string, password: string }))",
			Required: false,
		}),
	}

	return l
}

// Covers reports whether an entry logs in to the target, either the registry
// host itself or a path the target lives under.
func (l *Login) Covers(target string) bool {
	return slices.ContainsFunc(l.Entries, func(e LoginEntry) bool {
		return target == e.Uri || strings.HasPrefix(target, e.Uri+"/")
	})
}

func (l *Login) New(p *Plumber) *TaskList {
	tl := &TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, _ *TaskList) error {
			for i := range l.Entries {
				if err := p.Validate(&l.Entries[i]); err != nil {
					return err
				}
			}

			// kustomize runs helm with a temporary HELM_CONFIG_HOME, which hides the
			// registry config a login writes by default. It never overrides
			// HELM_REGISTRY_CONFIG though, and the build runs in this process, so
			// pinning it here is what the chart pulls read the login back from.
			if _, ok := os.LookupEnv("HELM_REGISTRY_CONFIG"); l.options.HelmRegistryConfig != "" && len(l.Entries) > 0 && !ok {
				return os.Setenv("HELM_REGISTRY_CONFIG", l.options.HelmRegistryConfig)
			}

			return nil
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				l.login(tl).Job(),
			)
		})
}

func (l *Login) login(tl *TaskList) *Task {
	return tl.CreateTask("login", "container").
		ShouldDisable(func(_ *Task) bool {
			return len(l.Entries) == 0
		}).
		ShouldRunBefore(func(_ context.Context, t *Task) error {
			for _, e := range l.Entries {
				t.Plumber.AppendSecrets(e.Password)
			}

			return nil
		}).
		Set(func(_ context.Context, t *Task) error {
			for _, e := range l.Entries {
				t.CreateCommand(
					l.options.Command[0],
					slices.Concat(l.options.Command[1:], []string{e.Uri, "--username", e.Username, "--password-stdin"})...,
				).
					SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDefault).
					Set(func(_ context.Context, c *Command) error {
						c.Log.Info(fmt.Sprintf("Logging in to container registry: %s", e.Uri))

						return nil
					}).
					SetStdin(func(_ *Command) io.Reader {
						return strings.NewReader(e.Password)
					}).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}
