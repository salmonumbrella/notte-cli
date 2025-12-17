package api

import (
	"context"
	"fmt"
	"net/http"
)

const DefaultBaseURL = "https://api.notte.cc"

// NotteClient wraps the generated client with auth
type NotteClient struct {
	client  *ClientWithResponses
	baseURL string
	apiKey  string
}

// NewClient creates a new Notte API client
func NewClient(apiKey string) (*NotteClient, error) {
	return NewClientWithURL(apiKey, DefaultBaseURL)
}

// NewClientWithURL creates a client with custom base URL
func NewClientWithURL(apiKey, baseURL string) (*NotteClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// Create HTTP client with auth header
	httpClient := &http.Client{
		Transport: &authTransport{
			apiKey: apiKey,
			base:   http.DefaultTransport,
		},
	}

	client, err := NewClientWithResponses(baseURL, WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &NotteClient{
		client:  client,
		baseURL: baseURL,
		apiKey:  apiKey,
	}, nil
}

// authTransport adds Authorization header to all requests
type authTransport struct {
	apiKey string
	base   http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.apiKey)
	return t.base.RoundTrip(req)
}

// Client returns the underlying generated client for direct access
func (c *NotteClient) Client() *ClientWithResponses {
	return c.client
}

// Context helper for commands
func DefaultContext() context.Context {
	return context.Background()
}
