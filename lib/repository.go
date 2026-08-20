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

import (
	"context"
	"fmt"
	"net/http"
)

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
func (c *Client) GetRepository(ctx context.Context, namespace, repository string) (RepositoryWithTags, error) {
	if namespace == "" {
		return RepositoryWithTags{}, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return RepositoryWithTags{}, fmt.Errorf("repository is required")
	}

	repoURL := c.buildURL("/repository/%s/%s", namespace, repository)
	req, err := newRequest(ctx, http.MethodGet, repoURL, nil)
	if err != nil {
		return RepositoryWithTags{}, fmt.Errorf("failed to create request for repository: %w", err)
	}

	var repo Repository
	if err := c.get(req, &repo); err != nil {
		return RepositoryWithTags{}, fmt.Errorf("failed to fetch repository details: %w", err)
	}

	tags, err := c.ListTags(ctx, namespace, repository, 0, false)
	if err != nil {
		return RepositoryWithTags{}, fmt.Errorf("failed to fetch repository tags: %w", err)
	}

	return RepositoryWithTags{
		Repository: repo,
		Tags:       *tags,
	}, nil
}

// CreateRepository creates a new repository
func (c *Client) CreateRepository(ctx context.Context, namespace, repository, visibility, description string) (*Repository, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}
	if visibility == "" {
		return nil, fmt.Errorf("visibility is required")
	}

	req, err := newRequestWithBody(ctx, http.MethodPost, c.buildURL("/repository"), CreateRepositoryRequest{
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
func (c *Client) UpdateRepository(ctx context.Context, namespace, repository, description, visibility string) (*Repository, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}

	updateReq := UpdateRepositoryRequest{}

	// Only include fields that are not empty
	if description != "" {
		updateReq.Description = description
	}
	if visibility != "" {
		updateReq.Visibility = visibility
	}

	req, err := newRequestWithBody(ctx, http.MethodPut, c.buildURL("/repository/%s/%s", namespace, repository), updateReq)
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
func (c *Client) DeleteRepository(ctx context.Context, namespace, repository string) error {
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return fmt.Errorf("repository is required")
	}

	req, err := newRequest(ctx, http.MethodDelete, c.buildURL("/repository/%s/%s", namespace, repository), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete repository request: %w", err)
	}

	if err := c.delete(req); err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}

	return nil
}

// ListRepositories lists all repositories visible to the user
func (c *Client) ListRepositories(ctx context.Context, namespace string, public, starred, popularity bool, page, limit int) (*RepositoryList, error) {
	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/repository"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list repositories request: %w", err)
	}

	params := map[string]string{}
	if namespace != "" {
		params["namespace"] = namespace
	}
	if public {
		params["public"] = queryValueTrue
	}
	if starred {
		params["starred"] = queryValueTrue
	}
	if popularity {
		params["popularity"] = queryValueTrue
	}
	if page > 0 {
		params["page"] = fmt.Sprintf("%d", page)
	}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}
	addQueryParams(req, params)

	var repos RepositoryList
	if err := c.get(req, &repos); err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	return &repos, nil
}

// ListAllRepositories fetches all repositories by following pagination automatically.
func (c *Client) ListAllRepositories(ctx context.Context, namespace string, public, starred, popularity bool) ([]OrganizationRepository, error) {
	var all []OrganizationRepository
	page := 1
	const pageSize = 100

	for {
		repos, err := c.ListRepositories(ctx, namespace, public, starred, popularity, page, pageSize)
		if err != nil {
			return nil, err
		}

		all = append(all, repos.Repositories...)

		if !repos.HasAdditional {
			break
		}
		page++
	}

	return all, nil
}

// ListAllTags fetches all tags for a repository by following pagination automatically.
func (c *Client) ListAllTags(ctx context.Context, namespace, repository string, onlyActive bool) ([]Tag, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}

	var all []Tag
	page := 1
	const pageSize = 100

	for {
		tags, err := c.ListTagsPage(ctx, namespace, repository, pageSize, page, onlyActive)
		if err != nil {
			return nil, err
		}

		all = append(all, tags.Tags...)

		if !tags.HasAdditional {
			break
		}
		page++
	}

	return all, nil
}

// ChangeRepositoryVisibility changes the visibility (public/private) of a repository
func (c *Client) ChangeRepositoryVisibility(ctx context.Context, namespace, repository, visibility string) error {
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return fmt.Errorf("repository is required")
	}
	if visibility == "" {
		return fmt.Errorf("visibility is required")
	}

	body := struct {
		Visibility string `json:"visibility"`
	}{
		Visibility: visibility,
	}
	req, err := newRequestWithBody(ctx, http.MethodPost, c.buildURL("/repository/%s/%s/changevisibility", namespace, repository), body)
	if err != nil {
		return fmt.Errorf("failed to create change visibility request: %w", err)
	}

	if err := c.post(req, nil); err != nil {
		return fmt.Errorf("failed to change repository visibility: %w", err)
	}

	return nil
}

// ListTagsPage lists tags for a repository with pagination support.
func (c *Client) ListTagsPage(ctx context.Context, namespace, repository string, limit, page int, onlyActive bool) (*RepositoryTags, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/repository/%s/%s/tag/", namespace, repository), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list tags request: %w", err)
	}

	params := map[string]string{}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}
	if page > 0 {
		params["page"] = fmt.Sprintf("%d", page)
	}
	if onlyActive {
		params["onlyActiveTags"] = queryValueTrue
	}
	addQueryParams(req, params)

	var tags RepositoryTags
	if err := c.get(req, &tags); err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	return &tags, nil
}

// ListTags lists tags for a repository
func (c *Client) ListTags(ctx context.Context, namespace, repository string, limit int, onlyActive bool) (*RepositoryTags, error) {
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if repository == "" {
		return nil, fmt.Errorf("repository is required")
	}

	req, err := newRequest(ctx, http.MethodGet, c.buildURL("/repository/%s/%s/tag/", namespace, repository), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list tags request: %w", err)
	}

	params := map[string]string{}
	if limit > 0 {
		params["limit"] = fmt.Sprintf("%d", limit)
	}
	if onlyActive {
		params["onlyActiveTags"] = queryValueTrue
	}
	addQueryParams(req, params)

	var tags RepositoryTags
	if err := c.get(req, &tags); err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	return &tags, nil
}
