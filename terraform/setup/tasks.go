package setup

import (
	"github.com/cenk1cenk2/plumber/v6"
)

func GenerateTerraformEnvVars(tl *plumber.TaskList) *plumber.Task {
	return tl.CreateTask("environment").
		Set(func(t *plumber.Task) error {
			C.Env["TF_IN_AUTOMATION"] = "true"

			C.Env["TF_LOG"] = P.LogLevel

			C.Env["TF_VAR_CI_API_V4_URL"] = P.CiVariables.ApiUrl
			C.Env["TF_VAR_CI_PROJECT_ID"] = P.CiVariables.ProjectId

			t.Log.Debugf("Generated following environment variables for terraform to consume: %+v", C.Env)

			return nil
		})
}
