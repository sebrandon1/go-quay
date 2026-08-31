package cmd

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

const testPrintFieldValue = "sample"

func writeResponse(t *testing.T, w http.ResponseWriter, body []byte) {
	t.Helper()
	if _, err := w.Write(body); err != nil {
		t.Errorf("response write: %v", err)
	}
}

func copyStdout(t *testing.T, r io.Reader) string {
	t.Helper()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}
	return buf.String()
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w
	fn()
	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	os.Stdout = oldStdout
	return copyStdout(t, r)
}

func TestPrintJSON(t *testing.T) {
	testData := map[string]interface{}{
		"name":   testPrintFieldValue,
		"count":  42,
		"active": true,
	}

	output := captureStdout(t, func() {
		if err := printJSON(testData); err != nil {
			t.Fatalf("printJSON: %v", err)
		}
	})

	if !strings.Contains(output, `"name": "`+testPrintFieldValue+`"`) {
		t.Errorf("Expected name field in output, got: %s", output)
	}
	if !strings.Contains(output, `"count": 42`) {
		t.Errorf("Expected count field in output, got: %s", output)
	}
	if !strings.Contains(output, `"active": true`) {
		t.Errorf("Expected active field in output, got: %s", output)
	}

	lines := strings.Split(output, "\n")
	foundIndented := false
	for _, line := range lines {
		if strings.HasPrefix(line, "  ") {
			foundIndented = true
			break
		}
	}
	if !foundIndented {
		t.Error("Expected indented JSON output")
	}
}

func TestPrintJSONYAML(t *testing.T) {
	orig := outputFormat
	outputFormat = outputYAML
	t.Cleanup(func() { outputFormat = orig })

	output := captureStdout(t, func() {
		if err := printJSON(map[string]string{"field": testPrintFieldValue}); err != nil {
			t.Fatalf("printJSON: %v", err)
		}
	})

	if !strings.Contains(output, "field: "+testPrintFieldValue) {
		t.Errorf("expected YAML output, got: %s", output)
	}
	if strings.Contains(output, `"name"`) {
		t.Errorf("expected YAML not JSON, got: %s", output)
	}
}

func TestPrintJSONTableFallback(t *testing.T) {
	orig := outputFormat
	outputFormat = outputTable
	t.Cleanup(func() { outputFormat = orig })

	output := captureStdout(t, func() {
		if err := printJSON(map[string]string{"name": testPrintFieldValue}); err != nil {
			t.Fatalf("printJSON: %v", err)
		}
	})

	if !strings.Contains(output, `"name": "`+testPrintFieldValue+`"`) {
		t.Errorf("expected JSON fallback for table mode, got: %s", output)
	}
}

func TestPrintAsJSONMarshalError(t *testing.T) {
	err := printAsJSON(make(chan int))
	if err == nil {
		t.Fatal("expected marshal error for channel type")
	}
	if !strings.Contains(err.Error(), "marshaling JSON") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetClient(t *testing.T) {
	origToken := token
	origURL := quayURL
	origVersion := rootCmd.Version
	t.Cleanup(func() {
		token = origToken
		quayURL = origURL
		rootCmd.Version = origVersion
	})

	token = "test-token"
	quayURL = "https://custom.example.com/api/v1"
	rootCmd.Version = "v9.9.9"

	client, err := getClient()
	if err != nil {
		t.Fatalf("getClient: %v", err)
	}
	if client.BearerToken != token {
		t.Errorf("BearerToken = %q, want %q", client.BearerToken, token)
	}
	if client.BaseURL != quayURL {
		t.Errorf("BaseURL = %q, want %q", client.BaseURL, quayURL)
	}
	if client.Version != "v9.9.9" {
		t.Errorf("Version = %q, want v9.9.9", client.Version)
	}
}
