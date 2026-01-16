/*
Package lib provides Quay.io API client functionality.

This file covers REPOSITORY NOTIFICATION operations:

Notification Management:
  - GET    /api/v1/repository/{namespace}/{repository}/notification/              - GetNotifications()
  - GET    /api/v1/repository/{namespace}/{repository}/notification/{uuid}        - GetNotification()
  - POST   /api/v1/repository/{namespace}/{repository}/notification/              - CreateNotification()
  - DELETE /api/v1/repository/{namespace}/{repository}/notification/{uuid}        - DeleteNotification()
  - POST   /api/v1/repository/{namespace}/{repository}/notification/{uuid}/test   - TestNotification()
  - POST   /api/v1/repository/{namespace}/{repository}/notification/{uuid}/reset  - ResetNotification()

Repository notifications (webhooks) allow external services to be notified
of events such as image pushes, builds, and vulnerability scans.

Supported events:
  - repo_push: Image push to repository
  - build_queued: Build has been queued
  - build_start: Build has started
  - build_success: Build completed successfully
  - build_failure: Build failed
  - build_canceled: Build was canceled
  - vulnerability_found: New vulnerability discovered

Supported methods:
  - webhook: HTTP webhook
  - email: Email notification
  - slack: Slack notification
  - hipchat: HipChat notification
  - flowdock: Flowdock notification
*/
package lib

import (
	"fmt"
)

// GetNotifications retrieves all notifications for a repository
func (c *Client) GetNotifications(namespace, repository string) (*RepositoryNotifications, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/notification/", QuayURL, namespace, repository), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get notifications request: %w", err)
	}

	var notifications RepositoryNotifications
	if err := c.get(req, &notifications); err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	return &notifications, nil
}

// GetNotification retrieves a specific notification by UUID
func (c *Client) GetNotification(namespace, repository, uuid string) (*RepositoryNotification, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/notification/%s", QuayURL, namespace, repository, uuid), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get notification request: %w", err)
	}

	var notification RepositoryNotification
	if err := c.get(req, &notification); err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	return &notification, nil
}

// CreateNotification creates a new notification for a repository
func (c *Client) CreateNotification(namespace, repository string, notificationReq *CreateNotificationRequest) (*RepositoryNotification, error) {
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/repository/%s/%s/notification/", QuayURL, namespace, repository), notificationReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification request: %w", err)
	}

	var notification RepositoryNotification
	if err := c.post(req, &notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return &notification, nil
}

// DeleteNotification deletes a notification from a repository
func (c *Client) DeleteNotification(namespace, repository, uuid string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/repository/%s/%s/notification/%s", QuayURL, namespace, repository, uuid), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete notification request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	return nil
}

// TestNotification tests a notification by sending a test event
func (c *Client) TestNotification(namespace, repository, uuid string) error {
	req, err := newRequest("POST", fmt.Sprintf("%s/repository/%s/%s/notification/%s/test", QuayURL, namespace, repository, uuid), nil)
	if err != nil {
		return fmt.Errorf("failed to create test notification request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to test notification: %w", err)
	}

	return nil
}

// ResetNotification resets failure count for a notification
func (c *Client) ResetNotification(namespace, repository, uuid string) error {
	req, err := newRequest("POST", fmt.Sprintf("%s/repository/%s/%s/notification/%s/reset", QuayURL, namespace, repository, uuid), nil)
	if err != nil {
		return fmt.Errorf("failed to create reset notification request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to reset notification: %w", err)
	}

	return nil
}

// UpdateNotification updates an existing notification
func (c *Client) UpdateNotification(namespace, repository, uuid string, notificationReq *CreateNotificationRequest) (*RepositoryNotification, error) {
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/repository/%s/%s/notification/%s", QuayURL, namespace, repository, uuid), notificationReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create update notification request: %w", err)
	}

	var notification RepositoryNotification
	if err := c.post(req, &notification); err != nil {
		return nil, fmt.Errorf("failed to update notification: %w", err)
	}

	return &notification, nil
}
