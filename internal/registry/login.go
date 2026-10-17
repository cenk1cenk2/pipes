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

	// LoginFlagsOptions is what a pipe builds the flags of its container registry login onto.
	LoginFlagsOptions struct {
		Category string
		// Alias is the pipe specific variable, read before the shared one when a pipe has one.
		Alias       string
		Destination *[]LoginEntry
	}

	// LoginOptions is what a pipe builds its container registry login onto.
	LoginOptions struct {
		// Command is the login command of the pipe, called with the uri and the credentials.
		Command []string
		Entries []LoginEntry
		// HelmRegistryConfig pins HELM_REGISTRY_CONFIG for the process when it is unset.
		HelmRegistryConfig string
	}

	// Login is the container registry login of a pipe, run as a task of the task list it is part of.
	Login struct {
		options LoginOptions
	}
)

func NewLoginFlags(opts LoginFlagsOptions) []cli.Flag {
	envs := []string{"CONTAINER_LOGIN"}
	if opts.Alias != "" {
		envs = slices.Insert(envs, 0, opts.Alias)
	}

	return []cli.Flag{
		flags.YAMLFlag(opts.Destination, &cli.StringFlag{
			Category: opts.Category,
			Name:     "container.login",
			Sources:  cli.EnvVars(envs...),
			Usage:    "Container registries to login to. format(yaml([]struct{ uri: string, username: string, password: string }))",
			Required: false,
		}),
	}
}

func NewLogin(opts LoginOptions) *Login {
	return &Login{options: opts}
}

// hasCredentials reports whether the entry has both a username and a password,
// which a protected variable leaves empty on an unprotected ref.
func (e LoginEntry) hasCredentials() bool {
	return e.Username != "" && e.Password != ""
}

// Covers reports whether an entry logs in to the image, either the registry
// host itself or a path the image lives under.
func (l *Login) Covers(image *Image) bool {
	target := image.String()

	return slices.ContainsFunc(l.options.Entries, func(e LoginEntry) bool {
		return e.hasCredentials() && (target == e.Uri || strings.HasPrefix(target, e.Uri+"/"))
	})
}

// Job is the login of the pipe as a task of the task list it runs in.
func (l *Login) Job(tl *TaskList) Job {
	return l.login(tl).Job()
}

func (l *Login) login(tl *TaskList) *Task {
	return tl.CreateTask("login", "container").
		ShouldDisable(func(_ *Task) bool {
			return len(l.options.Entries) == 0
		}).
		ShouldRunBefore(func(_ context.Context, t *Task) error {
			for i := range l.options.Entries {
				if err := t.Plumber.Validate(&l.options.Entries[i]); err != nil {
					return err
				}

				if !l.options.Entries[i].hasCredentials() {
					t.Log.Warn(fmt.Sprintf("Container registry has no credentials, skipping the login: %s", l.options.Entries[i].Uri))

					continue
				}

				t.Plumber.AppendSecrets(l.options.Entries[i].Password)
			}

			if !slices.ContainsFunc(l.options.Entries, LoginEntry.hasCredentials) {
				return nil
			}

			// kustomize runs helm with a temporary HELM_CONFIG_HOME, which hides the
			// registry config a login writes by default. It never overrides
			// HELM_REGISTRY_CONFIG though, and the build runs in this process, so
			// pinning it here is what the chart pulls read the login back from.
			if _, ok := os.LookupEnv("HELM_REGISTRY_CONFIG"); l.options.HelmRegistryConfig != "" && !ok {
				return os.Setenv("HELM_REGISTRY_CONFIG", l.options.HelmRegistryConfig)
			}

			return nil
		}).
		Set(func(_ context.Context, t *Task) error {
			for _, e := range l.options.Entries {
				if !e.hasCredentials() {
					continue
				}

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
