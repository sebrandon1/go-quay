/*
Package cmd provides the command-line interface for go-quay.

This file contains the commands for the Prototype API:
  - go-quay get prototype list     - List all permission prototypes
  - go-quay get prototype info     - Get a specific prototype
  - go-quay get prototype create   - Create a new prototype
  - go-quay get prototype update   - Update a prototype
  - go-quay get prototype delete   - Delete a prototype
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sebrandon1/go-quay/lib"
	"github.com/spf13/cobra"
)

var (
	prototypeOrg          string
	prototypeUUID         string
	prototypeDelegateName string
	prototypeDelegateKind string
	prototypeRole         string
	confirmProtoDelete    bool
)

// prototypeCmd represents the prototype command
var prototypeCmd = &cobra.Command{
	Use:   "prototype",
	Short: "Manage organization permission prototypes",
	Long: `Manage default permission prototypes for an organization.

Prototypes define default permissions that are automatically applied to new
repositories created within an organization. They allow setting up permission
templates for users, teams, or robot accounts.`,
}

// prototypeListCmd lists all prototypes
var prototypeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all permission prototypes",
	Long:  `List all permission prototypes for an organization.`,
	Run: func(_ *cobra.Command, _ []string) {
		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		prototypes, err := client.GetPrototypes(prototypeOrg)
		if err != nil {
			fmt.Println("Error getting prototypes:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(prototypes, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

// prototypeInfoCmd gets a specific prototype
var prototypeInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get a specific prototype",
	Long:  `Get detailed information about a specific permission prototype.`,
	Run: func(_ *cobra.Command, _ []string) {
		if prototypeUUID == "" {
			fmt.Println("Error: --uuid is required")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		prototype, err := client.GetPrototype(prototypeOrg, prototypeUUID)
		if err != nil {
			fmt.Println("Error getting prototype:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(prototype, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

// prototypeCreateCmd creates a new prototype
var prototypeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new prototype",
	Long: `Create a new permission prototype for an organization.

Delegate kinds:
  - user: A specific user account
  - team: A team within the organization  
  - robot: A robot account

Roles:
  - read: Pull images
  - write: Pull and push images
  - admin: Full administrative access`,
	Run: func(_ *cobra.Command, _ []string) {
		if prototypeDelegateName == "" {
			fmt.Println("Error: --delegate-name is required")
			os.Exit(1)
		}
		if prototypeDelegateKind == "" {
			fmt.Println("Error: --delegate-kind is required")
			os.Exit(1)
		}
		if prototypeRole == "" {
			fmt.Println("Error: --role is required")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		createReq := &lib.CreatePrototypeRequest{
			Delegate: lib.PrototypeDelegateRequest{
				Name: prototypeDelegateName,
				Kind: prototypeDelegateKind,
			},
			Role: prototypeRole,
		}

		prototype, err := client.CreatePrototype(prototypeOrg, createReq)
		if err != nil {
			fmt.Println("Error creating prototype:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(prototype, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

// prototypeUpdateCmd updates a prototype
var prototypeUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a prototype",
	Long:  `Update an existing permission prototype's role.`,
	Run: func(_ *cobra.Command, _ []string) {
		if prototypeUUID == "" {
			fmt.Println("Error: --uuid is required")
			os.Exit(1)
		}
		if prototypeRole == "" {
			fmt.Println("Error: --role is required")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		updateReq := &lib.UpdatePrototypeRequest{
			Role: prototypeRole,
		}

		prototype, err := client.UpdatePrototype(prototypeOrg, prototypeUUID, updateReq)
		if err != nil {
			fmt.Println("Error updating prototype:", err)
			os.Exit(1)
		}

		output, err := json.MarshalIndent(prototype, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	},
}

// prototypeDeleteCmd deletes a prototype
var prototypeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a prototype",
	Long:  `Delete a permission prototype from an organization.`,
	Run: func(_ *cobra.Command, _ []string) {
		if prototypeUUID == "" {
			fmt.Println("Error: --uuid is required")
			os.Exit(1)
		}
		if !confirmProtoDelete {
			fmt.Println("Error: --confirm is required to delete a prototype")
			os.Exit(1)
		}

		client, err := lib.NewClient(token)
		if err != nil {
			fmt.Println("Error creating client:", err)
			os.Exit(1)
		}

		err = client.DeletePrototype(prototypeOrg, prototypeUUID)
		if err != nil {
			fmt.Println("Error deleting prototype:", err)
			os.Exit(1)
		}

		fmt.Printf("Prototype %s deleted successfully\n", prototypeUUID)
	},
}

func setupPrototypeFlags() {
	// Common flags
	for _, cmd := range []*cobra.Command{prototypeListCmd, prototypeInfoCmd, prototypeCreateCmd, prototypeUpdateCmd, prototypeDeleteCmd} {
		cmd.Flags().StringVarP(&prototypeOrg, "organization", "o", "", "Organization name")
		cmd.Flags().StringVarP(&token, "token", "t", "", "Quay.io API token")
	}

	// UUID flags
	for _, cmd := range []*cobra.Command{prototypeInfoCmd, prototypeUpdateCmd, prototypeDeleteCmd} {
		cmd.Flags().StringVar(&prototypeUUID, "uuid", "", "Prototype UUID")
	}

	// Create flags
	prototypeCreateCmd.Flags().StringVar(&prototypeDelegateName, "delegate-name", "", "Name of the delegate (user/team/robot)")
	prototypeCreateCmd.Flags().StringVar(&prototypeDelegateKind, "delegate-kind", "", "Kind of delegate (user, team, robot)")
	prototypeCreateCmd.Flags().StringVar(&prototypeRole, "role", "", "Permission role (read, write, admin)")

	// Update flags
	prototypeUpdateCmd.Flags().StringVar(&prototypeRole, "role", "", "New permission role (read, write, admin)")

	// Delete flags
	prototypeDeleteCmd.Flags().BoolVar(&confirmProtoDelete, "confirm", false, "Confirm deletion")
}

func init() {
	prototypeCmd.AddCommand(prototypeListCmd)
	prototypeCmd.AddCommand(prototypeInfoCmd)
	prototypeCmd.AddCommand(prototypeCreateCmd)
	prototypeCmd.AddCommand(prototypeUpdateCmd)
	prototypeCmd.AddCommand(prototypeDeleteCmd)

	setupPrototypeFlags()
}
