# Security Scanning with go-quay

This tutorial covers how to use the go-quay library to scan container images for security vulnerabilities.

## Prerequisites

- Completed [Getting Started](./01-getting-started.md)
- A repository with at least one pushed image
- Security scanning enabled on your Quay.io registry

## Understanding Security Scanning

Quay.io uses Clair to scan container images for vulnerabilities. When you push an image:

1. Quay queues the image for scanning
2. Clair analyzes each layer
3. Results are stored and accessible via API

## Scan Status Values

| Status | Description |
|--------|-------------|
| `scanned` | Scan complete, results available |
| `queued` | Scan is waiting to start |
| `scanning` | Scan is in progress |
| `unsupported` | Image type cannot be scanned |
| `failed` | Scan encountered an error |

## Getting Security Scan Results

### Basic Security Scan

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/sebrandon1/go-quay/lib"
)

func main() {
    client, _ := lib.NewClient(os.Getenv("QUAY_TOKEN"))

    // First, get the manifest digest for a tag
    repo, err := client.GetRepository("my-namespace", "my-app")
    if err != nil {
        log.Fatalf("Failed to get repository: %v", err)
    }

    // Find the 'latest' tag
    var manifestDigest string
    for _, tag := range repo.Tags.Tags {
        if tag.Name == "latest" {
            manifestDigest = tag.ManifestDigest
            break
        }
    }

    if manifestDigest == "" {
        log.Fatal("Tag 'latest' not found")
    }

    // Get security scan results
    security, err := client.GetManifestSecurity(
        "my-namespace",
        "my-app",
        manifestDigest,
        true, // include vulnerability details
    )
    if err != nil {
        log.Fatalf("Failed to get security info: %v", err)
    }

    fmt.Printf("Scan Status: %s\n", security.Status)
}
```

### Analyzing Vulnerability Data

```go
// Get security scan with vulnerability details
security, err := client.GetManifestSecurity(namespace, repo, digest, true)
if err != nil {
    log.Fatalf("Failed to get security scan: %v", err)
}

if security.Status != "scanned" {
    fmt.Printf("Scan not ready: %s\n", security.Status)
    return
}

// Count vulnerabilities by severity
vulnCounts := map[string]int{
    "Critical":   0,
    "High":       0,
    "Medium":     0,
    "Low":        0,
    "Negligible": 0,
    "Unknown":    0,
}

// Iterate through features (packages) and their vulnerabilities
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

// Print summary
fmt.Println("Vulnerability Summary:")
fmt.Printf("  Critical:   %d\n", vulnCounts["Critical"])
fmt.Printf("  High:       %d\n", vulnCounts["High"])
fmt.Printf("  Medium:     %d\n", vulnCounts["Medium"])
fmt.Printf("  Low:        %d\n", vulnCounts["Low"])
fmt.Printf("  Negligible: %d\n", vulnCounts["Negligible"])
```

### Detailed Vulnerability Report

```go
// Generate detailed report for critical/high vulnerabilities
if security.Data != nil && security.Data.Layer != nil {
    fmt.Println("\nCritical and High Severity Vulnerabilities:")
    fmt.Println("============================================")

    for _, feature := range security.Data.Layer.Features {
        for _, vuln := range feature.Vulnerabilities {
            if vuln.Severity == "Critical" || vuln.Severity == "High" {
                fmt.Printf("\n[%s] %s\n", vuln.Severity, vuln.Name)
                fmt.Printf("  Package: %s %s\n", feature.Name, feature.Version)

                if vuln.FixedBy != "" {
                    fmt.Printf("  Fixed in: %s\n", vuln.FixedBy)
                } else {
                    fmt.Println("  Fixed in: No fix available")
                }

                if vuln.Description != "" {
                    // Truncate long descriptions
                    desc := vuln.Description
                    if len(desc) > 200 {
                        desc = desc[:200] + "..."
                    }
                    fmt.Printf("  Description: %s\n", desc)
                }

                if vuln.Link != "" {
                    fmt.Printf("  Details: %s\n", vuln.Link)
                }
            }
        }
    }
}
```

## Building a Security Dashboard

Here's how to build a simple security dashboard for multiple images:

```go
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/sebrandon1/go-quay/lib"
)

type ImageSecurity struct {
    Name     string
    Tag      string
    Status   string
    Critical int
    High     int
    Medium   int
    Low      int
}

