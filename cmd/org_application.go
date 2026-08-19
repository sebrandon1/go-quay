package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Organization Applications
var orgApplicationsCmd = &cobra.Command{
	Use:   "applications",
	Short: "Get organization applications",
	Long:  `Get list of OAuth applications for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		applications, err := client.GetApplications(cmd.Context(), orgName)
		if err != nil {
			return fmt.Errorf("getting organization applications: %w", err)
		}
		return printJSON(applications)
	},
}

// Get Application
var applicationCmd = &cobra.Command{
	Use:   "application",
	Short: "Get application information",
	Long:  `Get detailed information about a specific OAuth application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		app, err := client.GetApplication(cmd.Context(), orgName, clientID)
		if err != nil {
			return fmt.Errorf("getting application: %w", err)
		}
		return printJSON(app)
	},
}

// Create Application
var createApplicationCmd = &cobra.Command{
	Use:   "create-application",
	Short: "Create an OAuth application",
	Long:  `Create a new OAuth application for an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		app, err := client.CreateApplication(cmd.Context(), orgName, appName, description, applicationURI, redirectURI)
		if err != nil {
			return fmt.Errorf("creating application: %w", err)
		}
		return printJSON(app)
	},
}

// Update Application
var updateApplicationCmd = &cobra.Command{
	Use:   "update-application",
	Short: "Update an OAuth application",
	Long:  `Update an existing OAuth application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		app, err := client.UpdateApplication(cmd.Context(), orgName, clientID, appName, description, applicationURI, redirectURI)
		if err != nil {
			return fmt.Errorf("updating application: %w", err)
		}
		return printJSON(app)
	},
}

// Delete Application
var deleteApplicationCmd = &cobra.Command{
	Use:   "delete-application",
	Short: "Delete an OAuth application",
	Long:  `Delete an OAuth application. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to delete an application")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.DeleteApplication(cmd.Context(), orgName, clientID)
		if err != nil {
			return fmt.Errorf("deleting application: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Application deleted successfully")
		return nil
	},
}

// Reset Application Client Secret
var resetApplicationSecretCmd = &cobra.Command{
	Use:   "reset-application-secret",
	Short: "Reset application client secret",
	Long:  `Reset the client secret for an OAuth application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		app, err := client.ResetApplicationClientSecret(cmd.Context(), orgName, clientID)
		if err != nil {
			return fmt.Errorf("resetting application secret: %w", err)
		}
		return printJSON(app)
	},
}

func initOrgApplicationFlags() {
	applicationCmd.Flags().StringVar(&clientID, "client-id", "", "Application client ID")
	_ = applicationCmd.MarkFlagRequired("client-id")

	createApplicationCmd.Flags().StringVar(&appName, "name", "", "Application name")
	_ = createApplicationCmd.MarkFlagRequired("name")
	createApplicationCmd.Flags().StringVar(&description, "description", "", "Application description")
	createApplicationCmd.Flags().StringVar(&applicationURI, "application-uri", "", "Application URI")
	createApplicationCmd.Flags().StringVar(&redirectURI, "redirect-uri", "", "Redirect URI")

	updateApplicationCmd.Flags().StringVar(&clientID, "client-id", "", "Application client ID")
	_ = updateApplicationCmd.MarkFlagRequired("client-id")
	updateApplicationCmd.Flags().StringVar(&appName, "name", "", "Application name")
	_ = updateApplicationCmd.MarkFlagRequired("name")
	updateApplicationCmd.Flags().StringVar(&description, "description", "", "Application description")
	updateApplicationCmd.Flags().StringVar(&applicationURI, "application-uri", "", "Application URI")
	updateApplicationCmd.Flags().StringVar(&redirectURI, "redirect-uri", "", "Redirect URI")

	deleteApplicationCmd.Flags().StringVar(&clientID, "client-id", "", "Application client ID")
	_ = deleteApplicationCmd.MarkFlagRequired("client-id")
	deleteApplicationCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")

	resetApplicationSecretCmd.Flags().StringVar(&clientID, "client-id", "", "Application client ID")
	_ = resetApplicationSecretCmd.MarkFlagRequired("client-id")
}

func init() {
	organizationCmd.AddCommand(orgApplicationsCmd)
	organizationCmd.AddCommand(applicationCmd)
	organizationCmd.AddCommand(createApplicationCmd)
	organizationCmd.AddCommand(updateApplicationCmd)
	organizationCmd.AddCommand(deleteApplicationCmd)
	organizationCmd.AddCommand(resetApplicationSecretCmd)

	initOrgApplicationFlags()
}
