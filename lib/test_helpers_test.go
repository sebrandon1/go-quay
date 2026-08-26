package lib

import (
	"context"
	"sync/atomic"
)

type mockReader struct {
	calls atomic.Int32
	repo  RepositoryWithTags
	err   error
}

func (m *mockReader) GetRepository(_ context.Context, _, _ string) (RepositoryWithTags, error) {
	m.calls.Add(1)
	return m.repo, m.err
}

func (m *mockReader) ListTags(_ context.Context, _, _ string, _ int, _ bool) (*RepositoryTags, error) {
	m.calls.Add(1)
	return &RepositoryTags{}, m.err
}

func (m *mockReader) ListAllTags(_ context.Context, _, _ string, _ bool) ([]Tag, error) {
	m.calls.Add(1)
	return nil, m.err
}