func main() {
    client, _ := lib.NewClient(os.Getenv("QUAY_TOKEN"))
    namespace := "my-org"

    // Images to scan
    images := []struct {
        repo string
        tag  string
    }{
        {"app-frontend", "latest"},
        {"app-backend", "latest"},
        {"app-worker", "latest"},
    }

    var results []ImageSecurity

    for _, img := range images {
        result := scanImage(client, namespace, img.repo, img.tag)
        results = append(results, result)
    }

    // Print dashboard
    fmt.Println("Security Dashboard")
    fmt.Println("==================")
    fmt.Printf("%-20s %-10s %-10s %8s %8s %8s %8s\n",
        "Image", "Tag", "Status", "Critical", "High", "Medium", "Low")
    fmt.Println(strings.Repeat("-", 80))

    for _, r := range results {
        fmt.Printf("%-20s %-10s %-10s %8d %8d %8d %8d\n",
            r.Name, r.Tag, r.Status, r.Critical, r.High, r.Medium, r.Low)
    }
}

func scanImage(client *lib.Client, namespace, repo, tag string) ImageSecurity {
    result := ImageSecurity{Name: repo, Tag: tag}

    // Get repository to find manifest digest
    repoInfo, err := client.GetRepository(namespace, repo)
    if err != nil {
        result.Status = "error"
        return result
    }

    // Find tag
    var digest string
    for _, t := range repoInfo.Tags.Tags {
        if t.Name == tag {
            digest = t.ManifestDigest
            break
        }
    }

    if digest == "" {
        result.Status = "not found"
        return result
    }

    // Get security scan
    security, err := client.GetManifestSecurity(namespace, repo, digest, true)
    if err != nil {
        result.Status = "error"
        return result
    }

    result.Status = security.Status

    if security.Status == "scanned" && security.Data != nil && security.Data.Layer != nil {
        for _, feature := range security.Data.Layer.Features {
            for _, vuln := range feature.Vulnerabilities {
                switch vuln.Severity {
                case "Critical":
                    result.Critical++
                case "High":
                    result.High++
                case "Medium":
                    result.Medium++
                case "Low":
                    result.Low++
                }
            }
        }
    }

    return result
}
```

## Automated Security Alerts

Integrate security scanning into your CI/CD pipeline:

```go
package main

import (
    "fmt"
    "os"

    "github.com/sebrandon1/go-quay/lib"
)

func main() {
    client, _ := lib.NewClient(os.Getenv("QUAY_TOKEN"))

    // Configuration
    namespace := os.Getenv("QUAY_NAMESPACE")
    repo := os.Getenv("QUAY_REPOSITORY")
    tag := os.Getenv("IMAGE_TAG")
    maxCritical := 0  // Fail if any critical
    maxHigh := 5      // Fail if more than 5 high

    // Scan the image
    result := scanImage(client, namespace, repo, tag)

    // Check thresholds
    exitCode := 0

    if result.Critical > maxCritical {
        fmt.Printf("FAIL: Found %d critical vulnerabilities (max: %d)\n",
            result.Critical, maxCritical)
        exitCode = 1
    }

    if result.High > maxHigh {
        fmt.Printf("FAIL: Found %d high vulnerabilities (max: %d)\n",
            result.High, maxHigh)
        exitCode = 1
    }

    if exitCode == 0 {
        fmt.Println("PASS: Image meets security requirements")
    }

    os.Exit(exitCode)
}
```

## Working with Manifest Labels

Add security metadata to your images:

```go
// Add a security scan label to a manifest
err := client.AddManifestLabel(
    "my-namespace",
    "my-app",
    manifestDigest,
    "security-scanned",
    "true",
    "",
)
if err != nil {
    log.Printf("Failed to add label: %v", err)
}

// List all labels on a manifest
labels, err := client.GetManifestLabels("my-namespace", "my-app", manifestDigest)
if err != nil {
    log.Fatalf("Failed to get labels: %v", err)
}

fmt.Println("Manifest Labels:")
for _, label := range labels.Labels {
    fmt.Printf("  %s = %s\n", label.Key, label.Value)
}
```

## Complete Example

See the [security-scan example](../../examples/security-scan/main.go) for a complete working program.

## Next Steps

- [CI/CD Automation](./04-ci-cd-automation.md) - Automate security scanning in pipelines
- [Organization Administration](./05-organization-admin.md) - Set up security policies
