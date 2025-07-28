/*
Package lib provides Quay.io API client functionality.

This file covers HTTP CLIENT and helper methods:

Client Setup:
  - NewClient(bearerToken string) (*Client, error)   - Create authenticated client

HTTP Helper Methods:
  - get(req *http.Request, v interface{}) error      - Execute GET requests
  - post(req *http.Request, v interface{}) error     - Execute POST requests
  - put(req *http.Request, v interface{}) error      - Execute PUT requests
  - delete(req *http.Request) error                  - Execute DELETE requests

Request Helpers:
  - newRequest(method, url string, body io.Reader) (*http.Request, error)
  - newRequestWithBody(method, url string, body interface{}) (*http.Request, error)
  - mustMarshal(v interface{}) []byte
  - decodeJSON(r io.Reader, v interface{}) error

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

	// Create a custom transport to force IPv4 connections
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
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

func (c *Client) get(req *http.Request, v interface{}) error {
	req.Header.Set("Authorization", "Bearer "+c.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("error closing response body:", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	return decodeJSON(resp.Body, v)
}

func (c *Client) post(req *http.Request, v interface{}) error {
	req.Header.Set("Authorization", "Bearer "+c.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("error closing response body:", err)
		}
	}()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	if v != nil {
		return decodeJSON(resp.Body, v)
	}

	return nil
}

func (c *Client) put(req *http.Request, v interface{}) error {
	req.Header.Set("Authorization", "Bearer "+c.BearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("error closing response body:", err)
		}
	}()

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

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("error closing response body:", err)
		}
	}()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(body))
	}

	return nil
}

func decodeJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

//nolint:unparam
func newRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}

// mustMarshal marshals a value to JSON and panics on error
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal JSON: %v", err))
	}
	return data
}

// newRequestWithBody creates a new HTTP request with JSON body
func newRequestWithBody(method, url string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(mustMarshal(body))
	}
	return http.NewRequest(method, url, bodyReader)
}
