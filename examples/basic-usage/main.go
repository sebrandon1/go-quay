// Example: Basic Usage of go-quay Library
//
// This example demonstrates fundamental operations with the go-quay library:
// - Initializing the client with authentication
// - Getting user information
// - Listing repositories
// - Getting repository details with tags
//
// Usage:
//
//	export QUAY_TOKEN="your-quay-api-token"
//	go run main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/sebrandon1/go-quay/lib"
)

func main() {
	// Step 1: Get authentication token from environment
	token := os.Getenv("QUAY_TOKEN")
	if token == "" {
		log.Fatal("QUAY_TOKEN environment variable is required")
	}

	// Step 2: Initialize the client
	client, err := lib.NewClient(token)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Println("=== go-quay Basic Usage Example ===\n")

	// Step 3: Get current user information
	fmt.Println("1. Getting current user information...")
	user, err := client.GetUser()
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	fmt.Printf("   Logged in as: %s\n", user.Username)
	fmt.Printf("   Email: %s\n", user.Email)
	fmt.Printf("   Verified: %v\n\n", user.Verified)

	// Step 4: List starred repositories
	fmt.Println("2. Getting starred repositories...")
	starred, err := client.GetStarredRepositories()
	if err != nil {
		log.Printf("   Could not get starred repos: %v\n", err)
	} else if len(starred.Repositories) == 0 {
		fmt.Println("   No starred repositories found")
	} else {
		fmt.Printf("   Found %d starred repositories:\n", len(starred.Repositories))
		for _, repo := range starred.Repositories {
			fmt.Printf("   - %s/%s\n", repo.Namespace, repo.Name)
		}
	}
	fmt.Println()

	// Step 5: List repositories (example with a namespace)
	// Note: Replace "your-namespace" with an actual namespace you have access to
	namespace := os.Getenv("QUAY_NAMESPACE")
	if namespace == "" {
		namespace = user.Username // Default to user's own namespace
	}

	fmt.Printf("3. Listing repositories in namespace '%s'...\n", namespace)
	repos, err := client.ListRepositories(namespace, false, false, 1, 10)
	if err != nil {
		log.Printf("   Could not list repositories: %v\n", err)
	} else if len(repos.Repositories) == 0 {
		fmt.Println("   No repositories found in this namespace")
	} else {
		fmt.Printf("   Found %d repositories:\n", len(repos.Repositories))
		for _, repo := range repos.Repositories {
			visibility := "private"
			if repo.IsPublic {
				visibility = "public"
			}
			fmt.Printf("   - %s/%s (%s)\n", repo.Namespace, repo.Name, visibility)
		}
	}
	fmt.Println()

	// Step 6: Get detailed repository information (if we have at least one repo)
	if repos != nil && len(repos.Repositories) > 0 {
		firstRepo := repos.Repositories[0]
		fmt.Printf("4. Getting details for repository '%s/%s'...\n", firstRepo.Namespace, firstRepo.Name)

		repoDetails, err := client.GetRepository(firstRepo.Namespace, firstRepo.Name)
		if err != nil {
			log.Printf("   Could not get repository details: %v\n", err)
		} else {
			fmt.Printf("   Name: %s\n", repoDetails.Name)
			fmt.Printf("   Namespace: %s\n", repoDetails.Namespace)
			fmt.Printf("   Description: %s\n", repoDetails.Description)
			fmt.Printf("   Public: %v\n", repoDetails.IsPublic)
			fmt.Printf("   Can Write: %v\n", repoDetails.CanWrite)
			fmt.Printf("   Can Admin: %v\n", repoDetails.CanAdmin)

			// Show tags if available
			if len(repoDetails.Tags.Tags) > 0 {
				fmt.Printf("   Tags (%d):\n", len(repoDetails.Tags.Tags))
				for i, tag := range repoDetails.Tags.Tags {
					if i >= 5 {
						fmt.Printf("   ... and %d more tags\n", len(repoDetails.Tags.Tags)-5)
						break
					}
					fmt.Printf("   - %s (modified: %s)\n", tag.Name, tag.LastModified)
				}
			} else {
				fmt.Println("   No tags found")
			}
		}
	}

	fmt.Println("\n=== Example Complete ===")
}

// prettyPrint outputs a struct as formatted JSON (utility function)
func prettyPrint(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal: %v", err)
		return
	}
	fmt.Println(string(data))
}
