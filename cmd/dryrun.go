package cmd

import (
	"fmt"
	"net/http"
)

// dryRunTransport blocks API calls and reports their method and URL through the
// error returned by http.Client.Do. Request bodies are never read or printed.
type dryRunTransport struct{}

func (dryRunTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	bodySummary := "no request body"
	if req.Body != nil && req.Body != http.NoBody {
		bodySummary = "request body omitted"
		if req.ContentLength > 0 {
			bodySummary = fmt.Sprintf("request body omitted (%d bytes)", req.ContentLength)
		}
	}
	return nil, fmt.Errorf("dry-run: would send %s; request not sent; %s", req.Method, bodySummary)
}
