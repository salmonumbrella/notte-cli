package api

import (
	"encoding/json"
	"time"
)

// FlexibleTime is a time.Time that can parse multiple timestamp formats.
// It handles both RFC3339 (with timezone) and timestamps without timezone info.
type FlexibleTime struct {
	time.Time
}

// Supported time formats in order of preference
var timeFormats = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05.999999",
	"2006-01-02T15:04:05",
	"2006-01-02",
}

// UnmarshalJSON implements json.Unmarshaler for FlexibleTime
func (ft *FlexibleTime) UnmarshalJSON(data []byte) error {
	// Handle null
	if string(data) == "null" {
		return nil
	}

	// Remove quotes
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// Empty string
	if s == "" {
		return nil
	}

	// Try each format
	var parseErr error
	for _, format := range timeFormats {
		t, err := time.Parse(format, s)
		if err == nil {
			ft.Time = t
			return nil
		}
		parseErr = err
	}

	return parseErr
}

// MarshalJSON implements json.Marshaler for FlexibleTime
func (ft FlexibleTime) MarshalJSON() ([]byte, error) {
	if ft.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(ft.Time.Format(time.RFC3339Nano))
}
