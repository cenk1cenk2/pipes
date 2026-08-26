package node

import (
	"fmt"
	"os"
	"strings"

	"github.com/cenk1cenk2/plumber/v6"
	"github.com/nochso/gomd/eol"
	ucli "github.com/urfave/cli/v3"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

//revive:disable:line-length-limit

const CATEGORY_LOGIN = "Login"

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

func NewLoginFlags(cfg *Login) []ucli.Flag {
	return []ucli.Flag{
		cli.JSONFlag(&ucli.StringFlag{
			Category: CATEGORY_LOGIN,
			Name:     "npm.login",
			Sources:  cli.EnvVars("NPM_LOGIN"),
			Usage:    "NPM registries to login. json([]struct { username: string, password: string, registry?: string, useHttps?: bool })",
			Required: false,
			Value:    "",
		}, &cfg.Entries),

		&ucli.StringSliceFlag{
			Category:    CATEGORY_LOGIN,
			Name:        "npm.npmrc_file",
			Sources:     cli.EnvVars("NPM_NPMRC_FILE"),
			Usage:       ".npmrc file to use.",
			Required:    false,
			Value:       []string{".npmrc"},
			Destination: &cfg.NpmRcFiles,
		},

		&ucli.StringFlag{
			Category:    CATEGORY_LOGIN,
			Name:        "npm.npmrc",
			Sources:     cli.EnvVars("NPM_NPMRC"),
			Usage:       "Direct contents of .npmrc file.",
			Required:    false,
			Value:       "",
			Destination: &cfg.NpmRc,
		},
	}
}

// LoginTaskList writes the configured credentials into the npmrc files and
// checks that the registries accept them.
func LoginTaskList(p *plumber.Plumber, cfg *Login) *plumber.TaskList {
	tl := &plumber.TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ *plumber.TaskList) error {
			return p.Validate(cfg)
		}).
		Set(func(tl *plumber.TaskList) plumber.Job {
			return plumber.JobSequence(
				generateNpmRc(tl, cfg).Job(),
				verifyNpmLogin(tl, cfg).Job(),
			)
		})
}

func generateNpmRc(tl *plumber.TaskList, cfg *Login) *plumber.Task {
	return tl.CreateTask("npmrc").
		ShouldDisable(func(_ *plumber.Task) bool {
			return cfg.Entries == nil && cfg.NpmRc == ""
		}).
		Set(func(t *plumber.Task) error {
			t.Log.Debugf(
				".npmrc file: %s", strings.Join(cfg.NpmRcFiles, ", "),
			)

			npmrc := []string{}

			if cfg.Entries != nil {
				t.Log.Infoln("Logging in to given registries with credentials.")

				for _, v := range cfg.Entries {
					t.Log.Infof(
						"Generating login credentials for the registry: %s",
						v.Registry,
					)

					npmrc = append(
						npmrc,
						fmt.Sprintf("//%s/:_authToken=%s", v.Registry, v.Token),
					)
				}
			}

			if cfg.NpmRc != "" {
				t.Log.Infoln("Appending directly to the given npmrc file.")

				npmrc = append(npmrc, strings.Split(cfg.NpmRc, eol.OSDefault().String())...)
			}

			for _, file := range cfg.NpmRcFiles {
				t.CreateSubtask(file).
					Set(
						func(st *plumber.Task) error {
							st.Log.Infof("Generating npmrc file: %s", file)

							f, err := os.OpenFile(file,
								os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

							if err != nil {
								return err
							}

							defer f.Close()

							if _, err := f.WriteString(strings.Join(npmrc, eol.OSDefault().String()) + eol.OSDefault().String()); err != nil {
								return err
							}

							return nil
						}).
					AddSelfToTheParentAsParallel()
			}

			return nil
		}).
		ShouldRunAfter(func(t *plumber.Task) error {
			return t.RunSubtasks()
		})
}

func verifyNpmLogin(tl *plumber.TaskList, cfg *Login) *plumber.Task {
	return tl.CreateTask("login").
		// Without an npmrc file there is nothing holding the credentials for npm
		// to read back, so there is nothing to verify either.
		ShouldDisable(func(_ *plumber.Task) bool {
			return cfg.Entries == nil || len(cfg.NpmRcFiles) == 0
		}).
		Set(func(t *plumber.Task) error {
			for _, v := range cfg.Entries {
				t.CreateCommand(
					"npm",
					"whoami",
				).
					SetLogLevel(plumber.LOG_LEVEL_DEBUG, plumber.LOG_LEVEL_DEFAULT, plumber.LOG_LEVEL_DEBUG).
					Set(func(c *plumber.Command) error {
						c.Log.Infof(
							"Checking login credentials for Npm registry: %s", v.Registry,
						)

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
		ShouldRunAfter(func(t *plumber.Task) error {
			return t.RunCommandJobAsJobParallel()
		})
}
