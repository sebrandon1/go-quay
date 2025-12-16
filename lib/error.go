/*
Package lib provides Quay.io API client functionality.

This file covers ERROR TYPE operations:

Error Information:
  - GET /api/v1/error/{error_type} - GetErrorType()

The Error API provides detailed information about specific error types
that can be returned by the Quay.io API.
*/
package lib

import (
	"fmt"
)

// GetErrorType retrieves details about a specific error type
func (c *Client) GetErrorType(errorType string) (*ErrorType, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/error/%s", QuayURL, errorType), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create error type request: %w", err)
	}

	var errType ErrorType
	if err := c.get(req, &errType); err != nil {
		return nil, fmt.Errorf("failed to get error type: %w", err)
	}

	return &errType, nil
}
