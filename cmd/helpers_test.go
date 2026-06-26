package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintJSON(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	testData := map[string]interface{}{
		"name":   "test",
		"count":  42,
		"active": true,
	}

	printJSON(testData)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

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
