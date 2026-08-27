package lib

import "context"

type RepositoryReader interface {
	GetRepository(ctx context.Context, namespace, repository string) (RepositoryWithTags, error)
	ListTags(ctx context.Context, namespace, repository string, limit int, onlyActive bool) (*RepositoryTags, error)
	ListAllTags(ctx context.Context, namespace, repository string, onlyActive bool) ([]Tag, error)
	GetManifestSecurity(ctx context.Context, namespace, repository, manifestRef string, vulnerabilities bool) (*SecurityScan, error)
	GetManifest(ctx context.Context, namespace, repository, manifestRef string) (*Manifest, error)
	GetManifestLabels(ctx context.Context, namespace, repository, manifestRef string) (*ManifestLabels, error)
}
