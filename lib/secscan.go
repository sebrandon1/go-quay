/*
Package lib provides Quay.io API client functionality.

This file covers SECURITY SCAN operations:

Security Scanning:
  - GET /api/v1/repository/{namespace}/{repository}/manifest/{manifestref}/security - GetManifestSecurity()

Security scan operations provide access to vulnerability information for
container images, including CVE details, severity levels, and fix versions.
*/
package lib

import (
	"fmt"
)

// GetManifestSecurity retrieves security scan information for a specific manifest
func (c *Client) GetManifestSecurity(namespace, repository, manifestRef string, vulnerabilities bool) (*SecurityScan, error) {
	url := fmt.Sprintf("%s/repository/%s/%s/manifest/%s/security", QuayURL, namespace, repository, manifestRef)

	// Add query parameter to include vulnerabilities
	if vulnerabilities {
		url += "?vulnerabilities=true"
	}

	req, err := newRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get manifest security request: %w", err)
	}

	var securityScan SecurityScan
	if err := c.get(req, &securityScan); err != nil {
		return nil, fmt.Errorf("failed to get manifest security: %w", err)
	}

	return &securityScan, nil
}
