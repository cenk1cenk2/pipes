package login

import (
	"github.com/cenk1cenk2/plumber/v6"
)

// ContainerRegistryLoginVerify is the counterpart of the shared login task, for a
// pipeline that carries no credentials and is relying on an ambient login: it
// proves that login actually exists before the build spends time on an image it
// cannot push.
func ContainerRegistryLoginVerify(tl *plumber.TaskList) *plumber.Task {
	return tl.CreateTask("login", "verify").
		ShouldDisable(func(_ *plumber.Task) bool {
			return P.Username != "" && P.Password != ""
		}).
		Set(func(t *plumber.Task) error {
			t.CreateCommand(
				"buildah",
				"login",
				P.Uri,
			).
				SetLogLevel(plumber.LOG_LEVEL_DEBUG, plumber.LOG_LEVEL_DEFAULT, plumber.LOG_LEVEL_DEFAULT).
				Set(func(c *plumber.Command) error {
					c.Log.Debugf(
						"Will verify authentication in to container registry: %s",
						P.Uri,
					)

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *plumber.Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}
