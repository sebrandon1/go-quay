package lib

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestCachedClient_HitAndMiss(t *testing.T) {
	mock := &mockReader{
		repo: RepositoryWithTags{
			Repository: Repository{Name: testPlaceholder, Namespace: testNamespace},
		},
	}
	cached := NewCachedClient(mock, WithCacheTTL(time.Hour))

	ctx := context.Background()

	// First call: cache miss
	repo, err := cached.GetRepository(ctx, testNamespace, testPlaceholder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.Name != testPlaceholder {
		t.Fatalf("expected name 'test', got %q", repo.Name)
	}
	if mock.calls.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", mock.calls.Load())
	}

	// Second call: cache hit
	repo, err = cached.GetRepository(ctx, testNamespace, testPlaceholder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.Name != testPlaceholder {
		t.Fatalf("expected name 'test', got %q", repo.Name)
	}
	if mock.calls.Load() != 1 {
		t.Fatalf("expected still 1 call after cache hit, got %d", mock.calls.Load())
	}
}

func TestCachedClient_Expiry(t *testing.T) {
	mock := &mockReader{
		repo: RepositoryWithTags{
			Repository: Repository{Name: testPlaceholder},
		},
	}
	cached := NewCachedClient(mock, WithCacheTTL(10*time.Millisecond))

	ctx := context.Background()

	if _, err := cached.GetRepository(ctx, testNamespace, testPlaceholder); err != nil {
		t.Fatalf("GetRepository: %v", err)
	}
	if mock.calls.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", mock.calls.Load())
	}

	time.Sleep(20 * time.Millisecond)

	if _, err := cached.GetRepository(ctx, testNamespace, testPlaceholder); err != nil {
		t.Fatalf("GetRepository: %v", err)
	}
	if mock.calls.Load() != 2 {
		t.Fatalf("expected 2 calls after TTL expiry, got %d", mock.calls.Load())
	}
}

func TestCachedClient_ClearCache(t *testing.T) {
	mock := &mockReader{
		repo: RepositoryWithTags{
			Repository: Repository{Name: testPlaceholder},
		},
	}
	cached := NewCachedClient(mock)

	ctx := context.Background()

	if _, err := cached.GetRepository(ctx, testNamespace, testPlaceholder); err != nil {
		t.Fatalf("GetRepository: %v", err)
	}
	cached.ClearCache()
	if _, err := cached.GetRepository(ctx, testNamespace, testPlaceholder); err != nil {
		t.Fatalf("GetRepository: %v", err)
	}

	if mock.calls.Load() != 2 {
		t.Fatalf("expected 2 calls after ClearCache, got %d", mock.calls.Load())
	}
}

func TestCachedClient_CleanupExpired(t *testing.T) {
	mock := &mockReader{
		repo: RepositoryWithTags{
			Repository: Repository{Name: testPlaceholder},
		},
	}
	cached := NewCachedClient(mock, WithCacheTTL(10*time.Millisecond))

	ctx := context.Background()

	if _, err := cached.GetRepository(ctx, testNamespace, testPlaceholder); err != nil {
		t.Fatalf("GetRepository: %v", err)
	}
	time.Sleep(20 * time.Millisecond)
	cached.CleanupExpired()

	if _, err := cached.GetRepository(ctx, testNamespace, testPlaceholder); err != nil {
		t.Fatalf("GetRepository: %v", err)
	}
	if mock.calls.Load() != 2 {
		t.Fatalf("expected 2 calls after CleanupExpired, got %d", mock.calls.Load())
	}
}

func TestCachedClient_ErrorNotCached(t *testing.T) {
	mock := &mockReader{
		err: fmt.Errorf("api error"),
	}
	cached := NewCachedClient(mock)

	ctx := context.Background()

	_, err := cached.GetRepository(ctx, testNamespace, testPlaceholder)
	if err == nil {
		t.Fatal("expected error")
	}

	// Fix the mock to succeed
	mock.err = nil
	mock.repo = RepositoryWithTags{Repository: Repository{Name: testPlaceholder}}

	// Should NOT be cached — should call inner again
	repo, err := cached.GetRepository(ctx, testNamespace, testPlaceholder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.Name != testPlaceholder {
		t.Fatalf("expected name %q, got %q", testPlaceholder, repo.Name)
	}
	if mock.calls.Load() != 2 {
		t.Fatalf("expected 2 calls (error not cached), got %d", mock.calls.Load())
	}
}

func TestCachedClient_ListTagsPassthrough(t *testing.T) {
	mock := &mockReader{}
	cached := NewCachedClient(mock)

	tags, err := cached.ListTags(context.Background(), testNamespace, testRepoName, 10, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tags == nil {
		t.Fatal("expected non-nil tags")
	}
	if mock.calls.Load() != 1 {
		t.Fatalf("expected 1 call, got %d", mock.calls.Load())
	}
}

func TestCachedClient_ReaderPassthroughs(t *testing.T) {
	mock := &mockReader{}
	cached := NewCachedClient(mock)
	ctx := context.Background()
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "ListAllTags",
			call: func() error {
				_, err := cached.ListAllTags(ctx, testNamespace, testRepoName, true)
				return err
			},
		},
		{
			name: "GetManifestSecurity",
			call: func() error {
				_, err := cached.GetManifestSecurity(ctx, testNamespace, testRepoName, "latest", true)
				return err
			},
		},
		{
			name: "GetManifest",
			call: func() error {
				_, err := cached.GetManifest(ctx, testNamespace, testRepoName, "latest")
				return err
			},
		},
		{
			name: "GetManifestLabels",
			call: func() error {
				_, err := cached.GetManifestLabels(ctx, testNamespace, testRepoName, "latest")
				return err
			},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(); err != nil {
				t.Fatalf("passthrough returned error: %v", err)
			}
			if got := mock.calls.Load(); got != int32(i+1) {
				t.Errorf("inner call count = %d, want %d", got, i+1)
			}
		})
	}
}

func TestCachedClient_StartCleanupLoop(t *testing.T) {
	cached := NewCachedClient(&mockReader{}, WithCacheTTL(time.Hour))
	cached.cache["expired"] = cacheEntry{timestamp: time.Now().Add(-time.Hour)}

	ctx, cancel := context.WithCancel(context.Background())
	cached.StartCleanupLoop(ctx, 10*time.Millisecond)
	t.Cleanup(cancel)

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		cached.mu.RLock()
		_, exists := cached.cache["expired"]
		cached.mu.RUnlock()
		if !exists {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("cleanup loop did not remove the expired entry")
}

func TestCachedClient_DifferentKeys(t *testing.T) {
	mock := &mockReader{
		repo: RepositoryWithTags{
			Repository: Repository{Name: testPlaceholder},
		},
	}
	cached := NewCachedClient(mock)

	ctx := context.Background()

	if _, err := cached.GetRepository(ctx, "ns1", "repo1"); err != nil {
		t.Fatalf("GetRepository ns1: %v", err)
	}
	if _, err := cached.GetRepository(ctx, "ns2", "repo2"); err != nil {
		t.Fatalf("GetRepository ns2: %v", err)
	}

	if mock.calls.Load() != 2 {
		t.Fatalf("expected 2 calls for different keys, got %d", mock.calls.Load())
	}
}
