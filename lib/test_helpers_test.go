package lib

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"testing"
)

func mustMarshal(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	return data
}

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

func (m *mockReader) GetManifestSecurity(_ context.Context, _, _, _ string, _ bool) (*SecurityScan, error) {
	m.calls.Add(1)
	return nil, m.err
}

func (m *mockReader) GetManifest(_ context.Context, _, _, _ string) (*Manifest, error) {
	m.calls.Add(1)
	return nil, m.err
}

func (m *mockReader) GetManifestLabels(_ context.Context, _, _, _ string) (*ManifestLabels, error) {
	m.calls.Add(1)
	return nil, m.err
}
