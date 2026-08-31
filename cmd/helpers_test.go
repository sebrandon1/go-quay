package cmd

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

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
		"name":   "test",
		"count":  42,
		"active": true,
	}

	output := captureStdout(t, func() {
		if err := printJSON(testData); err != nil {
			t.Fatalf("printJSON: %v", err)
		}
	})

	// Verify JSON formatting
	if !strings.Contains(output, `"name": "test"`) {
		t.Errorf("Expected name field in output, got: %s", output)
	}
	if !strings.Contains(output, `"count": 42`) {
		t.Errorf("Expected count field in output, got: %s", output)
	}
	if !strings.Contains(output, `"active": true`) {
		t.Errorf("Expected active field in output, got: %s", output)
	}

	// Verify it's indented (has leading spaces)
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
