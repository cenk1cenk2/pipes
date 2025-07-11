package login

import (
	"fmt"
	"os"
	"strings"

	. "github.com/cenk1cenk2/plumber/v6"
	"github.com/nochso/gomd/eol"
)

func GenerateNpmRc(tl *TaskList) *Task {
	return tl.CreateTask("npmrc").
		ShouldDisable(func(t *Task) bool {
			return P.Npm.Login == nil && P.Npm.NpmRc == ""
		}).
		Set(func(t *Task) error {
			t.Log.Debugf(
				".npmrc file: %s", strings.Join(P.Npm.NpmRcFile, ", "),
			)

			npmrc := []string{}

			if P.Npm.Login != nil {
				t.Log.Infoln("Logging in to given registries with credentials.")

				for _, v := range P.Npm.Login {
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

			if P.Npm.NpmRc != "" {
				t.Log.Infoln("Appending directly to the given npmrc file.")

				npmrc = append(npmrc, strings.Split(P.Npm.NpmRc, eol.OSDefault().String())...)
			}

			for _, file := range P.Npm.NpmRcFile {
				t.CreateSubtask(file).
					Set(
						func(st *Task) error {
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
		ShouldRunAfter(func(t *Task) error {
			return t.RunSubtasks()
		})
}

func VerifyNpmLogin(tl *TaskList) *Task {
	return tl.CreateTask("login").
		ShouldDisable(func(t *Task) bool {
			return P.Npm.Login == nil
		}).
		Set(func(t *Task) error {
			for _, v := range P.Npm.Login {
				t.CreateCommand(
					"npm",
					"whoami",
				).
					SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
					Set(func(c *Command) error {
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
							P.Npm.NpmRcFile[0],
							"--registry",
							url,
						)

						return nil
					}).
					AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task) error {
			return t.RunCommandJobAsJobParallel()
		})
}
