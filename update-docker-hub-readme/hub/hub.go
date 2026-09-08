// Package hub is the Docker Hub API narrowed to the two calls this pipe makes,
// so the update can be driven without a registry to talk to.
package hub

import (
	"bytes"
	"context"
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"fmt"
	"io"
	"net/http"
)

const JSONRequest = "application/json"

// The login endpoint belongs to the account API, while the configured address
// points at the repositories the readme is pushed to.
const LoginAddress = "https://hub.docker.com/v2/users/login/"

type (
	// Readme is the pair of descriptions a repository page shows.
	Readme struct {
		Description string
		Full        string
	}

	// Result is the update response narrowed to what decides whether the readme landed.
	Result struct {
		CanEdit         bool
		Description     string
		FullDescription string
		StatusCode      int
	}

	ClientAdapter interface {
		Login(ctx context.Context, username, password string) (string, error)
		UpdateReadme(ctx context.Context, token, repository string, readme Readme) (Result, error)
	}
)

type (
	credentials struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}

	loginResponse struct {
		Token string `json:"token"`
	}

	updateRequest struct {
		Description string `json:"description"`
		Readme      string `json:"full_description"`
	}

	updateResponse struct {
		CanEdit         bool   `json:"can_edit"`
		Description     string `json:"description"`
		FullDescription string `json:"full_description"`
	}
)

type client struct {
	address      string
	loginAddress string
	userAgent    string
	client       *http.Client
}

var _ ClientAdapter = (*client)(nil)

func NewClient(address, userAgent string) ClientAdapter {
	return &client{
		address:      address,
		loginAddress: LoginAddress,
		userAgent:    userAgent,
		client:       &http.Client{},
	}
}

func (c *client) Login(ctx context.Context, username, password string) (string, error) {
	body, err := json.Marshal(credentials{
		Username: username,
		Password: password,
	}, jsontext.EscapeForHTML(true))

	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.loginAddress,
		bytes.NewReader(body),
	)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", JSONRequest)

	res, err := c.client.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	response := loginResponse{}
	if err := json.Unmarshal(body, &response, json.RejectUnknownMembers(true)); err != nil {
		return "", err
	}

	return response.Token, nil
}

func (c *client) UpdateReadme(
	ctx context.Context,
	token, repository string,
	readme Readme,
) (Result, error) {
	body, err := json.Marshal(updateRequest{
		Description: readme.Description,
		Readme:      readme.Full,
	}, jsontext.EscapeForHTML(true))

	if err != nil {
		return Result{}, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		fmt.Sprintf("%s/%s/", c.address, repository),
		bytes.NewReader(body),
	)

	if err != nil {
		return Result{}, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", JSONRequest)
	req.Header.Set("Authorization", fmt.Sprintf("JWT %s", token))

	res, err := c.client.Do(req)

	if err != nil {
		return Result{}, err
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)

	if err != nil {
		return Result{}, err
	}

	response := updateResponse{}
	// the repository payload carries far more than the fields above, so this decode
	// stays lenient where the login one rejects unknown members.
	if err := json.Unmarshal(body, &response); err != nil {
		return Result{}, fmt.Errorf("Response unexpected: %w > %s", err, string(body))
	}

	return Result{
		CanEdit:         response.CanEdit,
		Description:     response.Description,
		FullDescription: response.FullDescription,
		StatusCode:      res.StatusCode,
	}, nil
}
