// check-swagger-alignment.go - Compares implemented endpoints against Swagger spec
// Usage: go run scripts/check-swagger-alignment.go [options]
//
// This tool fetches the OpenAPI/Swagger spec from the Quay.io API and compares
// it against the endpoints implemented in the lib/ directory.

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// SwaggerSpec represents the OpenAPI/Swagger specification
type SwaggerSpec struct {
	Paths map[string]map[string]json.RawMessage `json:"paths"`
}

// Operation represents an API operation
type Operation struct {
	Tags        []string `json:"tags"`
	Summary     string   `json:"summary"`
	OperationID string   `json:"operationId"`
	Deprecated  bool     `json:"deprecated"`
	Description string   `json:"description"`
}

// Endpoint represents an API endpoint
type Endpoint struct {
	Method     string
	Path       string
	Tags       []string
	Summary    string
	Deprecated bool
}

// ImplementedEndpoint represents an endpoint found in source code
type ImplementedEndpoint struct {
	Method     string
	Path       string
	SourceFile string
	Function   string
}

func main() {
	swaggerURL := flag.String("swagger-url", "https://quay.io/api/v1/discovery", "URL to fetch Swagger spec")
	libPath := flag.String("lib-path", "./lib", "Path to lib directory")
	baseURLVar := flag.String("base-url-var", "QuayURL", "Variable name for base URL in source")
	outputFormat := flag.String("output", "text", "Output format: text, json, markdown")
	flag.Parse()

	// Fetch and parse swagger spec
	fmt.Println("Fetching Swagger spec from", *swaggerURL, "...")
	specEndpoints, err := fetchSwaggerEndpoints(*swaggerURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching swagger spec: %v\n", err)
		os.Exit(1)
	}

	// Parse implemented endpoints from source
	fmt.Println("Scanning lib/ for implemented endpoints...")
	implEndpoints, err := scanSourceEndpoints(*libPath, *baseURLVar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning source: %v\n", err)
		os.Exit(1)
	}

	// Compare and generate report
	report := generateReport(specEndpoints, implEndpoints)

	switch *outputFormat {
	case "json":
		printJSONReport(report)
	case "markdown":
		printMarkdownReport(report)
	default:
		printTextReport(report)
	}
}

func fetchSwaggerEndpoints(url string) ([]Endpoint, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var spec SwaggerSpec
	if err := json.NewDecoder(resp.Body).Decode(&spec); err != nil {
		return nil, fmt.Errorf("failed to parse swagger spec: %w", err)
	}

	var endpoints []Endpoint
	httpMethods := map[string]bool{
		"get": true, "post": true, "put": true, "delete": true, "patch": true,
	}

	for path, methods := range spec.Paths {
		for method, rawOp := range methods {
			methodLower := strings.ToLower(method)
			if !httpMethods[methodLower] {
				// Skip non-HTTP method fields (x-name, x-path, x-tag, parameters, etc.)
				continue
			}

			var op Operation
			if err := json.Unmarshal(rawOp, &op); err != nil {
				// Skip if we can't parse as an operation
				continue
			}

			endpoints = append(endpoints, Endpoint{
				Method:     strings.ToUpper(method),
				Path:       path,
				Tags:       op.Tags,
				Summary:    op.Summary,
				Deprecated: op.Deprecated,
			})
		}
	}

	sort.Slice(endpoints, func(i, j int) bool {
		if endpoints[i].Path == endpoints[j].Path {
			return endpoints[i].Method < endpoints[j].Method
		}
		return endpoints[i].Path < endpoints[j].Path
	})

	return endpoints, nil
}

