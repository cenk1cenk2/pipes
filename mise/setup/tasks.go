package setup

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
)

// ParseVersion takes the version out of the banner of mise --version, which
// trails it with the platform and the build date.
func ParseVersion(output string) string {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return ""
	}

	return fields[0]
}

func initialize(tl *TaskList) *Task {
	return tl.CreateTask("init").
		Set(func(_ context.Context, t *Task) error {
			C.Cwd = P.Cwd

			t.Log.Debug(fmt.Sprintf("Working directory: %s", C.Cwd))

			executable, err := exec.LookPath("mise")
			if err != nil {
				return fmt.Errorf("Can not find the mise binary of the image: %w", err)
			}

			C.Executable = executable

			t.Log.Debug(fmt.Sprintf("mise binary: %s", C.Executable))

			return nil
		})
}

func version(tl *TaskList) *Task {
	return tl.CreateTask("version").
		Set(func(_ context.Context, t *Task) error {
			var output string

			t.CreateCommand(C.Executable, "--version").
				SetLogLevel(LogLevelDebug, LogLevelDebug, LogLevelDebug).
				CaptureStdout(&output).
				ShouldRunAfter(func(_ context.Context, c *Command) error {
					C.Version = ParseVersion(output)

					if C.Version == "" {
						return fmt.Errorf("Can not resolve the version of mise: %s", C.Executable)
					}

					c.Log.Info(fmt.Sprintf("mise version: %s", C.Version))

					return nil
				}).
				AddSelfToTheTask()

			return nil
		}).
		ShouldRunAfter(func(ctx context.Context, t *Task) error {
			return t.RunCommandJobAsJobSequence(ctx)
		})
}

// The commands run with the environment of the process masked and receive a full
// copy of it instead, since plumber appends that environment after the one a
// command sets and the last value of a key wins.
func environment(tl *TaskList) *Task {
	return tl.CreateTask("environment").
		Set(func(_ context.Context, t *Task) error {
			dataDir, err := filepath.Abs(P.DataDir)
			if err != nil {
				return fmt.Errorf("Cannot get absolute path of data dir: %s -> %w", P.DataDir, err)
			}

			C.DataDir = dataDir
			C.Binary = filepath.Join(dataDir, "bin", "mise")

			C.Env = ParseEnvironmentVariablesToMap()
			C.Env["MISE_DATA_DIR"] = C.DataDir

			// mise links the shims to the first mise on the path, not to the one that runs.
			path := filepath.Dir(C.Binary)
			if C.Env["PATH"] != "" {
				path = strings.Join([]string{path, C.Env["PATH"]}, string(os.PathListSeparator))
			}

			C.Env["PATH"] = path

			t.Log.Debug(fmt.Sprintf("Data directory: %s", C.DataDir))

			return nil
		})
}
