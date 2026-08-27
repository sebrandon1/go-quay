package lib

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

const (
	DefaultCacheTTL = 1 * time.Hour
)

type cacheEntry struct {
	repo      RepositoryWithTags
	timestamp time.Time
}

// CachedClient wraps a RepositoryReader with an in-memory TTL cache.
// Only GetRepository results are cached; ListTags and ListAllTags pass through.
type CachedClient struct {
	inner RepositoryReader
	cache map[string]cacheEntry
	mu    sync.RWMutex
	ttl   time.Duration
	group singleflight.Group
}

type CacheOption func(*CachedClient)

func WithCacheTTL(ttl time.Duration) CacheOption {
	return func(c *CachedClient) {
		c.ttl = ttl
	}
}

func NewCachedClient(inner RepositoryReader, opts ...CacheOption) *CachedClient {
	c := &CachedClient{
		inner: inner,
		cache: make(map[string]cacheEntry),
		ttl:   DefaultCacheTTL,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func cacheKey(namespace, repository string) string {
	return namespace + "\x00" + repository
}

func (c *CachedClient) GetRepository(ctx context.Context, namespace, repository string) (RepositoryWithTags, error) {
	key := cacheKey(namespace, repository)

	c.mu.RLock()
	if entry, ok := c.cache[key]; ok && time.Since(entry.timestamp) < c.ttl {
		c.mu.RUnlock()
		return entry.repo, nil
	}
	c.mu.RUnlock()

	// singleflight deduplicates concurrent fetches for the same key
	v, err, _ := c.group.Do(key, func() (any, error) {
		return c.inner.GetRepository(ctx, namespace, repository)
	})
	if err != nil {
		return RepositoryWithTags{}, err
	}

	repo := v.(RepositoryWithTags)

	c.mu.Lock()
	c.cache[key] = cacheEntry{repo: repo, timestamp: time.Now()}
	c.mu.Unlock()

	return repo, nil
}

func (c *CachedClient) ListTags(ctx context.Context, namespace, repository string, limit int, onlyActive bool) (*RepositoryTags, error) {
	return c.inner.ListTags(ctx, namespace, repository, limit, onlyActive)
}

func (c *CachedClient) ListAllTags(ctx context.Context, namespace, repository string, onlyActive bool) ([]Tag, error) {
	return c.inner.ListAllTags(ctx, namespace, repository, onlyActive)
}

func (c *CachedClient) GetManifestSecurity(ctx context.Context, namespace, repository, manifestRef string, vulnerabilities bool) (*SecurityScan, error) {
	return c.inner.GetManifestSecurity(ctx, namespace, repository, manifestRef, vulnerabilities)
}

func (c *CachedClient) GetManifest(ctx context.Context, namespace, repository, manifestRef string) (*Manifest, error) {
	return c.inner.GetManifest(ctx, namespace, repository, manifestRef)
}

func (c *CachedClient) GetManifestLabels(ctx context.Context, namespace, repository, manifestRef string) (*ManifestLabels, error) {
	return c.inner.GetManifestLabels(ctx, namespace, repository, manifestRef)
}

func (c *CachedClient) ClearCache() {
	c.mu.Lock()
	c.cache = make(map[string]cacheEntry)
	c.mu.Unlock()
}

func (c *CachedClient) CleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, entry := range c.cache {
		if time.Since(entry.timestamp) >= c.ttl {
			delete(c.cache, key)
		}
	}
}

// StartCleanupLoop runs periodic cache cleanup. Cancel the context to stop.
func (c *CachedClient) StartCleanupLoop(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.CleanupExpired()
			}
		}
	}()
}

var _ RepositoryReader = (*CachedClient)(nil)
