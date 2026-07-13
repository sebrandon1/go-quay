/*
Package lib provides Quay.io API client functionality.

This file covers REPOSITORY endpoints:

Repository Management:
  - POST /api/v1/repository                                - CreateRepository()
  - GET  /api/v1/repository                                - ListRepositories()
  - GET  /api/v1/repository/{namespace}/{repository}       - GetRepository()
  - PUT  /api/v1/repository/{namespace}/{repository}       - UpdateRepository()
  - DELETE /api/v1/repository/{namespace}/{repository}     - DeleteRepository()
  - GET  /api/v1/repository/{namespace}/{repository}/tag/  - ListTags()

GetRepository() combines repository details with tag information via ListTags().
ListRepositories() supports a popularity flag for pull count data.
*/
package lib

import "fmt"

type RepositoryTags struct {
	Tags          []Tag `json:"tags,omitempty"`
	Page          int   `json:"page,omitempty"`
	HasAdditional bool  `json:"has_additional,omitempty"`
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
	Repository
	Tags RepositoryTags `json:"tags,omitempty"`
}

// GetRepository returns a repository with tags information baked in
func (c *Client) GetRepository(namespace, repository string) (RepositoryWithTags, error) {
	repoURL := c.buildURL("/repository/%s/%s", namespace, repository)
	req, err := newRequest("GET", repoURL, nil)
	if err != nil {
		return RepositoryWithTags{}, fmt.Errorf("failed to create request for repository: %w", err)
	}

	var repo Repository
	if err := c.get(req, &repo); err != nil {
		return RepositoryWithTags{}, fmt.Errorf("failed to fetch repository details: %w", err)
	}

	tags, err := c.ListTags(namespace, repository, 0, false)
	if err != nil {
		return RepositoryWithTags{}, fmt.Errorf("failed to fetch repository tags: %w", err)
	}

	return RepositoryWithTags{
		Repository: repo,
		Tags:       *tags,
	}, nil
}

// CreateRepository creates a new repository
func (c *Client) CreateRepository(namespace, repository, visibility, description string) (*Repository, error) {
	req, err := newRequestWithBody("POST", c.BaseURL+"/repository", CreateRepositoryRequest{
		Repository:  repository,
		Namespace:   namespace,
		Visibility:  visibility,
		Description: description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create repository request: %w", err)
	}

	var repo Repository
	if err := c.post(req, &repo); err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	return &repo, nil
}

// UpdateRepository updates an existing repository
func (c *Client) UpdateRepository(namespace, repository, description, visibility string) (*Repository, error) {
	updateReq := UpdateRepositoryRequest{}

	// Only include fields that are not empty
	if description != "" {
		updateReq.Description = description
	}
	if visibility != "" {
		updateReq.Visibility = visibility
	}

	req, err := newRequestWithBody("PUT", c.buildURL("/repository/%s/%s", namespace, repository), updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create update repository request: %w", err)
	}

	var repo Repository
	if err := c.put(req, &repo); err != nil {
		return nil, fmt.Errorf("failed to update repository: %w", err)
	}

	return &repo, nil
}

// DeleteRepository deletes a repository
func (c *Client) DeleteRepository(namespace, repository string) error {
	req, err := newRequest("DELETE", c.buildURL("/repository/%s/%s", namespace, repository), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete repository request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}

	return nil
}

// ListRepositories lists all repositories visible to the user
func (c *Client) ListRepositories(namespace string, public, starred, popularity bool, page, limit int) (*RepositoryList, error) {
	req, err := newRequest("GET", c.buildURL("/repository"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list repositories request: %w", err)
	}

	q := req.URL.Query()
	if namespace != "" {
		q.Add("namespace", namespace)
	}
	if public {
		q.Add("public", "true")
	}
	if starred {
		q.Add("starred", "true")
	}
	if popularity {
		q.Add("popularity", "true")
	}
	if page > 0 {
		q.Add("page", fmt.Sprintf("%d", page))
	}
	if limit > 0 {
		q.Add("limit", fmt.Sprintf("%d", limit))
	}
	req.URL.RawQuery = q.Encode()

	var repos RepositoryList
	if err := c.get(req, &repos); err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	return &repos, nil
}

// ChangeRepositoryVisibility changes the visibility (public/private) of a repository
func (c *Client) ChangeRepositoryVisibility(namespace, repository, visibility string) error {
	body := struct {
		Visibility string `json:"visibility"`
	}{
		Visibility: visibility,
	}
	req, err := newRequestWithBody("POST", c.buildURL("/repository/%s/%s/changevisibility", namespace, repository), body)
	if err != nil {
		return fmt.Errorf("failed to create change visibility request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to change repository visibility: %w", err)
	}

	return nil
}

// ListTags lists tags for a repository
func (c *Client) ListTags(namespace, repository string, limit int, onlyActive bool) (*RepositoryTags, error) {
	req, err := newRequest("GET", c.buildURL("/repository/%s/%s/tag/", namespace, repository), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list tags request: %w", err)
	}

	q := req.URL.Query()
	if limit > 0 {
		q.Add("limit", fmt.Sprintf("%d", limit))
	}
	if onlyActive {
		q.Add("onlyActiveTags", "true")
	}
	req.URL.RawQuery = q.Encode()

	var tags RepositoryTags
	if err := c.get(req, &tags); err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	return &tags, nil
}
