package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Organization Teams
var orgTeamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Get organization teams",
	Long:  `Get list of all teams in an organization.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		teams, err := client.GetTeams(cmd.Context(), orgName)
		if err != nil {
			return fmt.Errorf("getting organization teams: %w", err)
		}
		return printJSON(teams)
	},
}

// Team Info
var teamInfoCmd = &cobra.Command{
	Use:   "team",
	Short: "Get team information",
	Long:  `Get detailed information about a specific team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		team, err := client.GetTeam(cmd.Context(), orgName, teamName)
		if err != nil {
			return fmt.Errorf("getting team: %w", err)
		}
		return printJSON(team)
	},
}

// Team Members
var teamMembersCmd = &cobra.Command{
	Use:   "team-members",
	Short: "Get team members",
	Long:  `Get list of all members in a specific team.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		members, err := client.GetTeamMembers(cmd.Context(), orgName, teamName)
		if err != nil {
			return fmt.Errorf("getting team members: %w", err)
		}
		return printJSON(members)
	},
}

// Invite Team Member
var inviteMemberCmd = &cobra.Command{
	Use:   "invite-member",
	Short: "Invite a member to a team",
	Long:  `Invite a member to a team by email.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.InviteTeamMember(cmd.Context(), orgName, teamName, email)
		if err != nil {
			return fmt.Errorf("inviting team member: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Team member invited successfully")
		return nil
	},
}

// Cancel Team Invite
var cancelInviteCmd = &cobra.Command{
	Use:   "cancel-invite",
	Short: "Cancel a team member invitation",
	Long:  `Cancel a pending team member invitation. Requires --confirm flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !confirm {
			return fmt.Errorf("must pass --confirm to cancel an invite")
		}
		client, err := getClient()
		if err != nil {
			return fmt.Errorf("creating client: %w", err)
		}
		err = client.DeleteTeamInvite(cmd.Context(), orgName, teamName, email)
		if err != nil {
			return fmt.Errorf("canceling team invite: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Team invite canceled successfully")
		return nil
	},
}

func initOrgTeamFlags() {
	teamInfoCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	_ = teamInfoCmd.MarkFlagRequired("team")

	teamMembersCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	_ = teamMembersCmd.MarkFlagRequired("team")
}

func initOrgInviteFlags() {
	inviteMemberCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	_ = inviteMemberCmd.MarkFlagRequired("team")
	inviteMemberCmd.Flags().StringVar(&email, "email", "", "Email address")
	_ = inviteMemberCmd.MarkFlagRequired("email")

	cancelInviteCmd.Flags().StringVarP(&teamName, "team", "T", "", "Team name")
	_ = cancelInviteCmd.MarkFlagRequired("team")
	cancelInviteCmd.Flags().StringVar(&email, "email", "", "Email address")
	_ = cancelInviteCmd.MarkFlagRequired("email")
	cancelInviteCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm cancellation")
}

func init() {
	organizationCmd.AddCommand(orgTeamsCmd)
	organizationCmd.AddCommand(teamInfoCmd)
	organizationCmd.AddCommand(teamMembersCmd)
	organizationCmd.AddCommand(inviteMemberCmd)
	organizationCmd.AddCommand(cancelInviteCmd)

	initOrgTeamFlags()
	initOrgInviteFlags()
}
