package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	return &Client{
		BearerToken: bearerToken,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
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

func decodeJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

//nolint:unparam
func newRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}
