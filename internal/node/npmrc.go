package node

import (
	"context"
	"fmt"
	"os"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/nochso/gomd/eol"
)

//revive:disable:line-length-limit

type (
	// LoginEntry is one registry the pipe authenticates against.
	LoginEntry struct {
		Username string `json:"username"           validate:"required"`
		Token    string `json:"token"              validate:"required"`
		Registry string `json:"registry,omitempty"                     default:"registry.npmjs.org"`
		UseHttps bool   `json:"useHttps,omitempty"                     default:"true"`
	}

	// Login is the npmrc a pipe writes before it reaches for a registry.
	Login struct {
		Entries    []LoginEntry
		NpmRcFiles []string
		NpmRc      string
	}
)

// LoginTaskList writes the configured credentials into the npmrc files and
// checks that the registries accept them.
func LoginTaskList(p *Plumber, cfg *Login) *TaskList {
	tl := &TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ context.Context, _ *TaskList) error {
			return p.Validate(cfg)
		}).
		Set(func(tl *TaskList) Job {
			return JobSequence(
				npmrc(tl, cfg).Job(),
				login(tl, cfg).Job(),
			)
		})
}

func npmrc(tl *TaskList, cfg *Login) *Task {
	return tl.CreateTask("npmrc").
		ShouldDisable(func(_ *Task) bool {
			return cfg.Entries == nil && cfg.NpmRc == ""
		}).
		Set(func(_ context.Context, t *Task) error {
			t.Log.Debug(fmt.Sprintf(".npmrc file: %s", strings.Join(cfg.NpmRcFiles, ", ")))

			lines := []string{}

			if cfg.Entries != nil {
				t.Log.Info("Logging in to given registries with credentials.")

				for _, v := range cfg.Entries {
					t.Log.Info(fmt.Sprintf("Generating login credentials for the registry: %s",
						v.Registry),
					)

					lines = append(
						lines,
						fmt.Sprintf("//%s/:_authToken=%s", v.Registry, v.Token),
					)
				}
			}

			if cfg.NpmRc != "" {
				t.Log.Info("Appending directly to the given npmrc file.")

				lines = append(lines, strings.Split(cfg.NpmRc, eol.OSDefault().String())...)
			}

			for _, file := range cfg.NpmRcFiles {
				t.CreateSubtask(file).
					Set(
						func(_ context.Context, st *Task) error {
							st.Log.Info(fmt.Sprintf("Generating npmrc file: %s", file))

							f, err := os.OpenFile(file,
								os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

							if err != nil {
								return err
							}

							defer f.Close()

							if _, err := f.WriteString(strings.Join(lines, eol.OSDefault().String()) + eol.OSDefault().String()); err != nil {
								return err
							}

							return nil
						}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunSubtasks(ctx)
		})
}

func login(tl *TaskList, cfg *Login) *Task {
	return tl.CreateTask("login").
		// without an npmrc file there is nothing holding the credentials for npm
		// to read back, so there is nothing to verify either.
		ShouldDisable(func(_ *Task) bool {
			return cfg.Entries == nil || len(cfg.NpmRcFiles) == 0
		}).
		Set(func(_ context.Context, t *Task) error {
			for _, v := range cfg.Entries {
				t.CreateCommand(
					"npm",
					"whoami",
				).
					SetLogLevel(LogLevelDebug, LogLevelDefault, LogLevelDebug).
					Set(func(_ context.Context, c *Command) error {
						c.Log.Info(fmt.Sprintf("Checking login credentials for Npm registry: %s", v.Registry))

						var url string

						if v.UseHttps {
							url = fmt.Sprintf("https://%s", v.Registry)
						} else {
							url = fmt.Sprintf("http://%s", v.Registry)
						}

						c.AppendArgs(
							"--configfile",
							cfg.NpmRcFiles[0],
							"--registry",
							url,
						)

						return nil
					}).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobParallel(ctx)
		})
}
