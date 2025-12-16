/*
Package lib provides Quay.io API client functionality.

This file covers MANIFEST operations:

Manifest Management:
  - GET    /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}                - GetManifest()
  - DELETE /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}                - DeleteManifest()

Manifest Labels:
  - GET    /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}/labels         - GetManifestLabels()
  - POST   /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}/labels         - AddManifestLabel()
  - GET    /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}/labels/{labelid} - GetManifestLabel()
  - DELETE /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}/labels/{labelid} - DeleteManifestLabel()

Manifest operations provide access to container image manifests, their layers,
configuration, and labels for inspection and management.
*/
package lib

import (
	"fmt"
)

// GetManifest retrieves detailed information about a specific manifest
func (c *Client) GetManifest(namespace, repository, manifestRef string) (*Manifest, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/manifest/%s", QuayURL, namespace, repository, manifestRef), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get manifest request: %w", err)
	}

	var manifest Manifest
	if err := c.get(req, &manifest); err != nil {
		return nil, fmt.Errorf("failed to get manifest: %w", err)
	}

	return &manifest, nil
}

// DeleteManifest deletes a specific manifest from a repository
func (c *Client) DeleteManifest(namespace, repository, manifestRef string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/repository/%s/%s/manifest/%s", QuayURL, namespace, repository, manifestRef), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete manifest request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete manifest: %w", err)
	}

	return nil
}

// GetManifestLabels retrieves all labels for a specific manifest
func (c *Client) GetManifestLabels(namespace, repository, manifestRef string) (*ManifestLabels, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/manifest/%s/labels", QuayURL, namespace, repository, manifestRef), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get manifest labels request: %w", err)
	}

	var labels ManifestLabels
	if err := c.get(req, &labels); err != nil {
		return nil, fmt.Errorf("failed to get manifest labels: %w", err)
	}

	return &labels, nil
}

// AddManifestLabel adds a label to a specific manifest
func (c *Client) AddManifestLabel(namespace, repository, manifestRef, key, value, mediaType string) (*ManifestLabel, error) {
	addReq := AddManifestLabelRequest{
		Key:   key,
		Value: value,
	}

	if mediaType != "" {
		addReq.MediaType = mediaType
	}

	req, err := newRequestWithBody("POST", fmt.Sprintf("%s/repository/%s/%s/manifest/%s/labels", QuayURL, namespace, repository, manifestRef), addReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create add manifest label request: %w", err)
	}

	var label ManifestLabel
	if err := c.post(req, &label); err != nil {
		return nil, fmt.Errorf("failed to add manifest label: %w", err)
	}

	return &label, nil
}

// GetManifestLabel retrieves a specific label from a manifest
func (c *Client) GetManifestLabel(namespace, repository, manifestRef, labelID string) (*ManifestLabel, error) {
	req, err := newRequest("GET", fmt.Sprintf("%s/repository/%s/%s/manifest/%s/labels/%s", QuayURL, namespace, repository, manifestRef, labelID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get manifest label request: %w", err)
	}

	var label ManifestLabel
	if err := c.get(req, &label); err != nil {
		return nil, fmt.Errorf("failed to get manifest label: %w", err)
	}

	return &label, nil
}

// DeleteManifestLabel deletes a specific label from a manifest
func (c *Client) DeleteManifestLabel(namespace, repository, manifestRef, labelID string) error {
	req, err := newRequest("DELETE", fmt.Sprintf("%s/repository/%s/%s/manifest/%s/labels/%s", QuayURL, namespace, repository, manifestRef, labelID), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete manifest label request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete manifest label: %w", err)
	}

	return nil
}
