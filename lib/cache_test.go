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
