package git

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	. "github.com/cenk1cenk2/plumber/v7"
)

const DefaultUsername = "oauth2"

// Repository is a working tree the write-back acts in. Every git call runs through
// it, so a lazily fetched blob of a partial clone authenticates like the clone did.
// Username pairs with the token and falls back to DefaultUsername.
type Repository struct {
	Dir      string
	Token    string
	Username string
}

// Clone fetches only the tip of the branch and only the destinations, so a large
// target costs about as much as the paths the sync writes into. The sparse patterns
// are not in cone mode, which only takes directories, and every destination is
// anchored at the root so it can not match a same-named path further down.
func (r Repository) Clone(ctx context.Context, t *Task, url string, branch string, destinations []string) error {
	if err := ensureNoCredentials(url); err != nil {
		return err
	}

	patterns := []string{}
	for _, destination := range destinations {
		patterns = append(patterns, fmt.Sprintf("/%s", filepath.ToSlash(filepath.Clean(destination))))
	}

	if err := r.command(t, "", "clone", "--depth", "1", "--single-branch", "--branch", branch, "--filter=blob:none", "--sparse", url, r.Dir).
		Run(ctx); err != nil {
		return fmt.Errorf("Can not clone the target branch: %s -> %w", branch, err)
	}

	if err := r.command(t, r.Dir, append([]string{"sparse-checkout", "set", "--no-cone", "--"}, patterns...)...).
		Run(ctx); err != nil {
		return fmt.Errorf("Can not narrow the clone to the destinations: %v -> %w", destinations, err)
	}

	return nil
}

// Worktree checks the branch out next to the checkout the pipeline already has. The
// worktree gets an index of its own, so the checkout the job runs in stays untouched.
func (r Repository) Worktree(ctx context.Context, t *Task, remote string, branch string, path string) (Repository, error) {
	if err := r.command(t, r.Dir, "fetch", "--depth", "1", remote, fmt.Sprintf("refs/heads/%s", branch)).
		Run(ctx); err != nil {
		return Repository{}, fmt.Errorf("Can not fetch the target branch: %s -> %w", branch, err)
	}

	if err := r.command(t, r.Dir, "worktree", "add", "--detach", path, "FETCH_HEAD").
		Run(ctx); err != nil {
		return Repository{}, fmt.Errorf("Can not add a worktree for the target branch: %s -> %w", branch, err)
	}

	return Repository{Dir: path, Token: r.Token}, nil
}

// Lease is the commit the branch points at on the remote, or empty when it does not
// exist yet, which a push leased on it reads as "must still not exist".
func (r Repository) Lease(ctx context.Context, t *Task, remote string, branch string) (string, error) {
	if err := ensureNoCredentials(remote); err != nil {
		return "", err
	}

	var output string

	if err := r.command(t, r.Dir, "ls-remote", "--heads", remote, fmt.Sprintf("refs/heads/%s", branch)).
		SetLogLevel(LogLevelDebug, LogLevelWarn, LogLevelDebug).
		CaptureStdout(&output).
		Run(ctx); err != nil {
		return "", fmt.Errorf("Can not look up the branch on the remote: %s -> %w", branch, err)
	}

	sha, _, _ := strings.Cut(strings.TrimSpace(output), "\t")

	return sha, nil
}

// Push overwrites the branch only while it still points at the lease, so a run that
// raced another one fails instead of discarding what the other pushed.
func (r Repository) Push(ctx context.Context, t *Task, remote string, branch string, lease string) error {
	if err := ensureNoCredentials(remote); err != nil {
		return err
	}

	ref := fmt.Sprintf("refs/heads/%s", branch)

	if err := r.command(t, r.Dir, "push", fmt.Sprintf("--force-with-lease=%s:%s", ref, lease), remote, fmt.Sprintf("HEAD:%s", ref)).
		Run(ctx); err != nil {
		return fmt.Errorf("Can not push the branch: %s -> %w", branch, err)
	}

	return nil
}