func scanSourceEndpoints(libPath, baseURLVar string) ([]ImplementedEndpoint, error) {
	var endpoints []ImplementedEndpoint

	// Patterns to match endpoint definitions
	// Matches: QuayURL + "/path" or fmt.Sprintf("%s/path", QuayURL)
	urlConcatPattern := regexp.MustCompile(baseURLVar + `\s*\+\s*"(/[^"]+)"`)
	sprintfPattern := regexp.MustCompile(`fmt\.Sprintf\s*\(\s*"%s(/[^"]+)"`)
	// Match HTTP method from http.NewRequest or newRequest calls
	methodPattern := regexp.MustCompile(`(?:http\.NewRequest|newRequest|newRequestWithBody)\s*\(\s*"(GET|POST|PUT|DELETE|PATCH)"`)

	err := filepath.Walk(libPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var currentMethod string
		lineNum := 0
		var currentFunc string

		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			// Track current function name
			if funcMatch := regexp.MustCompile(`func\s+(?:\([^)]+\)\s+)?(\w+)`).FindStringSubmatch(line); len(funcMatch) > 1 {
				currentFunc = funcMatch[1]
			}

			// Look for HTTP method
			if methodMatch := methodPattern.FindStringSubmatch(line); len(methodMatch) > 1 {
				currentMethod = methodMatch[1]
			}

			// Look for URL patterns
			var urlPath string
			if matches := urlConcatPattern.FindStringSubmatch(line); len(matches) > 1 {
				urlPath = matches[1]
			} else if matches := sprintfPattern.FindStringSubmatch(line); len(matches) > 1 {
				urlPath = matches[1]
			}

			if urlPath != "" {
				// Normalize the path (replace %s, %d with {})
				normalizedPath := normalizePath(urlPath)
				method := currentMethod
				if method == "" {
					method = "GET" // Default assumption
				}

				endpoints = append(endpoints, ImplementedEndpoint{
					Method:     method,
					Path:       normalizedPath,
					SourceFile: filepath.Base(path),
					Function:   currentFunc,
				})
			}
		}

		return scanner.Err()
	})

	return endpoints, err
}

func normalizePath(path string) string {
	// Replace %s, %d, %v with {}
	path = regexp.MustCompile(`%[sdv]`).ReplaceAllString(path, "{}")
	// Remove query string patterns (e.g., ?{} or ?query={})
	if idx := strings.Index(path, "?"); idx != -1 {
		path = path[:idx]
	}
	// Remove trailing slashes for comparison
	path = strings.TrimSuffix(path, "/")
	return path
}

func normalizeSpecPath(path string) string {
	// Remove /api/v1 prefix if present
	path = strings.TrimPrefix(path, "/api/v1")
	// The Quay API spec uses {repository} to represent namespace/repo as a single parameter
	// but the implementation uses separate /{namespace}/{repo} parameters
	// Replace {repository} with {}/{} to match the implementation pattern
	path = strings.ReplaceAll(path, "{repository}", "{}/{}")
	// Replace remaining {param_name} with {}
	path = regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(path, "{}")
	// Remove trailing slashes
	path = strings.TrimSuffix(path, "/")
	return path
}

// Report represents the comparison report
type Report struct {
	TotalSpec         int
	TotalImplemented  int
	Implemented       []Endpoint
	Missing           []Endpoint
	MissingDeprecated []Endpoint
	Extra             []ImplementedEndpoint
	ByTag             map[string]TagReport
}

// TagReport represents coverage for a specific tag/category
type TagReport struct {
	Total       int
	Implemented int
}

func generateReport(specEndpoints []Endpoint, implEndpoints []ImplementedEndpoint) Report {
	report := Report{
		TotalSpec: len(specEndpoints),
		ByTag:     make(map[string]TagReport),
	}

	// Create a set of implemented paths for quick lookup
	implSet := make(map[string]bool)
	for _, impl := range implEndpoints {
		key := impl.Method + " " + impl.Path
		implSet[key] = true
	}

	// Create a set of spec paths for finding extras
	specSet := make(map[string]bool)

	for _, spec := range specEndpoints {
		normalizedSpecPath := normalizeSpecPath(spec.Path)
		key := spec.Method + " " + normalizedSpecPath
		specSet[key] = true

		// Update tag stats
		for _, tag := range spec.Tags {
			tr := report.ByTag[tag]
			tr.Total++
			if implSet[key] {
				tr.Implemented++
			}
			report.ByTag[tag] = tr
		}

		// Check if implemented
		if implSet[key] {
			report.Implemented = append(report.Implemented, spec)
			report.TotalImplemented++
		} else {
			if spec.Deprecated {
				report.MissingDeprecated = append(report.MissingDeprecated, spec)
			} else {
				report.Missing = append(report.Missing, spec)
			}
		}
	}

	// Find extra implementations (in source but not in spec)
	for _, impl := range implEndpoints {
		key := impl.Method + " " + impl.Path
		if !specSet[key] {
			report.Extra = append(report.Extra, impl)
		}
	}

	return report
}

