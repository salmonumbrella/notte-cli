// internal/validate/validate.go
package validate

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// URL validates that a string is a valid HTTP/HTTPS URL
func URL(s string) error {
	if s == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	u, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme, got %q", u.Scheme)
	}

	if u.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	return nil
}

// JSON validates that a string is valid JSON
func JSON(s string) error {
	if s == "" {
		return fmt.Errorf("JSON cannot be empty")
	}

	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	return nil
}

// Browser validates browser type
func Browser(s string) error {
	valid := map[string]bool{
		"chromium": true,
		"firefox":  true,
		"webkit":   true,
	}

	if !valid[s] {
		return fmt.Errorf("invalid browser: expected chromium|firefox|webkit, got %q", s)
	}

	return nil
}

// PositiveInt validates that n > 0
func PositiveInt(n int) error {
	if n <= 0 {
		return fmt.Errorf("value must be positive, got %d", n)
	}
	return nil
}

// Duration validates that s is a valid Go duration string
func Duration(s string) error {
	if s == "" {
		return fmt.Errorf("duration cannot be empty")
	}

	_, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", s, err)
	}

	return nil
}

// OutputFormat validates output format flag
func OutputFormat(s string) error {
	valid := map[string]bool{
		"text": true,
		"json": true,
	}

	if !valid[s] {
		return fmt.Errorf("invalid output format: expected text|json, got %q", s)
	}

	return nil
}

// NonEmpty validates that a string is not empty
func NonEmpty(s, name string) error {
	if s == "" {
		return fmt.Errorf("%s cannot be empty", name)
	}
	return nil
}
