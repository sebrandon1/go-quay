/*
Package lib provides Quay.io API client functionality.

This file covers HTTP CLIENT and helper methods:

Client Setup:
  - NewClient(bearerToken string) (*Client, error)                - Create authenticated client (default URL)
  - NewClientWithURL(bearerToken, baseURL string) (*Client, error) - Create authenticated client with custom URL

HTTP Helper Methods:
  - do(req, v, acceptedStatuses...)          - Core HTTP executor (auth, headers, status check, decode)
  - get(req, v) / post(req, v) / put(req, v) / delete(req) - Thin wrappers with preset accepted statuses

Request Helpers:
  - newRequest(method, url string, body io.Reader) (*http.Request, error)
  - newRequestWithBody(method, url string, body any) (*http.Request, error)
  - decodeJSON(r io.Reader, v any) error

All HTTP methods include:
  - Bearer token authentication
  - Proper headers (Content-Type, Authorization)
  - Error handling for non-2xx responses
  - JSON marshaling/unmarshaling
*/
package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const maxErrorBodySize = 1 << 20 // 1 MB
const queryValueTrue = "true"

// DefaultQuayURL is the default Quay.io API base URL.
const DefaultQuayURL = "https://quay.io/api/v1"

type Client struct {
	BearerToken string
	BaseURL     string
	HTTPClient  *http.Client
}

func NewClientWithURL(bearerToken, baseURL string) (*Client, error) {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &Client{
		BearerToken: bearerToken,
		BaseURL:     baseURL,
		HTTPClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}, nil
}

func NewClient(bearerToken string) (*Client, error) {
	return NewClientWithURL(bearerToken, DefaultQuayURL)
}

func (c *Client) buildURL(pathFmt string, args ...any) string {
	escaped := make([]any, len(args))
	for i, a := range args {
		if s, ok := a.(string); ok {
			escaped[i] = url.PathEscape(s)
		} else {
			escaped[i] = a
		}
	}
	return c.BaseURL + fmt.Sprintf(pathFmt, escaped...)
}

func (c *Client) do(req *http.Request, v any, acceptedStatuses ...int) error {
	if c.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.BearerToken)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	accepted := false
	for _, s := range acceptedStatuses {
		if resp.StatusCode == s {
			accepted = true
			break
		}
	}
	if !accepted {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodySize))
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	if v != nil && resp.StatusCode != http.StatusNoContent {
		return decodeJSON(resp.Body, v)
	}

	return nil
}

func (c *Client) get(req *http.Request, v any) error {
	return c.do(req, v, http.StatusOK)
}

func (c *Client) post(req *http.Request, v any) error {
	return c.do(req, v, http.StatusOK, http.StatusCreated)
}

func (c *Client) put(req *http.Request, v any) error {
	return c.do(req, v, http.StatusOK, http.StatusCreated, http.StatusNoContent)
}

func (c *Client) delete(req *http.Request) error {
	return c.do(req, nil, http.StatusOK, http.StatusNoContent)
}

func addQueryParams(req *http.Request, params map[string]string) {
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
}

func decodeJSON(r io.Reader, v any) error {
	return json.NewDecoder(r).Decode(v)
}

//nolint:unparam
func newRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

// newRequestWithBody creates a new HTTP request with JSON body
func newRequestWithBody(method, url string, body any) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}
	return http.NewRequest(method, url, bodyReader)
}
