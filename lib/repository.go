package lib

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
	// Note: For some reason the API does not return tags in an array but as a map
	req, err := NewRequest("GET", "https://quay.io/api/v1/repository/"+namespace+"/"+repository, nil)
	if err != nil {
		return RepositoryWithTags{}, err
	}

	var repo Repository
	err = c.Get(req, &repo)
	if err != nil {
		return RepositoryWithTags{}, err
	}

	req, err = NewRequest("GET", "https://quay.io/api/v1/repository/"+namespace+"/"+repository+"/tag", nil)
	if err != nil {
		return RepositoryWithTags{}, err
	}

	var tags RepositoryTags
	err = c.Get(req, &tags)
	if err != nil {
		return RepositoryWithTags{}, err
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
