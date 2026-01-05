package cmd

import (
	"fmt"
	"os"
	"reflect"
)

func IsJSONOutput() bool {
	return outputFormat == "json"
}

func PrintInfo(message string) {
	if IsJSONOutput() {
		_, _ = fmt.Fprintln(os.Stderr, message)
		return
	}
	_, _ = fmt.Fprintln(os.Stdout, message)
}

func PrintResult(message string, data map[string]any) error {
	if IsJSONOutput() {
		if data == nil {
			data = map[string]any{}
		}
		if _, ok := data["message"]; !ok && message != "" {
			data["message"] = message
		}
		return GetFormatter().Print(data)
	}

	if message == "" {
		return nil
	}
	_, err := fmt.Fprintln(os.Stdout, message)
	return err
}

func PrintListOrEmpty(items any, emptyMsg string) (bool, error) {
	if items == nil {
		if IsJSONOutput() {
			return true, GetFormatter().Print([]any{})
		}
		if emptyMsg != "" {
			_, _ = fmt.Fprintln(os.Stdout, emptyMsg)
		}
		return true, nil
	}

	v := reflect.ValueOf(items)
	if v.Kind() != reflect.Slice {
		return false, nil
	}

	if v.Len() == 0 {
		if IsJSONOutput() {
			empty := reflect.MakeSlice(v.Type(), 0, 0).Interface()
			return true, GetFormatter().Print(empty)
		}
		if emptyMsg != "" {
			_, _ = fmt.Fprintln(os.Stdout, emptyMsg)
		}
		return true, nil
	}

	return false, nil
}
