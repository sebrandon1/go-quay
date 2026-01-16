// Example: Security Scanning with go-quay Library
//
// This example demonstrates how to use go-quay for security vulnerability scanning:
// - Fetching manifest information for images
// - Retrieving security scan results
// - Analyzing vulnerability data
// - Building a simple security report
//
// Usage:
//
//	export QUAY_TOKEN="your-quay-api-token"
//	go run main.go --namespace myorg --repository myapp --tag latest
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sebrandon1/go-quay/lib"
)

func main() {
	// Parse command line arguments
	namespace := flag.String("namespace", "", "Repository namespace (required)")
	repository := flag.String("repository", "", "Repository name (required)")
	tag := flag.String("tag", "latest", "Image tag to scan")
	flag.Parse()

	if *namespace == "" || *repository == "" {
		fmt.Println("Usage: go run main.go --namespace <namespace> --repository <repository> [--tag <tag>]")
		os.Exit(1)
	}

	// Get authentication token
	token := os.Getenv("QUAY_TOKEN")
	if token == "" {
		log.Fatal("QUAY_TOKEN environment variable is required")
	}

	// Initialize client
	client, err := lib.NewClient(token)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("=== Security Scan Report ===")
	fmt.Printf("Repository: %s/%s\n", *namespace, *repository)
	fmt.Printf("Tag: %s\n\n", *tag)

	// Step 1: Get repository and tag information
	fmt.Println("1. Fetching repository information...")
	repo, err := client.GetRepository(*namespace, *repository)
	if err != nil {
		log.Fatalf("Failed to get repository: %v", err)
	}

	// Find the specified tag
	var manifestDigest string
	for _, t := range repo.Tags.Tags {
		if t.Name == *tag {
			manifestDigest = t.ManifestDigest
			fmt.Printf("   Found tag '%s' with manifest: %s\n", t.Name, shortenDigest(manifestDigest))
			break
		}
	}

	if manifestDigest == "" {
		log.Fatalf("Tag '%s' not found in repository", *tag)
	}

	// Step 2: Get manifest details
	fmt.Println("\n2. Fetching manifest details...")
	manifest, err := client.GetManifest(*namespace, *repository, manifestDigest)
	if err != nil {
		log.Fatalf("Failed to get manifest: %v", err)
	}
	fmt.Printf("   Digest: %s\n", shortenDigest(manifest.Digest))
	fmt.Printf("   Layers: %d\n", len(manifest.Layers))

	// Step 3: Get security scan results
	fmt.Println("\n3. Fetching security scan results...")
	security, err := client.GetManifestSecurity(*namespace, *repository, manifestDigest, true)
	if err != nil {
		log.Fatalf("Failed to get security scan: %v", err)
	}

	// Analyze scan status
	fmt.Printf("   Scan Status: %s\n", security.Status)

	if security.Status == "queued" || security.Status == "scanning" {
		fmt.Println("\n   Security scan is still in progress. Please try again later.")
		return
	}

	if security.Status == "unsupported" {
		fmt.Println("\n   This image type is not supported for security scanning.")
		return
	}

	if security.Status != "scanned" {
		fmt.Printf("\n   Unexpected scan status: %s\n", security.Status)
		return
	}

	// Step 4: Analyze vulnerabilities
	fmt.Println("\n4. Vulnerability Analysis")
	fmt.Println("   " + strings.Repeat("-", 50))

	// Count vulnerabilities by severity
	vulnCounts := map[string]int{
		"Critical":   0,
		"High":       0,
		"Medium":     0,
		"Low":        0,
		"Negligible": 0,
		"Unknown":    0,
	}

	// Access vulnerability data from the security response
	if security.Data != nil && security.Data.Layer != nil {
		for _, feature := range security.Data.Layer.Features {
			for _, vuln := range feature.Vulnerabilities {
				severity := vuln.Severity
				if severity == "" {
					severity = "Unknown"
				}
				vulnCounts[severity]++
			}
		}
	}

	totalVulns := 0
	for _, count := range vulnCounts {
		totalVulns += count
	}

	fmt.Printf("   Total Vulnerabilities: %d\n\n", totalVulns)
	fmt.Println("   Breakdown by Severity:")
	fmt.Printf("   - Critical:   %d\n", vulnCounts["Critical"])
	fmt.Printf("   - High:       %d\n", vulnCounts["High"])
	fmt.Printf("   - Medium:     %d\n", vulnCounts["Medium"])
	fmt.Printf("   - Low:        %d\n", vulnCounts["Low"])
	fmt.Printf("   - Negligible: %d\n", vulnCounts["Negligible"])
	fmt.Printf("   - Unknown:    %d\n", vulnCounts["Unknown"])

	// Step 5: Show critical/high vulnerabilities (limited output)
	if vulnCounts["Critical"] > 0 || vulnCounts["High"] > 0 {
		fmt.Println("\n5. Critical and High Severity Vulnerabilities")
		fmt.Println("   " + strings.Repeat("-", 50))

		shown := 0
		maxShow := 10

		if security.Data != nil && security.Data.Layer != nil {
			for _, feature := range security.Data.Layer.Features {
				for _, vuln := range feature.Vulnerabilities {
					if vuln.Severity == "Critical" || vuln.Severity == "High" {
						if shown < maxShow {
							fmt.Printf("\n   [%s] %s\n", vuln.Severity, vuln.Name)
							fmt.Printf("   Package: %s %s\n", feature.Name, feature.Version)
							if vuln.FixedBy != "" {
								fmt.Printf("   Fixed in: %s\n", vuln.FixedBy)
							}
							if vuln.Link != "" {
								fmt.Printf("   Details: %s\n", vuln.Link)
							}
							shown++
						}
					}
				}
			}
		}

		if shown < vulnCounts["Critical"]+vulnCounts["High"] {
			fmt.Printf("\n   ... and %d more high/critical vulnerabilities\n",
				vulnCounts["Critical"]+vulnCounts["High"]-shown)
		}
	}

	// Summary
	fmt.Println("\n=== Scan Summary ===")
	if vulnCounts["Critical"] > 0 {
		fmt.Println("CRITICAL: This image has critical vulnerabilities that should be addressed immediately!")
	} else if vulnCounts["High"] > 0 {
		fmt.Println("WARNING: This image has high severity vulnerabilities.")
	} else if totalVulns > 0 {
		fmt.Println("INFO: This image has some vulnerabilities but no critical/high severity issues.")
	} else {
		fmt.Println("OK: No vulnerabilities detected in this image.")
	}
}

// shortenDigest returns a shortened version of a digest for display
func shortenDigest(digest string) string {
	if len(digest) > 19 {
		return digest[:19] + "..."
	}
	return digest
}
