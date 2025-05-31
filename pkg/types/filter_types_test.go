package types

import (
	"testing"
)

func TestAdHocFilterType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		filter   AdHocFilterType
		expected bool
	}{
		{"Valid include", AdHocFilterInclude, true},
		{"Valid exclude", AdHocFilterExclude, true},
		{"Valid only", AdHocFilterOnly, true},
		{"Invalid filter", AdHocFilterType("invalid"), false},
		{"Empty filter", AdHocFilterType(""), false},
		{"Case sensitive - wrong case", AdHocFilterType("Include"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.IsValid(); got != tt.expected {
				t.Errorf("AdHocFilterType.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseAdHocFilterType(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  AdHocFilterType
		expectErr bool
	}{
		{"Valid include", "include", AdHocFilterInclude, false},
		{"Valid exclude", "exclude", AdHocFilterExclude, false},
		{"Valid only", "only", AdHocFilterOnly, false},
		{"Invalid filter", "invalid", AdHocFilterType(""), true},
		{"Empty string", "", AdHocFilterType(""), true},
		{"Case sensitive - uppercase", "INCLUDE", AdHocFilterType(""), true},
		{"Case sensitive - mixed case", "Include", AdHocFilterType(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAdHocFilterType(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ParseAdHocFilterType() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ParseAdHocFilterType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAdHocFilterTypeConstants(t *testing.T) {
	// Test that the constants have the expected string values
	tests := []struct {
		name     string
		filter   AdHocFilterType
		expected string
	}{
		{"Include constant", AdHocFilterInclude, "include"},
		{"Exclude constant", AdHocFilterExclude, "exclude"},
		{"Only constant", AdHocFilterOnly, "only"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.filter) != tt.expected {
				t.Errorf("AdHocFilterType constant %s = %v, want %v", tt.name, string(tt.filter), tt.expected)
			}
		})
	}
}

func TestParseAdHocFilterType_ErrorMessage(t *testing.T) {
	// Test that error messages are descriptive
	_, err := ParseAdHocFilterType("invalid")
	if err == nil {
		t.Fatal("ParseAdHocFilterType() should return error for invalid input")
	}

	expectedMessage := "invalid ad-hoc filter type: invalid (must be one of: include, exclude, only)"
	if err.Error() != expectedMessage {
		t.Errorf("ParseAdHocFilterType() error message = %v, want %v", err.Error(), expectedMessage)
	}
}