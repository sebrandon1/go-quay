package lib

import (
	"testing"
	"time"
)

var _ RepositoryReader = (*mockReader)(nil)

func TestRepositoryReaderImplementations(t *testing.T) {
	client, err := NewClient(testTokenValue)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	mock := &mockReader{}
	cached := NewCachedClient(mock, WithCacheTTL(time.Hour))
	rateLimited := NewRateLimitedClient(mock, 100, 100)

	readers := []RepositoryReader{client, cached, rateLimited}
	for i, reader := range readers {
		if reader == nil {
			t.Fatalf("reader %d is nil", i)
		}
	}
}
