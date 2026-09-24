package client

import (
	"bytes"
	"context"
	json "encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	JSONRequest = "application/json"
	JSONAccept  = "application/vnd.github+json"
	APIVersion  = "2022-11-28"
)

type (
	// TokenRequest narrows an installation token down to the given repository
	// names and permissions; left empty, the token carries everything the
	// installation was granted.
	TokenRequest struct {
		Repositories []string          `json:"repositories,omitempty"`
		Permissions  map[string]string `json:"permissions,omitempty"`
	}

	CommitStatus struct {
		State       string `json:"state"`
		TargetUrl   string `json:"target_url,omitempty"`
		Description string `json:"description,omitempty"`
		Context     string `json:"context"`
	}

	ApplicationClientAdapter interface {
		CreateInstallationToken(ctx context.Context, jwt, installation string, request TokenRequest) (string, error)
		CreateCommitStatus(ctx context.Context, token, repository, sha string, status CommitStatus) error
	}
)

type (
	tokenResponse struct {
		Token string `json:"token"`
	}

	errorResponse struct {
		Message string `json:"message"`
	}
)

type applicationClient struct {
	address   string
	userAgent string
	client    *http.Client
}

var _ ApplicationClientAdapter = (*applicationClient)(nil)

func NewApplicationClient(address, userAgent string) ApplicationClientAdapter {
	return &applicationClient{
		address:   strings.TrimSuffix(address, "/"),
		userAgent: userAgent,
		client:    &http.Client{},
	}
}

func (c *applicationClient) CreateInstallationToken(ctx context.Context, jwt, installation string, request TokenRequest) (string, error) {
	body, err := c.do(
		ctx,
		fmt.Sprintf("%s/app/installations/%s/access_tokens", c.address, installation),
		jwt,
		request,
	)

	if err != nil {
		return "", fmt.Errorf("Can not create installation token: %w", err)
	}

	response := tokenResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("Can not create installation token, response unexpected: %w", err)
	}

	if response.Token == "" {
		return "", fmt.Errorf("Can not create installation token, response carries no token.")
	}

	return response.Token, nil
}

func (c *applicationClient) CreateCommitStatus(ctx context.Context, token, repository, sha string, status CommitStatus) error {
	if _, err := c.do(
		ctx,
		fmt.Sprintf("%s/repos/%s/statuses/%s", c.address, repository, sha),
		token,
		status,
	); err != nil {
		return fmt.Errorf("Can not create commit status: %w", err)
	}

	return nil
}

// do posts the payload and hands back the body of a successful response. A failed
// one only surfaces GitHub's message, since the rest of the exchange carries the
// credentials.
func (c *applicationClient) do(ctx context.Context, address, token string, payload any) ([]byte, error) {
	body, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, address, bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", JSONRequest)
	req.Header.Set("Accept", JSONAccept)
	req.Header.Set("X-GitHub-Api-Version", APIVersion)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusCreated {
		response := errorResponse{}
		_ = json.Unmarshal(body, &response)

		return nil, fmt.Errorf("GitHub responded with code: %d > %s", res.StatusCode, response.Message)
	}

	return body, nil
}
