package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const SessionUUIDHeader = "X-Session-Uuid"

type User struct {
	UserUUID uuid.UUID `json:"user_uuid"`
	Login    string    `json:"login"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Whoami(ctx context.Context, sessionUUID string) (*User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/auth/whoami", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(SessionUUIDHeader, sessionUUID)
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("iam whoami: status %d: %s", res.StatusCode, string(body))
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

var ErrUnauthorized = fmt.Errorf("unauthorized")
