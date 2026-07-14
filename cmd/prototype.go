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
	"fmt"

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
	Use:   subcmdList,
	Short: "List all permission prototypes",
	Long:  `List all permission prototypes for an organization.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		prototypes, err := client.GetPrototypes(prototypeOrg)
		if err != nil {
			return fmt.Errorf("getting prototypes: %w", err)
		}

		return printJSON(prototypes)
	},
}

// prototypeInfoCmd gets a specific prototype
var prototypeInfoCmd = &cobra.Command{
	Use:   subcmdInfo,
	Short: "Get a specific prototype",
	Long:  `Get detailed information about a specific permission prototype.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if prototypeUUID == "" {
			return fmt.Errorf("--uuid is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		prototype, err := client.GetPrototype(prototypeOrg, prototypeUUID)
		if err != nil {
			return fmt.Errorf("getting prototype: %w", err)
		}

		return printJSON(prototype)
	},
}

// prototypeCreateCmd creates a new prototype
var prototypeCreateCmd = &cobra.Command{
	Use:   subcmdCreate,
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
	RunE: func(_ *cobra.Command, _ []string) error {
		if prototypeDelegateName == "" {
			return fmt.Errorf("--delegate-name is required")
		}
		if prototypeDelegateKind == "" {
			return fmt.Errorf("--delegate-kind is required")
		}
		if prototypeRole == "" {
			return fmt.Errorf("--role is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
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
			return fmt.Errorf("creating prototype: %w", err)
		}

		return printJSON(prototype)
	},
}

// prototypeUpdateCmd updates a prototype
var prototypeUpdateCmd = &cobra.Command{
	Use:   subcmdUpdate,
	Short: "Update a prototype",
	Long:  `Update an existing permission prototype's role.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if prototypeUUID == "" {
			return fmt.Errorf("--uuid is required")
		}
		if prototypeRole == "" {
			return fmt.Errorf("--role is required")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		updateReq := &lib.UpdatePrototypeRequest{
			Role: prototypeRole,
		}

		prototype, err := client.UpdatePrototype(prototypeOrg, prototypeUUID, updateReq)
		if err != nil {
			return fmt.Errorf("updating prototype: %w", err)
		}

		return printJSON(prototype)
	},
}

// prototypeDeleteCmd deletes a prototype
var prototypeDeleteCmd = &cobra.Command{
	Use:   subcmdDelete,
	Short: "Delete a prototype",
	Long:  `Delete a permission prototype from an organization.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		if prototypeUUID == "" {
			return fmt.Errorf("--uuid is required")
		}
		if !confirmProtoDelete {
			return fmt.Errorf("--confirm is required to delete a prototype")
		}

		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}

		err = client.DeletePrototype(prototypeOrg, prototypeUUID)
		if err != nil {
			return fmt.Errorf("deleting prototype: %w", err)
		}

		fmt.Printf("Prototype %s deleted successfully\n", prototypeUUID)
		return nil
	},
}

func setupPrototypeFlags() {
	// Common flags
	for _, cmd := range []*cobra.Command{prototypeListCmd, prototypeInfoCmd, prototypeCreateCmd, prototypeUpdateCmd, prototypeDeleteCmd} {
		cmd.Flags().StringVarP(&prototypeOrg, "organization", "o", "", "Organization name")
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
