package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/notte-cli/internal/api"
)

var (
	functionsCreateFile        string
	functionsCreateName        string
	functionsCreateDescription string
	functionsCreateShared      bool
)

var functionsCmd = &cobra.Command{
	Use:   "functions",
	Short: "Manage functions",
	Long:  "List and create functions.",
}

var functionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List functions",
	RunE:  runFunctionsList,
}

var functionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new function",
	RunE:  runFunctionsCreate,
}

func init() {
	rootCmd.AddCommand(functionsCmd)
	functionsCmd.AddCommand(functionsListCmd)
	functionsCmd.AddCommand(functionsCreateCmd)

	functionsCreateCmd.Flags().StringVar(&functionsCreateFile, "file", "", "Path to function file (required)")
	_ = functionsCreateCmd.MarkFlagRequired("file")
	functionsCreateCmd.Flags().StringVar(&functionsCreateName, "name", "", "Function name")
	functionsCreateCmd.Flags().StringVar(&functionsCreateDescription, "description", "", "Function description")
	functionsCreateCmd.Flags().BoolVar(&functionsCreateShared, "shared", false, "Make function public")
}

func runFunctionsList(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	params := &api.ListFunctionsParams{}
	resp, err := client.Client().ListFunctionsWithResponse(ctx, params)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	formatter := GetFormatter()

	var items []api.GetFunctionResponse
	if resp.JSON200 != nil {
		items = resp.JSON200.Items
	}
	if printed, err := PrintListOrEmpty(items, "No functions found."); err != nil {
		return err
	} else if printed {
		return nil
	}

	return formatter.Print(items)
}

func runFunctionsCreate(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	// Open the function file
	file, err := os.Open(functionsCreateFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add file field
	part, err := writer.CreateFormFile("file", filepath.Base(functionsCreateFile))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file data: %w", err)
	}

	// Add optional fields
	if functionsCreateName != "" {
		if err := writer.WriteField("name", functionsCreateName); err != nil {
			return fmt.Errorf("failed to write name field: %w", err)
		}
	}
	if functionsCreateDescription != "" {
		if err := writer.WriteField("description", functionsCreateDescription); err != nil {
			return fmt.Errorf("failed to write description field: %w", err)
		}
	}
	if cmd.Flags().Changed("shared") {
		if err := writer.WriteField("shared", fmt.Sprintf("%t", functionsCreateShared)); err != nil {
			return fmt.Errorf("failed to write shared field: %w", err)
		}
	}

	_ = writer.Close()

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	params := &api.FunctionCreateParams{}
	resp, err := client.Client().FunctionCreateWithBodyWithResponse(
		ctx,
		params,
		writer.FormDataContentType(),
		&buf,
	)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	formatter := GetFormatter()
	return formatter.Print(resp.JSON200)
}
