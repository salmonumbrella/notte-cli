package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/notte-cli/internal/api"
)

var (
	profileID string
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Operate on a specific profile",
	Long:  "Subcommands for interacting with a specific profile. Use --id flag to specify the profile ID.",
}

var profileShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show profile details",
	Args:  cobra.NoArgs,
	RunE:  runProfileShow,
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete profile",
	Args:  cobra.NoArgs,
	RunE:  runProfileDelete,
}

func init() {
	rootCmd.AddCommand(profileCmd)

	profileCmd.AddCommand(profileShowCmd)
	profileCmd.AddCommand(profileDeleteCmd)

	// Persistent flag for profile ID (required for all subcommands)
	profileCmd.PersistentFlags().StringVar(&profileID, "id", "", "Profile ID (required)")
	_ = profileCmd.MarkPersistentFlagRequired("id")
}

func runProfileShow(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	params := &api.ProfileGetParams{}
	resp, err := client.Client().ProfileGetWithResponse(ctx, profileID, params)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	return GetFormatter().Print(resp.JSON200)
}

func runProfileDelete(cmd *cobra.Command, args []string) error {
	confirmed, err := ConfirmAction("profile", profileID)
	if err != nil {
		return err
	}
	if !confirmed {
		return PrintResult("Cancelled.", map[string]any{"cancelled": true})
	}

	client, err := GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	params := &api.ProfileDeleteParams{}
	resp, err := client.Client().ProfileDeleteWithResponse(ctx, profileID, params)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	return PrintResult(fmt.Sprintf("Profile %s deleted.", profileID), map[string]any{
		"id":     profileID,
		"status": "deleted",
	})
}
