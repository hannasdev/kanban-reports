package reports

import (
	"testing"
)

func TestReportType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		rt       ReportType
		expected bool
	}{
		{"Valid contributor", ReportTypeContributor, true},
		{"Valid epic", ReportTypeEpic, true},
		{"Valid product-area", ReportTypeProductArea, true},
		{"Valid team", ReportTypeTeam, true},
		{"Invalid type", ReportType("invalid"), false},
		{"Empty type", ReportType(""), false},
		{"Case sensitive - wrong case", ReportType("Contributor"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rt.IsValid(); got != tt.expected {
				t.Errorf("ReportType.IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseReportType(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  ReportType
		expectErr bool
	}{
		{"Valid contributor", "contributor", ReportTypeContributor, false},
		{"Valid epic", "epic", ReportTypeEpic, false},
		{"Valid product-area", "product-area", ReportTypeProductArea, false},
		{"Valid team", "team", ReportTypeTeam, false},
		{"Invalid type", "invalid", ReportType(""), true},
		{"Empty string", "", ReportType(""), true},
		{"Case sensitive - uppercase", "CONTRIBUTOR", ReportType(""), true},
		{"Case sensitive - mixed case", "Epic", ReportType(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseReportType(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("ParseReportType() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ParseReportType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestReportTypeConstants(t *testing.T) {
	// Test that the constants have the expected string values
	tests := []struct {
		name     string
		rt       ReportType
		expected string
	}{
		{"Contributor constant", ReportTypeContributor, "contributor"},
		{"Epic constant", ReportTypeEpic, "epic"},
		{"Product area constant", ReportTypeProductArea, "product-area"},
		{"Team constant", ReportTypeTeam, "team"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.rt) != tt.expected {
				t.Errorf("ReportType constant %s = %v, want %v", tt.name, string(tt.rt), tt.expected)
			}
		})
	}
}