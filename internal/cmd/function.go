package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/notte-cli/internal/api"
)

var (
	functionID             string
	functionUpdateFile     string
	functionRunID          string
	functionMetadataJSON   string
	functionCronExpression string
)

var functionCmd = &cobra.Command{
	Use:   "function",
	Short: "Operate on a specific function",
	Long:  "Subcommands for interacting with a specific function. Use --id flag to specify the function ID.",
}

var functionShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show function details",
	Args:  cobra.NoArgs,
	RunE:  runFunctionShow,
}

var functionUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update function",
	Args:  cobra.NoArgs,
	RunE:  runFunctionUpdate,
}

var functionDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete function",
	Args:  cobra.NoArgs,
	RunE:  runFunctionDelete,
}

var functionRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the function",
	Args:  cobra.NoArgs,
	RunE:  runFunctionRun,
}

var functionRunsCmd = &cobra.Command{
	Use:   "runs",
	Short: "List function runs",
	Args:  cobra.NoArgs,
	RunE:  runFunctionRuns,
}

var functionForkCmd = &cobra.Command{
	Use:   "fork",
	Short: "Fork/duplicate the function",
	Args:  cobra.NoArgs,
	RunE:  runFunctionFork,
}

var functionRunStopCmd = &cobra.Command{
	Use:   "run-stop",
	Short: "Stop a function run",
	Args:  cobra.NoArgs,
	RunE:  runFunctionRunStop,
}

var functionRunMetadataCmd = &cobra.Command{
	Use:   "run-metadata",
	Short: "Get function run metadata",
	Args:  cobra.NoArgs,
	RunE:  runFunctionRunMetadata,
}

var functionRunMetadataUpdateCmd = &cobra.Command{
	Use:   "run-metadata-update",
	Short: "Update function run metadata",
	Args:  cobra.NoArgs,
	Example: `  # Direct JSON
  notte function run-metadata-update --id <function-id> --run-id <run-id> --data '{"key": "value"}'

  # From file
  notte function run-metadata-update --id <function-id> --run-id <run-id> --data @metadata.json

  # From stdin
  echo '{"key": "value"}' | notte function run-metadata-update --id <function-id> --run-id <run-id>`,
	RunE: runFunctionRunMetadataUpdate,
}

var functionScheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Set a cron schedule for the function",
	Args:  cobra.NoArgs,
	RunE:  runFunctionSchedule,
}

var functionUnscheduleCmd = &cobra.Command{
	Use:   "unschedule",
	Short: "Remove the schedule from the function",
	Args:  cobra.NoArgs,
	RunE:  runFunctionUnschedule,
}

func init() {
	rootCmd.AddCommand(functionCmd)

	functionCmd.AddCommand(functionShowCmd)
	functionCmd.AddCommand(functionUpdateCmd)
	functionCmd.AddCommand(functionDeleteCmd)
	functionCmd.AddCommand(functionRunCmd)
	functionCmd.AddCommand(functionRunsCmd)
	functionCmd.AddCommand(functionForkCmd)
	functionCmd.AddCommand(functionRunStopCmd)
	functionCmd.AddCommand(functionRunMetadataCmd)
	functionCmd.AddCommand(functionRunMetadataUpdateCmd)
	functionCmd.AddCommand(functionScheduleCmd)
	functionCmd.AddCommand(functionUnscheduleCmd)

	// Persistent flag for function ID (required for all subcommands)
	functionCmd.PersistentFlags().StringVar(&functionID, "id", "", "Function ID (required)")
	_ = functionCmd.MarkPersistentFlagRequired("id")

	functionUpdateCmd.Flags().StringVar(&functionUpdateFile, "file", "", "Path to updated function file (required)")
	_ = functionUpdateCmd.MarkFlagRequired("file")

	// Flags for run control commands
	functionRunStopCmd.Flags().StringVar(&functionRunID, "run-id", "", "Run ID (required)")
	_ = functionRunStopCmd.MarkFlagRequired("run-id")

	functionRunMetadataCmd.Flags().StringVar(&functionRunID, "run-id", "", "Run ID (required)")
	_ = functionRunMetadataCmd.MarkFlagRequired("run-id")

	functionRunMetadataUpdateCmd.Flags().StringVar(&functionRunID, "run-id", "", "Run ID (required)")
	_ = functionRunMetadataUpdateCmd.MarkFlagRequired("run-id")
	functionRunMetadataUpdateCmd.Flags().StringVar(&functionMetadataJSON, "data", "", "JSON metadata, @file, or '-' for stdin")

	functionScheduleCmd.Flags().StringVar(&functionCronExpression, "cron", "", "Cron expression (required)")
	_ = functionScheduleCmd.MarkFlagRequired("cron")
}

