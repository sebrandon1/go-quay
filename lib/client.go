package lib

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Client struct {
	BearerToken string
}

func NewClient(bearerToken string) (*Client, error) {
	if bearerToken == "" {
		return nil, errors.New("bearer token is required")
	}

	return &Client{
		BearerToken: bearerToken,
	}, nil
}

func (c *Client) Get(req *http.Request, v interface{}) error {
	req.Header.Set("Authorization", "Bearer "+c.BearerToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return decodeJSON(resp.Body, v)
}

func decodeJSON(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, url, body)
}
