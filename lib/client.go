/*
Package lib provides Quay.io API client functionality.

This file covers HTTP CLIENT and helper methods:

Client Setup:
  - NewClient(bearerToken string) (*Client, error)   - Create authenticated client

HTTP Helper Methods:
  - get(req *http.Request, v any) error      - Execute GET requests
  - post(req *http.Request, v any) error     - Execute POST requests
  - put(req *http.Request, v any) error      - Execute PUT requests
  - delete(req *http.Request) error                  - Execute DELETE requests

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
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type Client struct {
	BearerToken string
	HTTPClient  *http.Client
}

func NewClient(bearerToken string) (*Client, error) {
	if bearerToken == "" {
		return nil, errors.New("bearer token is required")
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &Client{
		BearerToken: bearerToken,
		HTTPClient: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}, nil
}

func (c *Client) get(req *http.Request, v any) error {
	req.Header.Set("Authorization", "Bearer "+c.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	return decodeJSON(resp.Body, v)
}

func (c *Client) post(req *http.Request, v any) error {
	req.Header.Set("Authorization", "Bearer "+c.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	if v != nil {
		return decodeJSON(resp.Body, v)
	}

	return nil
}

func (c *Client) put(req *http.Request, v any) error {
	req.Header.Set("Authorization", "Bearer "+c.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	if v != nil && resp.StatusCode != http.StatusNoContent {
		return decodeJSON(resp.Body, v)
	}

	return nil
}

func (c *Client) delete(req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+c.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
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
