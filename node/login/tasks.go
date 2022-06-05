package login

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nochso/gomd/eol"
	"github.com/sirupsen/logrus"
	. "gitlab.kilic.dev/libraries/plumber/v3"
)

type Ctx struct {
	NpmLogin []NpmLoginJson
}

func Decode(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("init").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Npm.Login == ""
		}).
		Set(func(t *Task[Pipe]) error {
			// unmarshal npm logins and use the default registry for ones that are not defined
			t.Log.Infoln("Npm login credentials are specified, initiating login process.")

			if err := json.Unmarshal([]byte(t.Pipe.Npm.Login), &t.Pipe.Ctx.NpmLogin); err != nil {
				t.Log.Fatalln("Can not decode Npm registry login credentials.")
			}

			if err := tl.Validate(&t.Pipe.Ctx); err != nil {
				return err
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunSubtasks()
		})
}

func GenerateNpmRc(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("npmrc").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Npm.Login == "" && t.Pipe.Npm.NpmRc == ""
		}).
		Set(func(t *Task[Pipe]) error {
			t.Log.Debugf(
				".npmrc file: %s", strings.Join(t.Pipe.Npm.NpmRcFile.Value(), ", "),
			)

			npmrc := []string{}

			if t.Pipe.Npm.Login != "" {
				t.Log.Infoln("Logging in to given registries with credentials.")

				for _, v := range t.Pipe.Ctx.NpmLogin {
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

			for _, file := range t.Pipe.Npm.NpmRcFile.Value() {
				t.CreateSubtask("generate").
					Set(
						func(file string) TaskFn[Pipe] {
							return func(st *Task[Pipe]) error {
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
							}
						}(file)).
					AddSelfToParent(func(pt *Task[Pipe], st *Task[Pipe]) {
						pt.ExtendSubtask(func(j Job) Job {
							return tl.JobParallel(j, st.Job())
						})
					})
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
			return t.Pipe.Npm.Login == ""
		}).
		Set(func(t *Task[Pipe]) error {
			for _, v := range t.Pipe.Ctx.NpmLogin {
				t.CreateCommand("npm", "whoami").
					SetLogLevel(logrus.DebugLevel, 0, logrus.DebugLevel).
					Set(func(v NpmLoginJson) CommandFn[Pipe] {
						return func(c *Command[Pipe]) error {
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
								t.Pipe.Npm.NpmRcFile.Value()[0],
								"--registry",
								url,
							)

							return nil
						}
					}(v)).
					AddSelfToTheTask()

			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobParallel()
		})
}
