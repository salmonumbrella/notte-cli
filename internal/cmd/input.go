package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func readJSONInput(cmd *cobra.Command, value string, flagName string) ([]byte, error) {
	input := strings.TrimSpace(value)
	if input == "" {
		return readFromStdin(cmd, flagName)
	}

	if strings.HasPrefix(input, "@") {
		path := strings.TrimPrefix(input, "@")
		if path == "" {
			return nil, fmt.Errorf("invalid %s value: missing file path after @", flagName)
		}
		if path == "-" {
			return readFromStdin(cmd, flagName)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s file %q: %w", flagName, path, err)
		}
		if len(bytes.TrimSpace(data)) == 0 {
			return nil, fmt.Errorf("%s file %q is empty", flagName, path)
		}
		return data, nil
	}

	if input == "-" {
		return readFromStdin(cmd, flagName)
	}

	return []byte(input), nil
}

func readFromStdin(cmd *cobra.Command, flagName string) ([]byte, error) {
	in := cmd.InOrStdin()
	if !stdinHasData(in) {
		return nil, fmt.Errorf("%s is required (use --%s, --%s @file, or pipe JSON via stdin)", flagName, flagName, flagName)
	}

	data, err := io.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s from stdin: %w", flagName, err)
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, fmt.Errorf("%s input is empty", flagName)
	}
	return data, nil
}

func stdinHasData(r io.Reader) bool {
	file, ok := r.(*os.File)
	if !ok {
		return true
	}
	info, err := file.Stat()
	if err != nil {
		return true
	}
	return (info.Mode() & os.ModeCharDevice) == 0
}
