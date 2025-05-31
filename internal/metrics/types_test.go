package metrics

import (
	"testing"
)

func TestMetricsType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		mt       MetricsType
		expected bool
	}{
		{"Valid lead-time", MetricsTypeLeadTime, true},
		{"Valid throughput", MetricsTypeThroughput, true},
		{"Valid flow", MetricsTypeFlow, true},
		{"Valid estimation", MetricsTypeEstimation, true},
		{"Valid age", MetricsTypeAge, true},
		{"Valid improvement", MetricsTypeImprovement, true},
		{"Valid all", MetricsTypeAll, true},
		{"Invalid type", MetricsType("invalid"), false},
		{"Empty type", MetricsType(""), false},
		{"Case sensitive - wrong case", MetricsType("Lead-Time"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mt.IsValid(); got != tt.expected {
				t.Errorf("MetricsType.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseMetricsType(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  MetricsType
		expectErr bool
	}{
		{"Valid lead-time", "lead-time", MetricsTypeLeadTime, false},
		{"Valid throughput", "throughput", MetricsTypeThroughput, false},
		{"Valid flow", "flow", MetricsTypeFlow, false},
		{"Valid estimation", "estimation", MetricsTypeEstimation, false},
		{"Valid age", "age", MetricsTypeAge, false},
		{"Valid improvement", "improvement", MetricsTypeImprovement, false},
		{"Valid all", "all", MetricsTypeAll, false},
		{"Empty string (valid)", "", MetricsType(""), false}, // Empty is valid for no metrics
		{"Invalid type", "invalid", MetricsType(""), true},
		{"Case sensitive - uppercase", "LEAD-TIME", MetricsType(""), true},
		{"Case sensitive - mixed case", "Lead-Time", MetricsType(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMetricsType(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ParseMetricsType() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ParseMetricsType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMetricsTypeConstants(t *testing.T) {
	// Test that the constants have the expected string values
	tests := []struct {
		name     string
		mt       MetricsType
		expected string
	}{
		{"Lead time constant", MetricsTypeLeadTime, "lead-time"},
		{"Throughput constant", MetricsTypeThroughput, "throughput"},
		{"Flow constant", MetricsTypeFlow, "flow"},
		{"Estimation constant", MetricsTypeEstimation, "estimation"},
		{"Age constant", MetricsTypeAge, "age"},
		{"Improvement constant", MetricsTypeImprovement, "improvement"},
		{"All constant", MetricsTypeAll, "all"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.mt) != tt.expected {
				t.Errorf("MetricsType constant %s = %v, want %v", tt.name, string(tt.mt), tt.expected)
			}
		})
	}
}

func TestPeriodType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		pt       PeriodType
		expected bool
	}{
		{"Valid week", PeriodTypeWeek, true},
		{"Valid month", PeriodTypeMonth, true},
		{"Invalid type", PeriodType("invalid"), false},
		{"Empty type", PeriodType(""), false},
		{"Case sensitive - wrong case", PeriodType("Week"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pt.IsValid(); got != tt.expected {
				t.Errorf("PeriodType.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParsePeriodType(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  PeriodType
		expectErr bool
	}{
		{"Valid week", "week", PeriodTypeWeek, false},
		{"Valid month", "month", PeriodTypeMonth, false},
		{"Invalid type", "invalid", PeriodType(""), true},
		{"Empty string", "", PeriodType(""), true},
		{"Case sensitive - uppercase", "WEEK", PeriodType(""), true},
		{"Case sensitive - mixed case", "Week", PeriodType(""), true},
		{"Valid but unusual", "day", PeriodType(""), true}, // Not supported
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePeriodType(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ParsePeriodType() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ParsePeriodType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPeriodTypeConstants(t *testing.T) {
	// Test that the constants have the expected string values
	tests := []struct {
		name     string
		pt       PeriodType
		expected string
	}{
		{"Week constant", PeriodTypeWeek, "week"},
		{"Month constant", PeriodTypeMonth, "month"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.pt) != tt.expected {
				t.Errorf("PeriodType constant %s = %v, want %v", tt.name, string(tt.pt), tt.expected)
			}
		})
	}
}

func TestParseMetricsType_ErrorMessage(t *testing.T) {
	// Test that error messages are descriptive
	_, err := ParseMetricsType("invalid")
	if err == nil {
		t.Fatal("ParseMetricsType() should return error for invalid input")
	}

	expectedMessage := "invalid report type: invalid"
	if err.Error() != expectedMessage {
		t.Errorf("ParseMetricsType() error message = %v, want %v", err.Error(), expectedMessage)
	}
}

func TestParsePeriodType_ErrorMessage(t *testing.T) {
	// Test that error messages are descriptive
	_, err := ParsePeriodType("invalid")
	if err == nil {
		t.Fatal("ParsePeriodType() should return error for invalid input")
	}

	expectedMessage := "invalid period type: invalid (must be one of: week, month)"
	if err.Error() != expectedMessage {
		t.Errorf("ParsePeriodType() error message = %v, want %v", err.Error(), expectedMessage)
	}
}

func TestMetricsType_StringConversion(t *testing.T) {
	// Test that MetricsType can be properly converted to string
	tests := []struct {
		name     string
		mt       MetricsType
		expected string
	}{
		{"Lead time to string", MetricsTypeLeadTime, "lead-time"},
		{"Throughput to string", MetricsTypeThroughput, "throughput"},
		{"Flow to string", MetricsTypeFlow, "flow"},
		{"Estimation to string", MetricsTypeEstimation, "estimation"},
		{"Age to string", MetricsTypeAge, "age"},
		{"Improvement to string", MetricsTypeImprovement, "improvement"},
		{"All to string", MetricsTypeAll, "all"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(tt.mt)
			if got != tt.expected {
				t.Errorf("string(MetricsType) = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPeriodType_StringConversion(t *testing.T) {
	// Test that PeriodType can be properly converted to string
	tests := []struct {
		name     string
		pt       PeriodType
		expected string
	}{
		{"Week to string", PeriodTypeWeek, "week"},
		{"Month to string", PeriodTypeMonth, "month"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(tt.pt)
			if got != tt.expected {
				t.Errorf("string(PeriodType) = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseMetricsType_EmptyString(t *testing.T) {
	// Test that empty string is valid (represents no metrics)
	got, err := ParseMetricsType("")
	if err != nil {
		t.Errorf("ParseMetricsType() with empty string should not return error, got: %v", err)
	}

	if got != MetricsType("") {
		t.Errorf("ParseMetricsType() with empty string = %v, want empty MetricsType", got)
	}
}

func TestMetricsType_RoundTrip(t *testing.T) {
	// Test that we can parse a string to MetricsType and back
	original := "lead-time"
	
	parsed, err := ParseMetricsType(original)
	if err != nil {
		t.Fatalf("ParseMetricsType() error = %v", err)
	}

	backToString := string(parsed)
	if backToString != original {
		t.Errorf("Round trip failed: %v -> %v -> %v", original, parsed, backToString)
	}
}

func TestPeriodType_RoundTrip(t *testing.T) {
	// Test that we can parse a string to PeriodType and back
	original := "week"
	
	parsed, err := ParsePeriodType(original)
	if err != nil {
		t.Fatalf("ParsePeriodType() error = %v", err)
	}

	backToString := string(parsed)
	if backToString != original {
		t.Errorf("Round trip failed: %v -> %v -> %v", original, parsed, backToString)
	}
}