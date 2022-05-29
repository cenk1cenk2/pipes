package login

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nochso/gomd/eol"
	"github.com/sirupsen/logrus"
	"github.com/workanator/go-floc/v3"
	. "gitlab.kilic.dev/libraries/plumber/v2"
)

type Ctx struct {
	NpmLogin []NpmLoginJson
}

func Decode(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("decode").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Npm.Login == ""
		}).
		Set(func(t *Task[Pipe], c floc.Control) error {
			// unmarshal npm logins and use the default registry for ones that are not defined
			t.Log.Debugln("Npm login credentials are specified, initiating login context.")

			if err := json.Unmarshal([]byte(t.Pipe.Npm.Login), &t.Pipe.Ctx.NpmLogin); err != nil {
				t.Log.Fatalln("Can not decode Npm registry login credentials.")
			}

			for i := range t.Pipe.Ctx.NpmLogin {
				t.CreateSubtask("validate").Set(func(st *Task[Pipe], c floc.Control) error {
					return st.TaskList.Validate(&st.Pipe.Ctx.NpmLogin[i])
				}).ToParent(t, func(pt, st *Task[Pipe]) {
					pt.ExtendSubtask(func(j floc.Job) floc.Job {
						return tl.JobParallel(j, st.Job())
					})
				})
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe], c floc.Control) error {
			return t.RunSubtasks()
		})
}

func GenerateNpmRc(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("npmrc").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Npm.Login == "" && t.Pipe.Npm.NpmRc == ""
		}).
		Set(func(t *Task[Pipe], c floc.Control) error {
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
				t.CreateSubtask("generate").Set(func(st *Task[Pipe], c floc.Control) error {
					st.Log.Debugf("Creating npmrc file: %s", file)

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
				}).ToParent(t, func(pt, st *Task[Pipe]) {
					pt.ExtendSubtask(func(j floc.Job) floc.Job {
						return pt.TaskList.JobParallel(j, st.Job())
					})
				})
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe], c floc.Control) error {
			return t.RunSubtasks()
		})
}

func VerifyNpmLogin(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("login").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.Npm.Login == ""
		}).
		Set(func(t *Task[Pipe], c floc.Control) error {
			for _, v := range t.Pipe.Ctx.NpmLogin {
				t.CreateCommand("npm", "whoami").
					SetLogLevel(logrus.DebugLevel, 0).
					Set(func(c *Command[Pipe]) error {

						t.Log.Infof(
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
					}).AddSelfToTheTask()
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe], c floc.Control) error {
			return t.RunCommandJobAsJobParallel()
		})
}
