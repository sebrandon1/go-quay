package lib

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

const (
	DefaultRateLimit = 5.0
	DefaultRateBurst = 10
)

type RateLimitedClient struct {
	inner   RepositoryReader
	limiter *rate.Limiter
}

func NewRateLimitedClient(inner RepositoryReader, ratePerSecond float64, burst int) *RateLimitedClient {
	return &RateLimitedClient{
		inner:   inner,
		limiter: rate.NewLimiter(rate.Limit(ratePerSecond), burst),
	}
}

func (c *RateLimitedClient) GetRepository(ctx context.Context, namespace, repository string) (RepositoryWithTags, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return RepositoryWithTags{}, err
	}
	return c.inner.GetRepository(ctx, namespace, repository)
}

func (c *RateLimitedClient) ListTags(ctx context.Context, namespace, repository string, limit int, onlyActive bool) (*RepositoryTags, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	return c.inner.ListTags(ctx, namespace, repository, limit, onlyActive)
}

func (c *RateLimitedClient) ListAllTags(ctx context.Context, namespace, repository string, onlyActive bool) ([]Tag, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	return c.inner.ListAllTags(ctx, namespace, repository, onlyActive)
}

func (c *RateLimitedClient) GetManifestSecurity(ctx context.Context, namespace, repository, manifestRef string, vulnerabilities bool) (*SecurityScan, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	return c.inner.GetManifestSecurity(ctx, namespace, repository, manifestRef, vulnerabilities)
}

func (c *RateLimitedClient) GetManifest(ctx context.Context, namespace, repository, manifestRef string) (*Manifest, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	return c.inner.GetManifest(ctx, namespace, repository, manifestRef)
}

func (c *RateLimitedClient) GetManifestLabels(ctx context.Context, namespace, repository, manifestRef string) (*ManifestLabels, error) {
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, err
	}
	return c.inner.GetManifestLabels(ctx, namespace, repository, manifestRef)
}

func NewCachedRateLimitedClient(base RepositoryReader, cacheTTL time.Duration, ratePerSecond float64, burst int) *CachedClient {
	rl := NewRateLimitedClient(base, ratePerSecond, burst)
	return NewCachedClient(rl, WithCacheTTL(cacheTTL))
}

var _ RepositoryReader = (*RateLimitedClient)(nil)
