package registry

import (
	"io"
	"slices"
	"strings"

	"github.com/cenk1cenk2/plumber/v6"
	"gitlab.kilic.dev/devops/pipes/internal/cli"
)

// LoginStep is the login stage of a pipe: the registry flags it declares and the
// task list that authenticates with them. It is declared once per pipe, so every
// command that publishes composes the same step rather than its own copy of the
// flags.
func LoginStep(spec Spec, creds *Credentials, binary string, args ...string) cli.Step {
	return cli.Step{
		Flags: NewFlags(spec, creds),
		New: func(p *plumber.Plumber) *plumber.TaskList {
			return LoginTaskList(p, creds, binary, args...)
		},
	}
}

// LoginTaskList logs in to the registry with the given binary, which is invoked
// as "<binary> <args...> <uri> --username <username> --password-stdin".
func LoginTaskList(p *plumber.Plumber, creds *Credentials, binary string, args ...string) *plumber.TaskList {
	tl := &plumber.TaskList{}

	return tl.New(p).
		SetRuntimeDepth(3).
		ShouldRunBefore(func(_ *plumber.TaskList) error {
			return cli.Validated(p, creds)
		}).
		Set(func(tl *plumber.TaskList) plumber.Job {
			return plumber.JobSequence(
				LoginTask(tl, creds, binary, args...).Job(),
			)
		})
}

// LoginTask hands the password to the binary over stdin so it never reaches an
// argument list. It is disabled without both a username and a password, since
// those pipelines are relying on an ambient login instead.
func LoginTask(tl *plumber.TaskList, creds *Credentials, binary string, args ...string) *plumber.Task {
	return tl.CreateTask("login").
		ShouldDisable(func(_ *plumber.Task) bool {
			return creds.Username == "" || creds.Password == ""
		}).
		Set(func(t *plumber.Task) error {
			t.Log.Infof("Logging in to registry: %s", creds.Uri)

			t.CreateCommand(binary, LoginArgs(creds, args...)...).
				SetLogLevel(plumber.LOG_LEVEL_DEBUG, plumber.LOG_LEVEL_DEBUG, plumber.LOG_LEVEL_DEFAULT).
				SetStdin(func(_ *plumber.Command) io.Reader {
					return strings.NewReader(creds.Password)
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(t *plumber.Task) error {
			return t.RunCommandJobAsJobSequence()
		})
}

// LoginArgs is everything after the binary. The password is deliberately absent:
// an argument list shows up in a process listing and in the command trace log.
func LoginArgs(creds *Credentials, args ...string) []string {
	return slices.Concat(args, []string{creds.Uri, "--username", creds.Username, "--password-stdin"})
}
