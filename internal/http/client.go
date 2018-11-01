package http

import (
	"net/http"
)

type (
	// BasicAuth contains basic HTTP authentication credentials.
	BasicAuth struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// ClientWithBasicAuth implements client with basic authorization
	ClientWithBasicAuth struct {
		httpClient HTTPClient
		auth       BasicAuth
	}

	// HTTPClient is an interface for testing a request object.
	HTTPClient interface {
		Do(*http.Request) (*http.Response, error)
	}
)

// NewClientWithBasicAuth returns httpClient which always enriches request with basic auth
func NewClientWithBasicAuth(auth BasicAuth) *ClientWithBasicAuth {
	return &ClientWithBasicAuth{
		auth:       auth,
		httpClient: http.DefaultClient,
	}
}

// Do executes passed HTTP request and returns a response. It enriches request with basic auth.
func (c *ClientWithBasicAuth) Do(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(c.auth.Username, c.auth.Password)
	return c.httpClient.Do(req)
}

func (c *ClientWithBasicAuth) WithClient(client HTTPClient) *ClientWithBasicAuth {
	c.httpClient = client
	return c
}