func runFunctionShow(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	params := &api.FunctionDownloadUrlParams{}
	resp, err := client.Client().FunctionDownloadUrlWithResponse(ctx, functionID, params)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	return GetFormatter().Print(resp.JSON200)
}

func runFunctionUpdate(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	// Open the function file
	file, err := os.Open(functionUpdateFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add file field
	part, err := writer.CreateFormFile("file", filepath.Base(functionUpdateFile))
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy file data: %w", err)
	}

	_ = writer.Close()

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	params := &api.FunctionUpdateParams{}
	resp, err := client.Client().FunctionUpdateWithBodyWithResponse(
		ctx,
		functionID,
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

	return GetFormatter().Print(resp.JSON200)
}

func runFunctionDelete(cmd *cobra.Command, args []string) error {
	confirmed, err := ConfirmAction("function", functionID)
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

	params := &api.FunctionDeleteParams{}
	resp, err := client.Client().FunctionDeleteWithResponse(ctx, functionID, params)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	return PrintResult(fmt.Sprintf("Function %s deleted.", functionID), map[string]any{
		"id":     functionID,
		"status": "deleted",
	})
}

func runFunctionRun(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	params := &api.FunctionRunStartParams{}
	resp, err := client.Client().FunctionRunStartWithResponse(ctx, functionID, params)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	return GetFormatter().Print(resp.JSON200)
}

func runFunctionRuns(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	params := &api.ListFunctionRunsByFunctionIdParams{}
	resp, err := client.Client().ListFunctionRunsByFunctionIdWithResponse(ctx, functionID, params)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	var items []api.GetFunctionRunResponse
	if resp.JSON200 != nil {
		items = resp.JSON200.Items
	}
	if printed, err := PrintListOrEmpty(items, "No function runs found."); err != nil {
		return err
	} else if printed {
		return nil
	}

	return GetFormatter().Print(items)
}

func runFunctionFork(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	params := &api.FunctionForkParams{}
	resp, err := client.Client().FunctionForkWithResponse(ctx, functionID, params)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	return GetFormatter().Print(resp.JSON200)
}

func runFunctionRunStop(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	params := &api.FunctionRunStopParams{}
	resp, err := client.Client().FunctionRunStopWithResponse(ctx, functionID, functionRunID, params)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	return GetFormatter().Print(resp.JSON200)
}

func runFunctionRunMetadata(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	params := &api.FunctionRunGetMetadataParams{}
	resp, err := client.Client().FunctionRunGetMetadataWithResponse(ctx, functionID, functionRunID, params)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	return GetFormatter().Print(resp.JSON200)
}

func runFunctionRunMetadataUpdate(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	metadataPayload, err := readJSONInput(cmd, functionMetadataJSON, "data")
	if err != nil {
		return err
	}

	// Parse the JSON metadata
	var metadata api.FunctionRunUpdateMetadataJSONRequestBody
	if err := json.Unmarshal(metadataPayload, &metadata); err != nil {
		return fmt.Errorf("failed to parse JSON metadata: %w", err)
	}

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	params := &api.FunctionRunUpdateMetadataParams{}
	resp, err := client.Client().FunctionRunUpdateMetadataWithResponse(ctx, functionID, functionRunID, params, metadata)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	return GetFormatter().Print(resp.JSON200)
}

func runFunctionSchedule(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	body := api.FunctionScheduleSetJSONRequestBody{
		Cron: functionCronExpression,
	}

	params := &api.FunctionScheduleSetParams{}
	resp, err := client.Client().FunctionScheduleSetWithResponse(ctx, functionID, params, body)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	return PrintResult(fmt.Sprintf("Function %s scheduled with cron expression: %s", functionID, functionCronExpression), map[string]any{
		"id":   functionID,
		"cron": functionCronExpression,
	})
}

func runFunctionUnschedule(cmd *cobra.Command, args []string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	ctx, cancel := GetContextWithTimeout(cmd.Context())
	defer cancel()

	params := &api.FunctionScheduleDeleteParams{}
	resp, err := client.Client().FunctionScheduleDeleteWithResponse(ctx, functionID, params)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}

	if err := HandleAPIResponse(resp.HTTPResponse); err != nil {
		return err
	}

	return PrintResult(fmt.Sprintf("Function %s schedule removed.", functionID), map[string]any{
		"id":     functionID,
		"status": "unscheduled",
	})
}
