package login

import (
	"context"

	"github.com/docker/docker/api/types"
	"gitlab.kilic.dev/devops/pipes/docker/setup"
	. "gitlab.kilic.dev/libraries/plumber/v4"
)

func DockerLogin(tl *TaskList[Pipe]) *Task[Pipe] {
	return tl.CreateTask("login").
		ShouldDisable(func(t *Task[Pipe]) bool {
			return t.Pipe.DockerRegistry.Username == "" ||
				t.Pipe.DockerRegistry.Password == ""
		}).
		Set(func(t *Task[Pipe]) error {
			t.Log.Infof(
				"Logging in to Docker registry: %s",
				t.Pipe.DockerRegistry.Registry,
			)

			result, err := setup.TL.Pipe.Ctx.Client.RegistryLogin(context.Background(), types.AuthConfig{
				ServerAddress: t.Pipe.DockerRegistry.Registry,
				Username:      t.Pipe.DockerRegistry.Username,
				Password:      t.Pipe.DockerRegistry.Password,
			})

			t.Log.Debugf("Result from Docker client: %+v", result)

			if err != nil {
				return err
			}

			return nil
		}).
		ShouldRunAfter(func(t *Task[Pipe]) error {
			return t.RunCommandJobAsJobSequence()
		})
}
