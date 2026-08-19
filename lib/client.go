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
  - newRequest(ctx, method, url string, body io.Reader) (*http.Request, error)
  - newRequestWithBody(ctx, method, url string, body any) (*http.Request, error)
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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const maxErrorBodySize = 1 << 20 // 1 MB
const queryValueTrue = "true"

// DefaultQuayURL is the default Quay.io API base URL.
const DefaultQuayURL = "https://quay.io/api/v1"

// RetryConfig controls automatic retry behavior for failed HTTP requests.
type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

type Client struct {
	BearerToken string
	BaseURL     string
	Version     string
	Retry       *RetryConfig
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
		Version:     "dev",
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
	c.setHeaders(req)

	maxAttempts := 1
	if c.Retry != nil && c.Retry.MaxRetries > 0 {
		maxAttempts = 1 + c.Retry.MaxRetries
	}

	var lastErr error
	for attempt := range maxAttempts {
		result, err := c.doOnce(req, v, acceptedStatuses)
		if err == nil {
			return result
		}

		lastErr = err
		if !c.shouldRetry(attempt, maxAttempts) {
			return err
		}

		retryable, retryAfter := isRetryableError(err)
		if !retryable {
			return err
		}

		if err := c.backoff(req.Context(), attempt, retryAfter); err != nil {
			return err
		}
	}

	return lastErr
}

func (c *Client) setHeaders(req *http.Request) {
	if c.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.BearerToken)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "go-quay/"+c.Version)
}

type retryableError struct {
	err        error
	retryAfter int
}

func (e *retryableError) Error() string { return e.err.Error() }
func (e *retryableError) Unwrap() error { return e.err }

func (c *Client) doOnce(req *http.Request, v any, acceptedStatuses []int) (error, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, &retryableError{err: err}
	}
	defer resp.Body.Close()

	if isAccepted(resp.StatusCode, acceptedStatuses) {
		if v != nil && resp.StatusCode != http.StatusNoContent {
			return decodeJSON(resp.Body, v), nil
		}
		return nil, nil
	}

	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodySize))
	apiErr := buildAPIError(resp.StatusCode, body)

	if isRetryableStatus(resp.StatusCode) {
		return nil, &retryableError{
			err:        apiErr,
			retryAfter: parseRetryAfter(resp.Header.Get("Retry-After")),
		}
	}

	return nil, apiErr
}

func isAccepted(status int, accepted []int) bool {
	for _, s := range accepted {
		if status == s {
			return true
		}
	}
	return false
}

func buildAPIError(status int, body []byte) error {
	var quayErr QuayError
	if json.Unmarshal(body, &quayErr) == nil && quayErr.Message != "" {
		quayErr.Status = status
		return &quayErr
	}
	return fmt.Errorf("unexpected status code: %d, response: %s", status, string(body))
}

func isRetryableError(err error) (bool, int) {
	var re *retryableError
	if errors.As(err, &re) {
		return true, re.retryAfter
	}
	return false, 0
}

func (c *Client) shouldRetry(attempt, maxAttempts int) bool {
	return c.Retry != nil && attempt < maxAttempts-1
}

func (c *Client) backoff(ctx context.Context, attempt int, retryAfterSecs int) error {
	var wait time.Duration
	if retryAfterSecs > 0 {
		wait = time.Duration(retryAfterSecs) * time.Second
	} else {
		wait = c.Retry.InitialBackoff * time.Duration(math.Pow(2, float64(attempt)))
	}

	if c.Retry.MaxBackoff > 0 && wait > c.Retry.MaxBackoff {
		wait = c.Retry.MaxBackoff
	}

	t := time.NewTimer(wait)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func isRetryableStatus(status int) bool {
	return status == http.StatusTooManyRequests || status >= http.StatusInternalServerError
}

func parseRetryAfter(header string) int {
	if header == "" {
		return 0
	}
	seconds, err := strconv.Atoi(header)
	if err != nil {
		return 0
	}
	return seconds
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

func newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

// newRequestWithBody creates a new HTTP request with JSON body
func newRequestWithBody(ctx context.Context, method, url string, body any) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}
	return newRequest(ctx, method, url, bodyReader)
}
