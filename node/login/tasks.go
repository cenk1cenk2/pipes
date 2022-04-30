package login

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/nochso/gomd/eol"
	"github.com/sirupsen/logrus"
	utils "gitlab.kilic.dev/libraries/go-utils/cli_utils"
)

type Ctx struct {
	NpmLogin []NpmLoginJson
}

var Context Ctx

func VerifyVariables() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "verify", Skip: Pipe.Npm.Login == ""},
		Task: func(t *utils.Task) error {
			err := utils.ValidateAndSetDefaults(t.Metadata, &Pipe)

			if err != nil {
				return err
			}

			// unmarshal npm logins and use the default registry for ones that are not defined
			t.Log.Debugln("Npm login credentials are specified, initiating login context.")

			err = json.Unmarshal([]byte(Pipe.Npm.Login), &Context.NpmLogin)

			if err != nil {
				t.Log.Fatalln("Can not decode Npm registry login credentials.")
			}

			var wg sync.WaitGroup
			wg.Add(len(Context.NpmLogin))
			errs := []error{}

			for i, v := range Context.NpmLogin {
				go func(i int, v NpmLoginJson) {
					defer wg.Done()

					err := utils.ValidateAndSetDefaults(t.Metadata, &Context.NpmLogin[i])

					if err != nil {
						errs = append(errs, err)
					}
				}(i, v)
			}

			wg.Wait()

			if len(errs) > 0 {
				for _, v := range errs {
					t.Log.Errorln(v)
				}

				t.Log.Fatalln("Errors encountered while validation.")
			}

			return nil
		},
	}
}

func GenerateNpmRc() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{
			Context: "generate-npmrc",
			Skip:    Pipe.Npm.Login == "" && Pipe.Npm.NpmRc == "",
		},
		Task: func(t *utils.Task) error {
			t.Log.Debugf(
				".npmrc file: %s", strings.Join(Pipe.Npm.NpmRcFile.Value(), ", "),
			)

			npmrc := []string{}

			if Pipe.Npm.Login != "" {
				t.Log.Infoln("Logging in to given registries with credentials.")

				var wg sync.WaitGroup
				wg.Add(len(Context.NpmLogin))
				for i, v := range Context.NpmLogin {
					go func(i int, v NpmLoginJson) {
						defer wg.Done()

						t.Log.Infof(
							"Generating login credentials for the registry: %s",
							v.Registry,
						)

						npmrc = append(
							npmrc,
							fmt.Sprintf("//%s/:_authToken=%s", v.Registry, v.Token),
						)
					}(i, v)
				}

				wg.Wait()
			}

			if Pipe.Npm.NpmRc != "" {
				t.Log.Infoln("Appending the given npmrc file.")

				npmrc = append(npmrc, strings.Split(Pipe.Npm.NpmRc, eol.OSDefault().String())...)
			}

			var wg sync.WaitGroup
			wg.Add(len(Pipe.Npm.NpmRcFile.Value()))
			errs := []error{}

			for i, v := range Pipe.Npm.NpmRcFile.Value() {
				go func(i int, file string) {
					defer wg.Done()

					t.Log.Debugf("Creating npmrc file: %s", file)

					f, err := os.OpenFile(file,
						os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

					if err != nil {
						errs = append(errs, err)

						return
					}

					defer f.Close()
					if _, err := f.WriteString(strings.Join(npmrc, eol.OSDefault().String()) + eol.OSDefault().String()); err != nil {
						errs = append(errs, err)

						return
					}
				}(i, v)
			}

			wg.Wait()

			if len(errs) > 0 {
				for _, v := range errs {
					t.Log.Errorln(v)
				}

				t.Log.Fatalln("Errors encountered while creating npmrc files.")
			}

			return nil
		},
	}
}

func VerifyNpmLogin() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{
			Context:        "verify-login",
			Skip:           Pipe.Npm.Login == "",
			StdOutLogLevel: logrus.DebugLevel,
		},
		Task: func(t *utils.Task) error {
			var wg sync.WaitGroup
			wg.Add(len(Context.NpmLogin))

			for i, v := range Context.NpmLogin {
				go func(i int, v NpmLoginJson) {
					defer wg.Done()

					t.Log.Infof(
						"Checking login credentials for Npm registry: %s", v.Registry,
					)

					cmd := exec.Command("npm", "whoami")

					var url string

					if v.UseHttps {
						url = fmt.Sprintf("https://%s", v.Registry)
					} else {
						url = fmt.Sprintf("http://%s", v.Registry)
					}

					cmd.Args = append(
						cmd.Args,
						"--configfile",
						Pipe.Npm.NpmRcFile.Value()[0],
						"--registry",
						url,
					)

					t.Commands = append(t.Commands, cmd)
				}(i, v)
			}

			wg.Wait()

			return nil
		},
	}
}