func printTextReport(report Report) {
	fmt.Println()
	fmt.Println("=============================================================")
	fmt.Println("                   API COVERAGE REPORT")
	fmt.Println("              go-quay vs Quay.io API Spec")
	fmt.Println("=============================================================")
	fmt.Println()

	coverage := float64(report.TotalImplemented) / float64(report.TotalSpec) * 100
	fmt.Println("SUMMARY")
	fmt.Println("-------")
	fmt.Printf("Total API Endpoints:     %d\n", report.TotalSpec)
	fmt.Printf("Implemented:             %d (%.1f%%)\n", report.TotalImplemented, coverage)
	fmt.Printf("Missing:                 %d\n", len(report.Missing))
	fmt.Printf("Missing (Deprecated):    %d\n", len(report.MissingDeprecated))
	fmt.Printf("Extra (not in spec):     %d\n", len(report.Extra))
	fmt.Println()

	// By tag
	if len(report.ByTag) > 0 {
		fmt.Println("BY CATEGORY")
		fmt.Println("-----------")

		// Sort tags by name
		var tags []string
		for tag := range report.ByTag {
			tags = append(tags, tag)
		}
		sort.Strings(tags)

		for _, tag := range tags {
			tr := report.ByTag[tag]
			pct := float64(tr.Implemented) / float64(tr.Total) * 100
			bar := progressBar(pct, 20)
			fmt.Printf("%-20s %s %2d/%-3d %5.1f%%\n", tag, bar, tr.Implemented, tr.Total, pct)
		}
		fmt.Println()
	}

	// Missing endpoints
	if len(report.Missing) > 0 {
		fmt.Println("MISSING ENDPOINTS")
		fmt.Println("-----------------")
		for _, ep := range report.Missing {
			fmt.Printf("- %-6s %s\n", ep.Method, ep.Path)
		}
		fmt.Println()
	}

	// Missing deprecated
	if len(report.MissingDeprecated) > 0 {
		fmt.Println("MISSING ENDPOINTS (Deprecated - Low Priority)")
		fmt.Println("----------------------------------------------")
		for _, ep := range report.MissingDeprecated {
			fmt.Printf("- %-6s %s\n", ep.Method, ep.Path)
		}
		fmt.Println()
	}

	// Extra implementations
	if len(report.Extra) > 0 {
		fmt.Println("EXTRA IMPLEMENTATIONS (Not in spec)")
		fmt.Println("------------------------------------")
		for _, ep := range report.Extra {
			fmt.Printf("- %-6s %-40s (%s)\n", ep.Method, ep.Path, ep.Function)
		}
		fmt.Println()
	}

	fmt.Println("=============================================================")
}

func progressBar(pct float64, width int) string {
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	return "[" + strings.Repeat("=", filled) + strings.Repeat(" ", width-filled) + "]"
}

func printJSONReport(report Report) {
	data, _ := json.MarshalIndent(report, "", "  ")
	fmt.Println(string(data))
}

func printMarkdownReport(report Report) {
	coverage := float64(report.TotalImplemented) / float64(report.TotalSpec) * 100

	fmt.Println("# API Coverage Report")
	fmt.Println()
	fmt.Println("## Summary")
	fmt.Println()
	fmt.Printf("| Metric | Value |\n")
	fmt.Printf("|--------|-------|\n")
	fmt.Printf("| Total API Endpoints | %d |\n", report.TotalSpec)
	fmt.Printf("| Implemented | %d (%.1f%%) |\n", report.TotalImplemented, coverage)
	fmt.Printf("| Missing | %d |\n", len(report.Missing))
	fmt.Printf("| Missing (Deprecated) | %d |\n", len(report.MissingDeprecated))
	fmt.Println()

	if len(report.Missing) > 0 {
		fmt.Println("## Missing Endpoints")
		fmt.Println()
		fmt.Println("| Method | Path |")
		fmt.Println("|--------|------|")
		for _, ep := range report.Missing {
			fmt.Printf("| %s | %s |\n", ep.Method, ep.Path)
		}
		fmt.Println()
	}
}
