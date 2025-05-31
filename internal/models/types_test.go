package models

import (
	"strings"
	"testing"
	"time"
)

func TestFilterField_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		field    FilterField
		expected bool
	}{
		{"Valid completed_at", FilterFieldCompletedAt, true},
		{"Valid created_at", FilterFieldCreatedAt, true},
		{"Valid started_at", FilterFieldStartedAt, true},
		{"Invalid field", FilterField("invalid"), false},
		{"Empty field", FilterField(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.IsValid(); got != tt.expected {
				t.Errorf("FilterField.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseFilterField(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  FilterField
		expectErr bool
	}{
		{"Valid completed_at", "completed_at", FilterFieldCompletedAt, false},
		{"Valid created_at", "created_at", FilterFieldCreatedAt, false},
		{"Valid started_at", "started_at", FilterFieldStartedAt, false},
		{"Invalid field", "invalid", FilterField(""), true},
		{"Empty string", "", FilterField(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFilterField(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ParseFilterField() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ParseFilterField() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFilterField_GetItemDate(t *testing.T) {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	twoDaysAgo := now.AddDate(0, 0, -2)

	item := KanbanItem{
		IsCompleted: true,
		CreatedAt:   twoDaysAgo,
		StartedAt:   yesterday,
		CompletedAt: now,
	}

	tests := []struct {
		name         string
		field        FilterField
		item         KanbanItem
		expectedDate time.Time
		expectedOk   bool
	}{
		{
			name:         "CompletedAt field with completed item",
			field:        FilterFieldCompletedAt,
			item:         item,
			expectedDate: now,
			expectedOk:   true,
		},
		{
			name:         "CompletedAt field with incomplete item",
			field:        FilterFieldCompletedAt,
			item:         KanbanItem{IsCompleted: false, CompletedAt: now},
			expectedDate: time.Time{},
			expectedOk:   false,
		},
		{
			name:         "CreatedAt field",
			field:        FilterFieldCreatedAt,
			item:         item,
			expectedDate: twoDaysAgo,
			expectedOk:   true,
		},
		{
			name:         "StartedAt field",
			field:        FilterFieldStartedAt,
			item:         item,
			expectedDate: yesterday,
			expectedOk:   true,
		},
		{
			name:         "CreatedAt field with zero time",
			field:        FilterFieldCreatedAt,
			item:         KanbanItem{},
			expectedDate: time.Time{},
			expectedOk:   false,
		},
		{
			name:         "Invalid field defaults to CompletedAt behavior",
			field:        FilterField("invalid"),
			item:         item,
			expectedDate: now,
			expectedOk:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDate, gotOk := tt.field.GetItemDate(tt.item)
			if !gotDate.Equal(tt.expectedDate) {
				t.Errorf("FilterField.GetItemDate() date = %v, want %v", gotDate, tt.expectedDate)
			}
			if gotOk != tt.expectedOk {
				t.Errorf("FilterField.GetItemDate() ok = %v, want %v", gotOk, tt.expectedOk)
			}
		})
	}
}

func TestParseDelimiter(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  DelimiterType
		expectErr bool
	}{
		{"comma", "comma", DelimiterComma, false},
		{"tab", "tab", DelimiterTab, false},
		{"semicolon", "semicolon", DelimiterSemicolon, false},
		{"auto", "auto", DelimiterAuto, false},
		{"invalid", "invalid", DelimiterAuto, true},
		{"empty", "", DelimiterAuto, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDelimiter(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ParseDelimiter() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got.Value != tt.expected.Value || got.Name != tt.expected.Name {
				t.Errorf("ParseDelimiter() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDetectDelimiterType(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected DelimiterType
	}{
		{
			name:     "Comma delimited",
			content:  "name,age,city\nJohn,25,NYC\nJane,30,LA",
			expected: DelimiterComma,
		},
		{
			name:     "Tab delimited",
			content:  "name\tage\tcity\nJohn\t25\tNYC\nJane\t30\tLA",
			expected: DelimiterTab,
		},
		{
			name:     "Semicolon delimited",
			content:  "name;age;city\nJohn;25;NYC\nJane;30;LA",
			expected: DelimiterSemicolon,
		},
		{
			name:     "Mixed delimiters - comma wins",
			content:  "name,age;city\nJohn,25,NYC\nJane,30,LA",
			expected: DelimiterComma,
		},
		{
			name:     "More than 5 lines - only first 5 used",
			content:  strings.Repeat("a,b,c\n", 10),
			expected: DelimiterComma,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectDelimiterType(tt.content)
			if got.Value != tt.expected.Value {
				t.Errorf("DetectDelimiterType() = %v, want %v", got.Value, tt.expected.Value)
			}
		})
	}
}

func TestDetectDelimiterType_EdgeCases(t *testing.T) {
	// Test edge cases where behavior might be non-deterministic
	// We just verify that a valid delimiter is returned
	
	edgeCases := []struct {
		name    string
		content string
	}{
		{"Empty content", ""},
		{"No delimiters", "nameagecity\nJohn25NYC\nJane30LA"},
		{"Single line no delimiters", "singlelinenodelimiters"},
	}
	
	validDelimiters := map[rune]bool{
		',':  true,
		'\t': true,
		';':  true,
	}
	
	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			got := DetectDelimiterType(tc.content)
			
			// Just verify that a valid delimiter is returned
			if !validDelimiters[got.Value] {
				t.Errorf("DetectDelimiterType() returned invalid delimiter: %v", got.Value)
			}
			
			// Verify that the returned DelimiterType is consistent
			if got.AutoDetect != false {
				t.Errorf("DetectDelimiterType() should return non-auto-detect delimiter")
			}
		})
	}
}

func TestDelimiterConstants(t *testing.T) {
	tests := []struct {
		name      string
		delimiter DelimiterType
		wantValue rune
		wantName  string
		wantAuto  bool
	}{
		{"Comma", DelimiterComma, ',', "comma", false},
		{"Tab", DelimiterTab, '\t', "tab", false},
		{"Semicolon", DelimiterSemicolon, ';', "semicolon", false},
		{"Auto", DelimiterAuto, ',', "auto", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.delimiter.Value != tt.wantValue {
				t.Errorf("%s delimiter value = %v, want %v", tt.name, tt.delimiter.Value, tt.wantValue)
			}
			if tt.delimiter.Name != tt.wantName {
				t.Errorf("%s delimiter name = %v, want %v", tt.name, tt.delimiter.Name, tt.wantName)
			}
			if tt.delimiter.AutoDetect != tt.wantAuto {
				t.Errorf("%s delimiter auto = %v, want %v", tt.name, tt.delimiter.AutoDetect, tt.wantAuto)
			}
		})
	}
}