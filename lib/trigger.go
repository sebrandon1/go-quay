/*
Package lib provides Quay.io API client functionality.

This file covers BUILD TRIGGER operations:

Trigger Management:
  - GET    /api/v1/repository/{namespace}/{repository}/trigger/                       - GetTriggers()
  - GET    /api/v1/repository/{namespace}/{repository}/trigger/{trigger_uuid}         - GetTrigger()
  - DELETE /api/v1/repository/{namespace}/{repository}/trigger/{trigger_uuid}         - DeleteTrigger()
  - PUT    /api/v1/repository/{namespace}/{repository}/trigger/{trigger_uuid}         - UpdateTrigger()
  - POST   /api/v1/repository/{namespace}/{repository}/trigger/{trigger_uuid}/start   - StartTriggerBuild()
  - POST   /api/v1/repository/{namespace}/{repository}/trigger/{trigger_uuid}/activate - ActivateTrigger()

Build triggers allow automated image builds when code is pushed to connected
source repositories like GitHub, GitLab, or Bitbucket.

Supported trigger services:
  - github: GitHub repository
  - gitlab: GitLab repository
  - bitbucket: Bitbucket repository
  - custom-git: Custom git repository
*/
package lib

import (
	"fmt"
)

// GetTriggers retrieves all build triggers for a repository
func (c *Client) GetTriggers(namespace, repository string) (*BuildTriggers, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/trigger/", QuayURL, namespace, repository), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get triggers request: %w", err)
	}

	var triggers BuildTriggers
	if err := c.get(req, &triggers); err != nil {
		return nil, fmt.Errorf("failed to get triggers: %w", err)
	}

	return &triggers, nil
}

// GetTrigger retrieves a specific build trigger by UUID
func (c *Client) GetTrigger(namespace, repository, triggerUUID string) (*BuildTrigger, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/trigger/%s", QuayURL, namespace, repository, triggerUUID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get trigger request: %w", err)
	}

	var trigger BuildTrigger
	if err := c.get(req, &trigger); err != nil {
		return nil, fmt.Errorf("failed to get trigger: %w", err)
	}

	return &trigger, nil
}

// DeleteTrigger deletes a build trigger
func (c *Client) DeleteTrigger(namespace, repository, triggerUUID string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/repository/%s/%s/trigger/%s", QuayURL, namespace, repository, triggerUUID), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete trigger request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete trigger: %w", err)
	}

	return nil
}

// UpdateTrigger updates a build trigger (enable/disable)
func (c *Client) UpdateTrigger(namespace, repository, triggerUUID string, enabled bool) (*BuildTrigger, error) {
	body := map[string]interface{}{
		"enabled": enabled,
	}

	req, err := newRequestWithBody("PUT", fmt.Sprintf("%s/repository/%s/%s/trigger/%s", QuayURL, namespace, repository, triggerUUID), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create update trigger request: %w", err)
	}

	var trigger BuildTrigger
	if err := c.put(req, &trigger); err != nil {
		return nil, fmt.Errorf("failed to update trigger: %w", err)
	}

	return &trigger, nil
}

// StartTriggerBuild manually starts a build from a trigger
func (c *Client) StartTriggerBuild(namespace, repository, triggerUUID string, triggerReq *ManualTriggerRequest) (*Build, error) {
	var body interface{}
	if triggerReq != nil {
		body = triggerReq
	}

	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/repository/%s/%s/trigger/%s/start", QuayURL, namespace, repository, triggerUUID), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create start trigger build request: %w", err)
	}

	var build Build
	if err := c.post(req, &build); err != nil {
		return nil, fmt.Errorf("failed to start trigger build: %w", err)
	}

	return &build, nil
}

// ActivateTrigger activates a build trigger with configuration
func (c *Client) ActivateTrigger(namespace, repository, triggerUUID string, activateReq *ActivateTriggerRequest) (*BuildTrigger, error) {
	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/repository/%s/%s/trigger/%s/activate", QuayURL, namespace, repository, triggerUUID), activateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create activate trigger request: %w", err)
	}

	var trigger BuildTrigger
	if err := c.post(req, &trigger); err != nil {
		return nil, fmt.Errorf("failed to activate trigger: %w", err)
	}

	return &trigger, nil
}
