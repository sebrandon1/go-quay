/*
Package lib provides Quay.io API client functionality.

This file covers MESSAGES operations:

System Messages:
  - GET /api/v1/messages - GetMessages()

The Messages API returns system-wide messages for the authenticated user,
such as maintenance notifications or important announcements.
*/
package lib

import (
	"fmt"
)

// GetMessages retrieves system messages for the user
func (c *Client) GetMessages() (*Messages, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/messages", QuayURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create messages request: %w", err)
	}

	var messages Messages
	if err := c.get(req, &messages); err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	return &messages, nil
}

// CreateMessage creates a new global message (admin only)
func (c *Client) CreateMessage(content, severity, mediaType string) (*Message, error) {
	body := struct {
		Message struct {
			Content   string `json:"content"`
			Severity  string `json:"severity"`
			MediaType string `json:"media_type"`
		} `json:"message"`
	}{
		Message: struct {
			Content   string `json:"content"`
			Severity  string `json:"severity"`
			MediaType string `json:"media_type"`
		}{
			Content:   content,
			Severity:  severity,
			MediaType: mediaType,
		},
	}

	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/messages", QuayURL), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create message request: %w", err)
	}

	var message Message
	if err := c.post(req, &message); err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	return &message, nil
}
