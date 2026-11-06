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
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ScanStatus is the state of a Quay security scan.
type ScanStatus string

const (
	ScanQueued      ScanStatus = "queued"
	ScanScanning    ScanStatus = "scanning"
	ScanScanned     ScanStatus = "scanned"
	ScanFailed      ScanStatus = "failed"
	ScanUnsupported ScanStatus = "unsupported"
)

// Normalize lowercases and trims the status so comparisons are insensitive to
// how Quay spelled it.
func (s ScanStatus) Normalize() ScanStatus {
	return ScanStatus(strings.ToLower(strings.TrimSpace(string(s))))
}

// Terminal reports whether the scan can make no further progress. Unknown
// statuses are not terminal, so callers keep waiting rather than failing.
func (s ScanStatus) Terminal() bool {
	switch s.Normalize() {
	case ScanScanned, ScanFailed, ScanUnsupported:
		return true
	}
	return false
}

// GetManifestSecurity retrieves security scan information for a specific manifest
func (c *Client) GetManifestSecurity(ctx context.Context, namespace, repository, manifestRef string, vulnerabilities bool) (*SecurityScan, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	if manifestRef == "" {
		return nil, fmt.Errorf("manifestRef is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/repository/%s/%s/manifest/%s/security", namespace, repository, manifestRef))
	if err != nil {
		return nil, fmt.Errorf("failed to create get manifest security request: %w", err)
	}

	if vulnerabilities {
		addQueryParams(req, map[string]string{"vulnerabilities": queryValueTrue})
	}

	var securityScan SecurityScan
	if err := c.get(req, &securityScan); err != nil {
		return nil, fmt.Errorf("failed to get manifest security: %w", err)
	}

	return &securityScan, nil
}
