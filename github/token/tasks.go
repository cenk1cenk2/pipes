package token

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"regexp"
	"slices"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
	"github.com/joho/godotenv"

	"gitlab.kilic.dev/devops/pipes/github/client"
)

func mint(tl *TaskList) *Task {
	return tl.CreateTask("mint").
		Set(func(ctx context.Context, t *Task) error {
			t.Log.Info(fmt.Sprintf("Minting installation token for app: %s > %s", P.App.Id, P.App.InstallationId))

			token, err := P.App.Mint(ctx, t.Plumber, C.Client, client.TokenRequest{
				Repositories: P.Token.Repositories,
				Permissions:  P.Token.Permissions,
			})

			if err != nil {
				return err
			}

			C.Token = token

			return nil
		})
}

func write(tl *TaskList) *Task {
	return tl.CreateTask("write").
		Set(func(_ context.Context, t *Task) error {
			env, err := godotenv.Read(P.Token.File)

			if errors.Is(err, fs.ErrNotExist) {
				env = map[string]string{}
			} else if err != nil {
				return err
			}

			if P.Token.Variable != "" {
				env[P.Token.Variable] = C.Token
			}

			if P.Token.GitCredentialsVariable != "" {
				credentials := fmt.Sprintf("x-access-token:%s", C.Token)
				t.Plumber.AppendSecrets(credentials)

				env[P.Token.GitCredentialsVariable] = credentials
			}

			content, err := marshal(env)

			if err != nil {
				return err
			}

			f, err := os.OpenFile(P.Token.File, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)

			if err != nil {
				return err
			}

			defer f.Close()

			if _, err := f.WriteString(content); err != nil {
				return err
			}

			if P.Token.Variable != "" {
				t.Log.Info(fmt.Sprintf("Token written to dotenv file: %s > %s", P.Token.File, P.Token.Variable))
			}

			if P.Token.GitCredentialsVariable != "" {
				t.Log.Info(fmt.Sprintf("Git credential written to dotenv file: %s > %s", P.Token.File, P.Token.GitCredentialsVariable))
			}

			return nil
		})
}

// marshal writes the variables unquoted, since GitLab keeps the quotes of a
// dotenv value as part of it.
func marshal(env map[string]string) (string, error) {
	lines := make([]string, 0, len(env))
	// GitLab only takes dotenv keys made of letters, digits and underscores.
	variable := regexp.MustCompile(`^[A-Za-z0-9_]+$`)

	for _, key := range slices.Sorted(maps.Keys(env)) {
		if !variable.MatchString(key) {
			return "", fmt.Errorf("Variable name can only contain letters, digits and underscores: %s", key)
		}

		if strings.ContainsAny(env[key], "\r\n\"'`#") {
			return "", fmt.Errorf("Value of variable can not carry newlines, quotes or comments in a dotenv file: %s", key)
		}

		lines = append(lines, fmt.Sprintf("%s=%s\n", key, env[key]))
	}

	return strings.Join(lines, ""), nil
}
