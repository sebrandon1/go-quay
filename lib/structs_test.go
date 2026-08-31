package lib

import (
	"encoding/json"
	"strings"
	"testing"
)

func mustMarshalJSON(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	return data
}

func TestQuayErrorErrorWithoutDetail(t *testing.T) {
	err := &QuayError{Status: 404, Message: "not_found"}
	want := "quay API error (status 404): not_found"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestQuayErrorErrorWithDetail(t *testing.T) {
	err := &QuayError{Status: 403, Message: "denied", Detail: "insufficient permissions"}
	got := err.Error()
	if !strings.Contains(got, "403") || !strings.Contains(got, "denied") || !strings.Contains(got, "insufficient permissions") {
		t.Errorf("Error() = %q, want status, message, and detail", got)
	}
}

func TestQuayErrorStatusCode(t *testing.T) {
	err := &QuayError{Status: 429}
	if got := err.StatusCode(); got != 429 {
		t.Errorf("StatusCode() = %d, want 429", got)
	}
}

func roundTripJSON(t *testing.T, src, dst any) {
	t.Helper()
	data := mustMarshalJSON(t, src)
	if err := json.Unmarshal(data, dst); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
}

func TestRepositoryWithTagsJSONRoundTrip(t *testing.T) {
	src := RepositoryWithTags{
		Repository: Repository{
			Name:      testRepoName,
			Namespace: testNamespace,
			IsPublic:  true,
		},
		Tags: RepositoryTags{
			Tags: []Tag{{Name: testTagNameLatest, ManifestDigest: testDigestSHA256}},
		},
	}
	dst := &RepositoryWithTags{}
	roundTripJSON(t, src, dst)

	if dst.Name != testRepoName || dst.Namespace != testNamespace || !dst.IsPublic {
		t.Errorf("repository fields not preserved: %+v", dst.Repository)
	}
	if len(dst.Tags.Tags) != 1 || dst.Tags.Tags[0].Name != testTagNameLatest {
		t.Errorf("tags not preserved: %+v", dst.Tags)
	}
}

func TestOrganizationMembersJSONRoundTrip(t *testing.T) {
	src := OrganizationMembers{
		Members: []OrganizationMember{
			{Name: "alice", Kind: testKindUser},
			{Name: "robots", Kind: testKindTeam},
		},
	}
	dst := &OrganizationMembers{}
	roundTripJSON(t, src, dst)

	if len(dst.Members) != 2 {
		t.Fatalf("members len = %d, want 2", len(dst.Members))
	}
	if dst.Members[0].Name != "alice" {
		t.Errorf("first member = %q, want alice", dst.Members[0].Name)
	}
}

func TestSecurityScanJSONRoundTrip(t *testing.T) {
	src := SecurityScan{
		Status: testSecScanStatusQueued,
		Data: &SecurityData{
			Layer: &SecurityLayer{Name: "layer-1", NamespaceName: testNamespace},
		},
	}
	dst := &SecurityScan{}
	roundTripJSON(t, src, dst)

	if dst.Status != testSecScanStatusQueued {
		t.Errorf("status = %q, want %s", dst.Status, testSecScanStatusQueued)
	}
	if dst.Data == nil || dst.Data.Layer == nil || dst.Data.Layer.Name != "layer-1" {
		t.Errorf("nested layer not preserved: %+v", dst.Data)
	}
}

func TestQuayErrorJSONRoundTrip(t *testing.T) {
	src := QuayError{
		Status:      400,
		Message:     "invalid",
		Detail:      "bad request",
		ErrorDetail: map[string]any{"field": "name"},
	}
	dst := &QuayError{}
	roundTripJSON(t, src, dst)

	if dst.Status != 400 || dst.Message != "invalid" || dst.Detail != "bad request" {
		t.Errorf("QuayError fields not preserved: %+v", dst)
	}
	if dst.ErrorDetail["field"] != "name" {
		t.Errorf("ErrorDetail = %v, want field=name", dst.ErrorDetail)
	}
}

func TestStructOmitempty(t *testing.T) {
	repo := Repository{Name: testRepoName}
	data := mustMarshalJSON(t, repo)
	jsonStr := string(data)

	if strings.Contains(jsonStr, `"namespace"`) {
		t.Errorf("expected empty namespace omitted, got %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"name":"`+testRepoName+`"`) && !strings.Contains(jsonStr, `"name": "`+testRepoName+`"`) {
		t.Errorf("expected name in JSON, got %s", jsonStr)
	}
}
