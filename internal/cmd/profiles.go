package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/notte-cli/internal/api"
)

var (
	profilesCreateName string
)

var profilesCmd = &cobra.Command{
	Use:   "profiles",
	Short: "Manage browser profiles",
	Long:  "List and create browser profiles.",
}

var profilesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List profiles",
	RunE:  runProfilesList,
}

var profilesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new profile",
	RunE:  runProfilesCreate,
}

func init() {
	rootCmd.AddCommand(profilesCmd)
	profilesCmd.AddCommand(profilesListCmd)
	profilesCmd.AddCommand(profilesCreateCmd)

	profilesCreateCmd.Flags().StringVar(&profilesCreateName, "name", "", "Profile name")
}

func runProfilesList(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	params := &api.ProfileListParams{}
	resp, err := client.Client().ProfileListWithResponse(ctx, params)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	formatter := GetFormatter()

	var items []api.ProfileResponse
	if resp.JSON200 != nil {
		items = resp.JSON200.Items
	}
	if printed, err := PrintListOrEmpty(items, "No profiles found."); err != nil {
		return err
	} else if printed {
		return nil
	}

	return formatter.Print(items)
}

func runProfilesCreate(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	body := api.ProfileCreateJSONRequestBody{}
	if profilesCreateName != "" {
		body.Name = &profilesCreateName
	}

	params := &api.ProfileCreateParams{}
	resp, err := client.Client().ProfileCreateWithResponse(ctx, params, body)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	formatter := GetFormatter()
	return formatter.Print(resp.JSON200)
}
