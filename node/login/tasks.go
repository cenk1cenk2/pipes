package login

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	utils "github.com/cenk1cenk2/ci-cd-pipes/utils"
	"github.com/nochso/gomd/eol"
	"github.com/sirupsen/logrus"
)

type Ctx struct {
	NpmLogin []NpmLoginJson
}

var Context Ctx

func VerifyVariables() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "verify-login", Skip: Pipe.Npm.Login == ""},
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

func VerifyNpmLogin() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{
			Context:        "login",
			Skip:           Pipe.Npm.Login == "",
			StdOutLogLevel: logrus.DebugLevel,
		},
		Task: func(t *utils.Task) error {
			var wg sync.WaitGroup
			wg.Add(len(Context.NpmLogin))

			for i, v := range Context.NpmLogin {
				go func(i int, v NpmLoginJson) {
					defer wg.Done()

					t.Log.Infoln(
						fmt.Sprintf("Checking login credentials for Npm registry: %s", v.Registry),
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
						Pipe.Npm.NpmRc,
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

func GenerateNpmRc() utils.Task {
	return utils.Task{
		Metadata: utils.TaskMetadata{Context: "generate-npmrc", Skip: Pipe.Npm.Login == ""},
		Task: func(t *utils.Task) error {
			var wg sync.WaitGroup
			wg.Add(len(Context.NpmLogin))

			t.Log.Debugln(fmt.Sprintf(".npmrc file: %s", Pipe.Npm.NpmRc))

			npmrc := []string{}

			for i, v := range Context.NpmLogin {
				go func(i int, v NpmLoginJson) {
					defer wg.Done()

					t.Log.Infoln(
						fmt.Sprintf(
							"Generating login credentials for Npm registry: %s",
							v.Registry,
						),
					)

					npmrc = append(npmrc, fmt.Sprintf("//%s/:_authToken=%s", v.Registry, v.Token))
				}(i, v)
			}

			wg.Wait()

			f, err := os.OpenFile(Pipe.Npm.NpmRc,
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				return err
			}

			defer f.Close()
			if _, err := f.WriteString(strings.Join(npmrc, eol.OSDefault().String()) + eol.OSDefault().String()); err != nil {
				return err
			}

			return nil
		},
	}
}
