package lib

import (
	"context"
	"testing"
	"time"
)

func TestRateLimitedClient_Passthrough(t *testing.T) {
	mock := &mockReader{
		repo: RepositoryWithTags{
			Repository: Repository{Name: testPlaceholder},
		},
	}
	rl := NewRateLimitedClient(mock, 100, 100)

	ctx := context.Background()

	repo, err := rl.GetRepository(ctx, testNamespace, testPlaceholder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.Name != testPlaceholder {
		t.Fatalf("expected name %q, got %q", testPlaceholder, repo.Name)
	}
}

func TestRateLimitedClient_ContextCancelled(t *testing.T) {
	mock := &mockReader{
		repo: RepositoryWithTags{
			Repository: Repository{Name: testPlaceholder},
		},
	}
	// Very low rate to force waiting
	rl := NewRateLimitedClient(mock, 0.001, 1)

	ctx := context.Background()

	// Use the one burst token
	_, _ = rl.GetRepository(ctx, testNamespace, testPlaceholder)

	// Next call should block; cancel the context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := rl.GetRepository(ctx, testNamespace, "test2")
	if err == nil {
		t.Fatal("expected error from canceled context")
	}
}

func TestRateLimitedClient_ListTags(t *testing.T) {
	mock := &mockReader{}
	rl := NewRateLimitedClient(mock, 100, 100)

	tags, err := rl.ListTags(context.Background(), testNamespace, testRepoName, 10, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tags == nil {
		t.Fatal("expected non-nil tags")
	}
}

func TestRateLimitedClient_ListAllTags_ContextCanceled(t *testing.T) {
	mock := &mockReader{}
	rl := NewRateLimitedClient(mock, 0.001, 1)

	// Use the one burst token
	_, _ = rl.ListAllTags(context.Background(), testNamespace, testRepoName, true)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := rl.ListAllTags(ctx, testNamespace, testRepoName, true)
	if err == nil {
		t.Fatal("expected error from canceled context")
	}
}

func TestNewCachedRateLimitedClient(t *testing.T) {
	mock := &mockReader{
		repo: RepositoryWithTags{
			Repository: Repository{Name: testPlaceholder},
		},
	}
	client := NewCachedRateLimitedClient(mock, time.Hour, 100, 100)

	ctx := context.Background()

	// First call
	_, err := client.GetRepository(ctx, testNamespace, testPlaceholder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Second call should be cached
	_, err = client.GetRepository(ctx, testNamespace, testPlaceholder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mock.calls.Load() != 1 {
		t.Fatalf("expected 1 call (cached), got %d", mock.calls.Load())
	}
}
