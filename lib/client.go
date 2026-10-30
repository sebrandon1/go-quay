/*
Package lib provides Quay.io API client functionality.

This file covers HTTP CLIENT and helper methods:

Client Setup:
  - NewClient(bearerToken string) (*Client, error)                - Create authenticated client (default URL)
  - NewClientWithURL(bearerToken, baseURL string) (*Client, error) - Create authenticated client with custom URL

HTTP Helper Methods:
  - do(req, v, acceptedStatuses...)          - Core HTTP executor (auth, headers, status check, decode)
  - get(req, v) / post(req, v) / put(req, v) / delete(req) - Thin wrappers with preset accepted statuses

Request Helpers (thin wrappers over go-toolbox/httpx):
  - newRequest(ctx, method, url string) (*http.Request, error)
  - newRequestWithBody(ctx, method, url string, body any) (*http.Request, error)

All HTTP methods include:
  - Bearer token authentication
  - Proper headers (Content-Type, Authorization)
  - Error handling for non-2xx responses
  - JSON marshaling/unmarshaling
*/
package lib

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/sebrandon1/go-toolbox/httpx"
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

// RequestLogEntry describes one HTTP attempt made by Client.
type RequestLogEntry struct {
	Method     string
	URL        string
	Attempt    int
	StatusCode int
	Duration   time.Duration
}

type Client struct {
	BearerToken   string
	BaseURL       string
	Version       string
	Retry         *RetryConfig
	HTTPClient    *http.Client
	RequestLogger func(RequestLogEntry)
}

var _ RepositoryReader = (*Client)(nil)

func NewClientWithURL(bearerToken, baseURL string) (*Client, error) {
	return &Client{
		BearerToken: bearerToken,
		BaseURL:     baseURL,
		Version:     "dev",
		HTTPClient: httpx.NewHTTPClient(httpx.Options{
			Timeout:             30 * time.Second,
			DialTimeout:         10 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		}),
	}, nil
}

func NewClient(bearerToken string) (*Client, error) {
	return NewClientWithURL(bearerToken, DefaultQuayURL)
}

func (c *Client) buildURL(pathFmt string, args ...any) string {
	escaped := make([]any, len(args))
	for i, a := range args {
		if s, ok := a.(string); ok {
			escaped[i] = httpx.PathSeg(s)
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
		result, err := c.doOnce(req, v, acceptedStatuses, attempt+1)
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

func (c *Client) doOnce(req *http.Request, v any, acceptedStatuses []int, attempt int) (error, error) {
	start := time.Now()
	entry := RequestLogEntry{
		Method:  req.Method,
		URL:     redactRequestURL(req.URL),
		Attempt: attempt,
	}
	defer func() {
		entry.Duration = time.Since(start)
		if c.RequestLogger != nil {
			c.RequestLogger(entry)
		}
	}()

	resp, err := c.HTTPClient.Do(req)
	if resp != nil {
		entry.StatusCode = resp.StatusCode
	}
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, &retryableError{err: err}
	}
	defer resp.Body.Close()

	if isAccepted(resp.StatusCode, acceptedStatuses) {
		if v != nil && resp.StatusCode != http.StatusNoContent {
			return httpx.DecodeJSON(resp.Body, v), nil
		}
		return nil, nil
	}

	body := httpx.ReadErrorBody(resp.Body, maxErrorBodySize)
	apiErr := buildAPIError(resp.StatusCode, body)

	if isRetryableStatus(resp.StatusCode) {
		return nil, &retryableError{
			err:        apiErr,
			retryAfter: parseRetryAfter(resp.Header.Get("Retry-After")),
		}
	}

	return nil, apiErr
}

func redactRequestURL(requestURL *url.URL) string {
	return httpx.RedactURLQuery(requestURL)
}

func isAccepted(status int, accepted []int) bool {
	for _, s := range accepted {
		if status == s {
			return true
		}
	}
	return false
}

// buildAPIError converts a non-2xx response into an error.
//
// Quay uses two error body shapes: some endpoints return {"error": ...}, while
// the authentication and not-found responses return
// {"detail": ..., "error_message": ..., "error_type": ...}. Both become a
// *QuayError so callers can branch on the HTTP status instead of parsing text.
func buildAPIError(status int, body string) error {
	var quayErr QuayError
	if err := json.Unmarshal([]byte(body), &quayErr); err != nil {
		return fmt.Errorf("unexpected status code: %d, response: %s", status, body)
	}

	quayErr.Status = status
	if quayErr.Message == "" {
		quayErr.Message = cmp.Or(
			quayErr.ErrorMessage, quayErr.Detail, quayErr.ErrorType, http.StatusText(status))
	}

	return &quayErr
}

// StatusCode returns the HTTP status code carried by err. It reports false when
// err did not come from a Quay API response.
func StatusCode(err error) (int, bool) {
	var quayErr *QuayError
	if errors.As(err, &quayErr) {
		return quayErr.StatusCode(), true
	}

	return 0, false
}

// statusIs reports whether err carries one of the given HTTP statuses.
func statusIs(err error, statuses ...int) bool {
	code, ok := StatusCode(err)
	if !ok {
		return false
	}

	for _, status := range statuses {
		if code == status {
			return true
		}
	}

	return false
}

// IsNotFound reports whether err is a 404 response.
func IsNotFound(err error) bool {
	return statusIs(err, http.StatusNotFound)
}

// IsUnauthorized reports whether err is a 401 response.
//
// Quay answers 401 for repositories the caller may not see, which includes
// repositories that no longer exist: it deliberately does not distinguish
// "missing" from "private" in order to avoid leaking repository names.
func IsUnauthorized(err error) bool {
	return statusIs(err, http.StatusUnauthorized)
}

// IsForbidden reports whether err is a 403 response.
func IsForbidden(err error) bool {
	return statusIs(err, http.StatusForbidden)
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
	return int(httpx.ParseRetryAfter(header) / time.Second)
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
	httpx.AddQueryParams(req, params)
}

func newRequest(ctx context.Context, method, requestURL string) (*http.Request, error) {
	return httpx.NewRequest(ctx, method, requestURL)
}

// newRequestWithBody creates a new HTTP request with JSON body
func newRequestWithBody(ctx context.Context, method, requestURL string, body any) (*http.Request, error) {
	return httpx.NewRequestWithBody(ctx, method, requestURL, body)
}
