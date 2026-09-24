// Package client authenticates as a GitHub App and talks to the GitHub API with the
// installation token that buys.
package client

import (
	"context"
	"os"
	"strings"
	"time"

	. "github.com/cenk1cenk2/plumber/v7"
)

// Config is the GitHub App a pipe was configured with. The credentials come as a
// set, so a partial one fails validation instead of silently falling back.
type Config struct {
	Id             string `validate:"required_with=InstallationId PrivateKey"`
	InstallationId string `validate:"required_with=Id PrivateKey"`
	PrivateKey     string `validate:"required_with=Id InstallationId"`
	ApiUrl         string `validate:"required,url"`
}

// Enabled tells whether the app credentials were given at all.
func (c *Config) Enabled() bool {
	return c.Id != "" || c.InstallationId != "" || c.PrivateKey != ""
}

// Mint exchanges the app credentials for an installation token. The private key,
// the JWT and the token are all masked before anything can log them.
func (c *Config) Mint(ctx context.Context, p *Plumber, client ApplicationClientAdapter, request TokenRequest) (string, error) {
	pem, err := c.readPrivateKey()

	if err != nil {
		return "", err
	}

	// the logger re-indents every line of a multi-line message, so the key as a
	// whole never matches and only stays masked line by line.
	for line := range strings.Lines(pem) {
		p.AppendSecrets(strings.TrimSpace(line))
	}

	key, err := ParsePrivateKey([]byte(pem))

	if err != nil {
		return "", err
	}

	jwt, err := SignJWT(key, c.Id, time.Now())

	if err != nil {
		return "", err
	}

	p.AppendSecrets(jwt)

	token, err := client.CreateInstallationToken(ctx, jwt, c.InstallationId, request)

	if err != nil {
		return "", err
	}

	p.AppendSecrets(token)

	return token, nil
}

// readPrivateKey takes the key as a path first, since a GitLab file variable is
// the only way to keep a multi-line key masked in the job log.
func (c *Config) readPrivateKey() (string, error) {
	if strings.HasPrefix(strings.TrimSpace(c.PrivateKey), "-----BEGIN") {
		return c.PrivateKey, nil
	}

	content, err := os.ReadFile(c.PrivateKey)

	if err != nil {
		return "", err
	}

	return string(content), nil
}
