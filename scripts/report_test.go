package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

const (
	repositoriesTag = "repositories"
	missingAPIPath  = "/api/v1/missing"
)

func TestNormalizePaths(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "formats placeholders and removes query and trailing slash", path: "/repo/%s/tag/%d/%v/?page=%s", want: "/repo/{}/tag/{}/{}"},
		{name: "removes query before trailing slash handling", path: "/repo/?page=1", want: "/repo"},
		{name: "keeps paths without placeholders", path: "/repo/tag", want: "/repo/tag"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizePath(tt.path); got != tt.want {
				t.Errorf("normalizePath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestNormalizeSpecPath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "repository and named parameters", path: "/api/v1/repository/{repository}/tag/{tag}/", want: "/repository/{}/{}/tag/{}"},
		{name: "prefix absent", path: "/organization/{org}/repositories", want: "/organization/{}/repositories"},
		{name: "literal path", path: "/api/v1/status", want: "/status"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeSpecPath(tt.path); got != tt.want {
				t.Errorf("normalizeSpecPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestGenerateReport(t *testing.T) {
	specEndpoints := []Endpoint{
		{Method: http.MethodGet, Path: "/api/v1/repository/{repository}/", Tags: []string{repositoriesTag}, Summary: "List"},
		{Method: http.MethodPost, Path: "/api/v1/repository/{repository}/tag/{tag}", Tags: []string{repositoriesTag}, Summary: "Change"},
		{Method: http.MethodGet, Path: missingAPIPath, Tags: []string{"security"}, Summary: "Missing"},
		{Method: http.MethodGet, Path: "/api/v1/old", Tags: []string{"legacy"}, Summary: "Deprecated", Deprecated: true},
	}
	implEndpoints := []ImplementedEndpoint{
		{Method: http.MethodGet, Path: "/repository/{}/{}", Function: "ListRepositories"},
		{Method: http.MethodPost, Path: "/repository/{}/{}/tag/{}", Function: "ChangeTag"},
		{Method: http.MethodGet, Path: "/extra", Function: "ExtraEndpoint"},
	}

	got := generateReport(specEndpoints, implEndpoints)
	if got.TotalSpec != 4 || got.TotalImplemented != 2 {
		t.Errorf("report totals = (%d, %d), want (4, 2)", got.TotalSpec, got.TotalImplemented)
	}
	if len(got.Implemented) != 2 || len(got.Missing) != 1 || len(got.MissingDeprecated) != 1 || len(got.Extra) != 1 {
		t.Errorf("report list lengths = implemented %d, missing %d, deprecated %d, extra %d; want 2, 1, 1, 1", len(got.Implemented), len(got.Missing), len(got.MissingDeprecated), len(got.Extra))
	}
	if len(got.Missing) == 1 && len(got.MissingDeprecated) == 1 && len(got.Extra) == 1 &&
		(got.Missing[0].Path != missingAPIPath || got.MissingDeprecated[0].Path != "/api/v1/old" || got.Extra[0].Function != "ExtraEndpoint") {
		t.Errorf("report unmatched endpoints = missing %#v, deprecated %#v, extra %#v", got.Missing, got.MissingDeprecated, got.Extra)
	}
	wantByTag := map[string]TagReport{
		repositoriesTag: {Total: 2, Implemented: 2},
		"security":      {Total: 1},
		"legacy":        {Total: 1},
	}
	if !reflect.DeepEqual(got.ByTag, wantByTag) {
		t.Errorf("report.ByTag = %#v, want %#v", got.ByTag, wantByTag)
	}
}

func TestPrintTextReport(t *testing.T) {
	report := generateReport(
		[]Endpoint{
			{Method: http.MethodGet, Path: "/api/v1/implemented", Tags: []string{"zeta"}},
			{Method: http.MethodGet, Path: missingAPIPath, Tags: []string{"alpha"}},
			{Method: http.MethodGet, Path: "/api/v1/deprecated", Deprecated: true},
		},
		[]ImplementedEndpoint{
			{Method: http.MethodGet, Path: "/implemented"},
			{Method: http.MethodPost, Path: "/extra", Function: "CreateExtra"},
		},
	)

	got := captureReportOutput(t, func() { printTextReport(report) })
	for _, want := range []string{
		"API COVERAGE REPORT",
		"Total API Endpoints:     3",
		"Implemented:             1 (33.3%)",
		"alpha",
		"zeta",
		"MISSING ENDPOINTS",
		"GET    /api/v1/missing",
		"MISSING ENDPOINTS (Deprecated - Low Priority)",
		"GET    /api/v1/deprecated",
		"EXTRA IMPLEMENTATIONS (Not in spec)",
		"POST   /extra",
		"CreateExtra",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("printTextReport() output missing %q:\n%s", want, got)
		}
	}
}

func TestPrintJSONReport(t *testing.T) {
	want := Report{
		TotalSpec:        2,
		TotalImplemented: 1,
		Implemented:      []Endpoint{{Method: http.MethodGet, Path: "/implemented"}},
		Missing:          []Endpoint{{Method: http.MethodGet, Path: "/missing"}},
		ByTag:            map[string]TagReport{repositoriesTag: {Total: 2, Implemented: 1}},
	}

	gotOutput := captureReportOutput(t, func() { printJSONReport(want) })
	var got Report
	if err := json.Unmarshal([]byte(gotOutput), &got); err != nil {
		t.Fatalf("printJSONReport() emitted invalid JSON: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("printJSONReport() = %#v, want %#v", got, want)
	}
}

func TestPrintMarkdownReport(t *testing.T) {
	tests := []struct {
		name               string
		in                 Report
		want               []string
		wantMissingSection bool
	}{
		{
			name: "includes missing endpoint table",
			in: Report{
				TotalSpec:        2,
				TotalImplemented: 1,
				Missing:          []Endpoint{{Method: http.MethodGet, Path: "/missing"}},
			},
			want:               []string{"# API Coverage Report", "50.0%", "## Missing Endpoints", "| GET | /missing |"},
			wantMissingSection: true,
		},
		{
			name: "omits empty missing endpoint table",
			in:   Report{TotalSpec: 1, TotalImplemented: 1},
			want: []string{"100.0%", "| Missing | 0 |"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureReportOutput(t, func() { printMarkdownReport(tt.in) })
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("printMarkdownReport() output missing %q:\n%s", want, got)
				}
			}
			if hasMissingSection := strings.Contains(got, "## Missing Endpoints"); hasMissingSection != tt.wantMissingSection {
				t.Errorf("printMarkdownReport() missing section present = %t, want %t:\n%s", hasMissingSection, tt.wantMissingSection, got)
			}
		})
	}
}

func TestProgressBar(t *testing.T) {
	tests := []struct {
		name  string
		pct   float64
		width int
		want  string
	}{
		{name: "empty", pct: 0, width: 4, want: "[    ]"},
		{name: "partial", pct: 50, width: 4, want: "[==  ]"},
		{name: "full", pct: 100, width: 4, want: "[====]"},
		{name: "clamps above full", pct: 150, width: 4, want: "[====]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := progressBar(tt.pct, tt.width); got != tt.want {
				t.Errorf("progressBar(%v, %d) = %q, want %q", tt.pct, tt.width, got, tt.want)
			}
		})
	}
}

func captureReportOutput(t *testing.T, print func()) string {
	t.Helper()

	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	os.Stdout = writer
	t.Cleanup(func() {
		os.Stdout = oldStdout
		_ = writer.Close()
		_ = reader.Close()
	})

	print()
	if err := writer.Close(); err != nil {
		t.Fatalf("close stdout pipe: %v", err)
	}
	os.Stdout = oldStdout
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read stdout pipe: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close stdout reader: %v", err)
	}
	return string(output)
}
