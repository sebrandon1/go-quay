package lib

import "fmt"

type RepositoryTags struct {
	Tags []struct {
		Name           string `json:"name,omitempty"`
		Reversion      bool   `json:"reversion,omitempty"`
		StartTs        int    `json:"start_ts,omitempty"`
		ManifestDigest string `json:"manifest_digest,omitempty"`
		IsManifestList bool   `json:"is_manifest_list,omitempty"`
		Size           any    `json:"size,omitempty"`
		LastModified   string `json:"last_modified,omitempty"`
		EndTs          int    `json:"end_ts,omitempty"`
		Expiration     string `json:"expiration,omitempty"`
	} `json:"tags,omitempty"`
	Page          int  `json:"page,omitempty"`
	HasAdditional bool `json:"has_additional,omitempty"`
}

type Repository struct {
	Namespace      string `json:"namespace,omitempty"`
	Name           string `json:"name,omitempty"`
	Kind           string `json:"kind,omitempty"`
	Description    string `json:"description,omitempty"`
	IsPublic       bool   `json:"is_public,omitempty"`
	IsOrganization bool   `json:"is_organization,omitempty"`
	IsStarred      bool   `json:"is_starred,omitempty"`
	StatusToken    string `json:"status_token,omitempty"`
	TrustEnabled   bool   `json:"trust_enabled,omitempty"`
	TagExpirationS int    `json:"tag_expiration_s,omitempty"`
	IsFreeAccount  bool   `json:"is_free_account,omitempty"`
	State          string `json:"state,omitempty"`
	CanWrite       bool   `json:"can_write,omitempty"`
	CanAdmin       bool   `json:"can_admin,omitempty"`
}

type RepositoryWithTags struct {
	Namespace      string         `json:"namespace,omitempty"`
	Name           string         `json:"name,omitempty"`
	Kind           string         `json:"kind,omitempty"`
	Description    string         `json:"description,omitempty"`
	IsPublic       bool           `json:"is_public,omitempty"`
	IsOrganization bool           `json:"is_organization,omitempty"`
	IsStarred      bool           `json:"is_starred,omitempty"`
	StatusToken    string         `json:"status_token,omitempty"`
	TrustEnabled   bool           `json:"trust_enabled,omitempty"`
	TagExpirationS int            `json:"tag_expiration_s,omitempty"`
	IsFreeAccount  bool           `json:"is_free_account,omitempty"`
	State          string         `json:"state,omitempty"`
	CanWrite       bool           `json:"can_write,omitempty"`
	CanAdmin       bool           `json:"can_admin,omitempty"`
	Tags           RepositoryTags `json:"tags,omitempty"`
}

// GetRepository returns a repository with tags information baked in
func (c *Client) GetRepository(namespace, repository string) (RepositoryWithTags, error) {
	baseURL := "https://quay.io/api/v1/repository/"

	// Fetch repository details
	repoURL := fmt.Sprintf("%s%s/%s", baseURL, namespace, repository)
	req, err := newRequest("GET", repoURL, nil)
	if err != nil {
		return RepositoryWithTags{}, fmt.Errorf("failed to create request for repository: %w", err)
	}

	var repo Repository
	if err := c.get(req, &repo); err != nil {
		return RepositoryWithTags{}, fmt.Errorf("failed to fetch repository details: %w", err)
	}

	// Fetch repository tags
	tagsURL := fmt.Sprintf("%s%s/%s/tag", baseURL, namespace, repository)
	req, err = newRequest("GET", tagsURL, nil)
	if err != nil {
		return RepositoryWithTags{}, fmt.Errorf("failed to create request for tags: %w", err)
	}

	var tags RepositoryTags
	if err := c.get(req, &tags); err != nil {
		return RepositoryWithTags{}, fmt.Errorf("failed to fetch repository tags: %w", err)
	}

	repoWithTags := RepositoryWithTags{
		Namespace:      repo.Namespace,
		Name:           repo.Name,
		Kind:           repo.Kind,
		Description:    repo.Description,
		IsPublic:       repo.IsPublic,
		IsOrganization: repo.IsOrganization,
		IsStarred:      repo.IsStarred,
		StatusToken:    repo.StatusToken,
		TrustEnabled:   repo.TrustEnabled,
		TagExpirationS: repo.TagExpirationS,
		IsFreeAccount:  repo.IsFreeAccount,
		State:          repo.State,
		CanWrite:       repo.CanWrite,
		CanAdmin:       repo.CanAdmin,
		Tags:           tags,
	}

	return repoWithTags, nil
}
