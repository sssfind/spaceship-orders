package env

import (
	"errors"
	"os"
	"strings"
)

type iamHTTP struct {
	baseURL string
}

func NewIAMHTTPConfig() (*iamHTTP, error) {
	baseURL := strings.TrimRight(os.Getenv("IAM_HTTP_BASE_URL"), "/")
	if baseURL == "" {
		return nil, errors.New("IAM_HTTP_BASE_URL is not set")
	}
	return &iamHTTP{baseURL: baseURL}, nil
}

func (cfg *iamHTTP) BaseURL() string { return cfg.baseURL }
