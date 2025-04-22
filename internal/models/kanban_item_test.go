package models

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name     string
		timeStr  string
		wantTime time.Time
		wantErr  bool
	}{
		{
			name:     "Valid timestamp",
			timeStr:  "2024/05/07 03:49:34",
			wantTime: time.Date(2024, 5, 7, 3, 49, 34, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "Empty string",
			timeStr:  "",
			wantTime: time.Time{},
			wantErr:  false,
		},
		{
			name:     "Invalid format",
			timeStr:  "07-05-2024",
			wantTime: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, err := ParseTime(tt.timeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !gotTime.Equal(tt.wantTime) {
				t.Errorf("ParseTime() = %v, want %v", gotTime, tt.wantTime)
			}
		})
	}
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		name     string
		boolStr  string
		expected bool
	}{
		{"TRUE value", "TRUE", true},
		{"True value", "True", true},
		{"true value", "true", true},
		{"FALSE value", "FALSE", false},
		{"False value", "False", false},
		{"false value", "false", false},
		{"Empty string", "", false},
		{"Random string", "random", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseBool(tt.boolStr); got != tt.expected {
				t.Errorf("ParseBool(%q) = %v, want %v", tt.boolStr, got, tt.expected)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		name      string
		floatStr  string
		expected  float64
	}{
		{"Valid float", "3.14", 3.14},
		{"Integer as float", "42", 42.0},
		{"Empty string", "", 0.0},
		{"Invalid float", "not-a-number", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseFloat(tt.floatStr); got != tt.expected {
				t.Errorf("ParseFloat(%q) = %v, want %v", tt.floatStr, got, tt.expected)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name    string
		intStr  string
		expected int
	}{
		{"Valid integer", "42", 42},
		{"Empty string", "", 0},
		{"Invalid integer", "not-a-number", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseInt(tt.intStr); got != tt.expected {
				t.Errorf("ParseInt(%q) = %v, want %v", tt.intStr, got, tt.expected)
			}
		})
	}
}

func TestParseStringList(t *testing.T) {
	tests := []struct {
		name    string
		listStr string
		want    []string
	}{
		{
			name:    "Comma separated list",
			listStr: "one,two,three",
			want:    []string{"one", "two", "three"},
		},
		{
			name:    "Single item",
			listStr: "just-one",
			want:    []string{"just-one"},
		},
		{
			name:    "Empty string",
			listStr: "",
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseStringList(tt.listStr)
			if len(got) != len(tt.want) {
				t.Errorf("ParseStringList() = %v, want %v", got, tt.want)
				return
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("ParseStringList()[%d] = %v, want %v", i, v, tt.want[i])
				}
			}
		})
	}
}

func TestParseOwners(t *testing.T) {
	tests := []struct {
		name      string
		ownersStr string
		want      []string
	}{
		{
			name:      "Comma separated",
			ownersStr: "john.doe@example.com,jane.smith@example.com",
			want:      []string{"john.doe@example.com", "jane.smith@example.com"},
		},
		{
			name:      "Semicolon separated",
			ownersStr: "john.doe@example.com;jane.smith@example.com",
			want:      []string{"john.doe@example.com", "jane.smith@example.com"},
		},
		{
			name:      "Space separated",
			ownersStr: "john.doe@example.com jane.smith@example.com",
			want:      []string{"john.doe@example.com", "jane.smith@example.com"},
		},
		{
			name:      "Single owner",
			ownersStr: "john.doe@example.com",
			want:      []string{"john.doe@example.com"},
		},
		{
			name:      "Empty string",
			ownersStr: "",
			want:      []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseOwners(tt.ownersStr)
			if len(got) != len(tt.want) {
				t.Errorf("ParseOwners(%q) = %v, want %v", tt.ownersStr, got, tt.want)
				return
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("ParseOwners(%q)[%d] = %v, want %v", tt.ownersStr, i, v, tt.want[i])
				}
			}
		})
	}
}

func TestParseCustomFields(t *testing.T) {
	tests := []struct {
		name            string
		customFieldsStr string
		want            map[string]string
	}{
		{
			name:            "Valid custom fields",
			customFieldsStr: "key1=value1;key2=value2",
			want:            map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name:            "Single field",
			customFieldsStr: "key=value",
			want:            map[string]string{"key": "value"},
		},
		{
			name:            "Empty string",
			customFieldsStr: "",
			want:            map[string]string{},
		},
		{
			name:            "Malformed field",
			customFieldsStr: "key1=value1;malformed",
			want:            map[string]string{"key1": "value1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseCustomFields(tt.customFieldsStr)
			if len(got) != len(tt.want) {
				t.Errorf("ParseCustomFields() = %v, want %v", got, tt.want)
				return
			}
			for k, v := range got {
				if wantV, ok := tt.want[k]; !ok || v != wantV {
					t.Errorf("ParseCustomFields()[%q] = %v, want %v", k, v, wantV)
				}
			}
		})
	}
}