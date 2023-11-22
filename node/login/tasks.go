package login

import (
	"fmt"
	"os"
	"strings"

	"github.com/nochso/gomd/eol"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

func GenerateNpmRc(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("npmrc").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Npm.Login == nil && t.Pipe.Npm.NpmRc == ""
		}).
		Set(func(t *Task[Pipe]) error {
			t.Log.Debugf(
				".npmrc file: %s", strings.Join(t.Pipe.Npm.NpmRcFile, ", "),
			)

			npmrc := []string{}

			if t.Pipe.Npm.Login != nil {
				t.Log.Infoln("Logging in to given registries with credentials.")

				for _, v := range t.Pipe.Npm.Login {
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

			if t.Pipe.Npm.NpmRc != "" {
				t.Log.Infoln("Appending directly to the given npmrc file.")

				npmrc = append(npmrc, strings.Split(t.Pipe.Npm.NpmRc, eol.OSDefault().String())...)
			}

			for _, file := range t.Pipe.Npm.NpmRcFile {
				func(file string) {
					t.CreateSubtask(file).
						Set(
							func(st *Task[Pipe]) error {
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
				}(file)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}

func VerifyNpmLogin(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("login").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Npm.Login == nil
		}).
		Set(func(t *Task[Pipe]) error {
			for _, v := range t.Pipe.Npm.Login {
				func(v NpmLoginJson) {
					t.CreateCommand(
						"npm",
						"whoami",
					).
						SetLogLevel(LOG_LEVEL_DEBUG, LOG_LEVEL_DEFAULT, LOG_LEVEL_DEBUG).
						Set(func(c *Command[Pipe]) error {
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
								t.Pipe.Npm.NpmRcFile[0],
								"--registry",
								url,
							)

							return nil
						}).
						AddSelfToTheTask()
				}(v)
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobParallel()
		})
}
