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

func Unmarshal(tl *TaskList[Pipe, Ctx]) *Task[Pipe, Ctx] {
	t := Task[Pipe, Ctx]{}

	return t.New(tl, "unmarshal").ShouldDisable(func(t *Task[Pipe, Ctx]) bool {
		return t.Pipe.Npm.Login == ""
	}).Set(func(t *Task[Pipe, Ctx], c floc.Control) error {
		// unmarshal npm logins and use the default registry for ones that are not defined
		t.Log.Debugln("Npm login credentials are specified, initiating login context.")

		if err := json.Unmarshal([]byte(t.Pipe.Npm.Login), &t.Context.NpmLogin); err != nil {
			t.Log.Fatalln("Can not decode Npm registry login credentials.")
		}

		for i := range t.Context.NpmLogin {
			st := Task[Pipe, Ctx]{}

			st.New(tl, "validate").Set(func(t *Task[Pipe, Ctx], c floc.Control) error {
				return t.TaskList.Validate(&t.Context.NpmLogin[i])
			})

			t.SetSubtask(st.Job())
		}

		return nil
	}).ShouldRunAfter(func(t *Task[Pipe, Ctx], c floc.Control) error {
		return t.RunSubtasks()
	})
}

func GenerateNpmRc(tl *TaskList[Pipe, Ctx]) *Task[Pipe, Ctx] {
	t := Task[Pipe, Ctx]{}

	return t.New(tl, "npmrc").ShouldDisable(func(t *Task[Pipe, Ctx]) bool {
		return t.Pipe.Npm.Login == "" && t.Pipe.Npm.NpmRc == ""
	}).Set(func(t *Task[Pipe, Ctx], c floc.Control) error {
		t.Log.Debugf(
			".npmrc file: %s", strings.Join(t.Pipe.Npm.NpmRcFile.Value(), ", "),
		)

		npmrc := []string{}

		if t.Pipe.Npm.Login != "" {
			t.Log.Infoln("Logging in to given registries with credentials.")

			for _, v := range t.Context.NpmLogin {
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
			st := Task[Pipe, Ctx]{}

			st.New(tl, "generate").Set(func(t *Task[Pipe, Ctx], c floc.Control) error {
				t.Log.Debugf("Creating npmrc file: %s", file)

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
			})

			t.ExtendSubtask(func(j floc.Job) floc.Job {
				return t.TaskList.JobParallel(j, st.Job())
			})
		}

		return nil
	}).ShouldRunAfter(func(t *Task[Pipe, Ctx], c floc.Control) error {
		return t.RunSubtasks()
	})
}

func VerifyNpmLogin(tl *TaskList[Pipe, Ctx]) *Task[Pipe, Ctx] {
	t := Task[Pipe, Ctx]{}

	return t.New(tl, "login").ShouldDisable(func(t *Task[Pipe, Ctx]) bool {
		return t.Pipe.Npm.Login == ""
	}).Set(func(t *Task[Pipe, Ctx], c floc.Control) error {
		for _, v := range t.Context.NpmLogin {
			cmd := Command[Pipe, Ctx]{}

			cmd.New(t, "npm", "whoami").SetLogLevel(logrus.DebugLevel, 0).Set(func(c *Command[Pipe, Ctx]) error {

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
			})

			t.AddCommands(cmd)
		}

		return nil
	}).ShouldRunAfter(func(t *Task[Pipe, Ctx], c floc.Control) error {
		return t.RunCommandJobAsJobParallel()
	})
}