// Stage forces the destinations into the index, since a destination the target
// ignores would otherwise stage nothing and come out as an empty diff without an
// error. A sparse clone refuses paths it considers outside its definition otherwise.
// A destination neither on disk nor in the index, like one an absent source left
// empty in a target that never had it, is a pathspec git refuses, so it is left out.
func (r Repository) Stage(ctx context.Context, t *Task, destinations []string) error {
	existing := []string{}

	for _, destination := range destinations {
		exists, err := r.exists(ctx, t, destination)
		if err != nil {
			return err
		}

		if !exists {
			t.Log.Warn(fmt.Sprintf("Destination does not exist, leaving it out of the commit: %s", destination))

			continue
		}

		existing = append(existing, destination)
	}

	if len(existing) == 0 {
		return fmt.Errorf("None of the destinations exist, there is nothing to stage: %v", destinations)
	}

	if err := r.command(t, r.Dir, append([]string{"add", "--sparse", "--all", "--force", "--"}, existing...)...).
		Run(ctx); err != nil {
		return fmt.Errorf("Can not stage the destinations: %v -> %w", existing, err)
	}

	return nil
}

// A destination the sync removed from disk still exists while the index tracks it,
// since staging it is what records the deletion.
func (r Repository) exists(ctx context.Context, t *Task, destination string) (bool, error) {
	if _, err := os.Lstat(filepath.Join(r.Dir, destination)); err == nil {
		return true, nil
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("Can not check the destination: %s -> %w", destination, err)
	}

	var output string

	if err := r.command(t, r.Dir, "ls-files", "--", destination).
		SetLogLevel(LogLevelDebug, LogLevelWarn, LogLevelDebug).
		CaptureStdout(&output).
		Run(ctx); err != nil {
		return false, fmt.Errorf("Can not look up the destination in the index: %s -> %w", destination, err)
	}

	return strings.TrimSpace(output) != "", nil
}

// Patch writes the staged changes to path; an empty file means there is nothing to
// publish. A relative path resolves against the working directory of the job, not
// the repository the diff runs in.
func (r Repository) Patch(ctx context.Context, t *Task, path string) error {
	output, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("Can not resolve the patch path: %s -> %w", path, err)
	}

	if err := r.command(t, r.Dir, "diff", "--cached", "--binary", fmt.Sprintf("--output=%s", output)).
		Run(ctx); err != nil {
		return fmt.Errorf("Can not write the staged changes: %s -> %w", path, err)
	}

	return nil
}

// The token travels as an extra header through the environment, since an argument or
// a URL shows up in the process listing, the command log and the remote's config.
func (r Repository) command(t *Task, dir string, args ...string) *Command {
	c := t.CreateCommand("git", args...).
		SetDir(dir)

	if r.Token == "" {
		return c
	}

	username := r.Username
	if username == "" {
		username = DefaultUsername
	}

	// the inherited environment is appended after the one set here and wins on a
	// duplicate key, which would drop the header without git ever complaining.
	return c.
		ShouldRunBefore(func(_ context.Context, _ *Command) error {
			if _, ok := os.LookupEnv("GIT_CONFIG_COUNT"); ok {
				return fmt.Errorf("GIT_CONFIG_COUNT is already set in the environment, which would override the git credentials.")
			}

			return nil
		}).
		AppendEnvironment(map[string]string{
			"GIT_CONFIG_COUNT":   "1",
			"GIT_CONFIG_KEY_0":   "http.extraHeader",
			"GIT_CONFIG_VALUE_0": fmt.Sprintf("Authorization: Basic %s", base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "%s:%s", username, r.Token))),
		})
}

// The remote ends up in the command log, so one carrying credentials of its own, like
// CI_REPOSITORY_URL does, is refused instead of printed.
func ensureNoCredentials(remote string) error {
	u, err := url.Parse(remote)
	if err != nil {
		return fmt.Errorf("Can not parse the remote URL.")
	}

	if u.User != nil {
		return fmt.Errorf("Remote URL can not carry credentials, pass the token instead.")
	}

	return nil
}
