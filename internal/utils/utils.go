package utils

import (
	"encoding/json"
	"github.com/charmbracelet/x/ansi"
)

func Wordwrap(content string, width int) string {
	// ADHD: these breakpoints are silly.
	var breakpoints string = " ,-"
	return ansi.Wordwrap(content, width, breakpoints)
}

// Clamp ensures that the value is within the specified min and max range.
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func UnmarshalJSONToMap(data string) (map[string]any, error) {
	var result map[string]any
	err := json.Unmarshal([]byte(data), &result)
	return result, err
}
