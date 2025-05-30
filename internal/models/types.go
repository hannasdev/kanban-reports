package models

import (
	"fmt"
	"strings"
	"time"
)

// FilterField defines which date field to use for filtering
type FilterField string

const (
		// FilterFieldCompletedAt uses the completed_at date for filtering
		FilterFieldCompletedAt FilterField = "completed_at"
		// FilterFieldCreatedAt uses the created_at date for filtering
		FilterFieldCreatedAt FilterField = "created_at"
		// FilterFieldStartedAt uses the started_at date for filtering
		FilterFieldStartedAt FilterField = "started_at"
)

// IsValid checks if a FilterField is valid
func (ff FilterField) IsValid() bool {
		switch ff {
		case FilterFieldCompletedAt, FilterFieldCreatedAt, FilterFieldStartedAt:
				return true
		}
		return false
}

// ParseFilterField converts a string to a FilterField with validation
func ParseFilterField(s string) (FilterField, error) {
		ff := FilterField(s)
		if !ff.IsValid() {
				return "", fmt.Errorf("invalid filter field: %s (must be one of: completed_at, created_at, started_at)", s)
		}
		return ff, nil
}

// GetItemDate returns the appropriate date from the KanbanItem based on this filter field
func (ff FilterField) GetItemDate(item KanbanItem) (time.Time, bool) {
	switch ff {
	case FilterFieldCompletedAt:
			// For completed_at, only return a date if the item is actually completed
			if !item.IsCompleted || item.CompletedAt.IsZero() {
					return time.Time{}, false
			}
			return item.CompletedAt, true
	case FilterFieldCreatedAt:
			return item.CreatedAt, !item.CreatedAt.IsZero()
	case FilterFieldStartedAt:
			return item.StartedAt, !item.StartedAt.IsZero()
	default:
			// Default to completed_at
			if !item.IsCompleted || item.CompletedAt.IsZero() {
					return time.Time{}, false
			}
			return item.CompletedAt, true
	}
}

// DelimiterType defines the type of delimiter used in CSV files
type DelimiterType struct {
	Value     rune
	Name      string
	AutoDetect bool
}

// Common delimiter constants
var (
	DelimiterComma     = DelimiterType{Value: ',', Name: "comma", AutoDetect: false}
	DelimiterTab       = DelimiterType{Value: '\t', Name: "tab", AutoDetect: false}
	DelimiterSemicolon = DelimiterType{Value: ';', Name: "semicolon", AutoDetect: false}
	DelimiterAuto      = DelimiterType{Value: ',', Name: "auto", AutoDetect: true} // Default to comma but will auto-detect
)

// ParseDelimiter converts a string to a DelimiterType with validation
func ParseDelimiter(s string) (DelimiterType, error) {
	switch s {
	case "comma":
			return DelimiterComma, nil
	case "tab":
			return DelimiterTab, nil
	case "semicolon":
			return DelimiterSemicolon, nil
	case "auto":
			return DelimiterAuto, nil
	default:
			return DelimiterAuto, fmt.Errorf("invalid delimiter type: %s (must be one of: comma, tab, semicolon, auto)", s)
	}
}

// DetectDelimiterType automatically detects the delimiter from sample content
func DetectDelimiterType(content string) DelimiterType {
	// Take a sample of up to first 5 lines to improve detection
	lines := strings.SplitN(content, "\n", 6)
	
	// Ensure we don't try to access beyond the slice length
	if len(lines) > 5 {
			lines = lines[:5]
	}
	
	delimiterCounts := map[rune]int{
			',': 0, 
			'\t': 0, 
			';': 0,
	}
	
	// Count delimiters in each line
	for _, line := range lines {
			delimiterCounts[','] += strings.Count(line, ",")
			delimiterCounts['\t'] += strings.Count(line, "\t")
			delimiterCounts[';'] += strings.Count(line, ";")
	}
	
	// Find delimiter with highest count
	maxCount := -1
	bestDelimiter := ','
	
	for delimiter, count := range delimiterCounts {
			if count > maxCount {
					maxCount = count
					bestDelimiter = delimiter
			}
	}
	
	// Map the rune back to our delimiter type
	switch bestDelimiter {
	case ',':
			return DelimiterComma
	case '\t':
			return DelimiterTab
	case ';':
			return DelimiterSemicolon
	default:
			return DelimiterComma // Fallback to comma
	}
}

