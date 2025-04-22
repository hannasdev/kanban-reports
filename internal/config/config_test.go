package config

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/hannasdev/kanban-reports/internal/reports"
)

func TestParseFlags(t *testing.T) {
	// Save original command line arguments and restore after test
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	oldFlagCommandLine := flag.CommandLine
	defer func() { flag.CommandLine = oldFlagCommandLine }()

	// Test cases
	testCases := []struct {
		name      string
		args      []string
		expectErr bool
		validate  func(*Config) bool
	}{
		{
			name:      "Missing CSV path",
			args:      []string{"cmd", "--type", "contributor"},
			expectErr: true,
		},
		{
			name:      "Valid contributor report",
			args:      []string{"cmd", "--csv", "test.csv", "--type", "contributor"},
			expectErr: false,
			validate: func(cfg *Config) bool {
				return cfg.CSVPath == "test.csv" && 
				       cfg.ReportType == reports.ReportTypeContributor &&
				       !cfg.IsMetricsReport()
			},
		},
		{
			name:      "Valid metrics report",
			args:      []string{"cmd", "--csv", "test.csv", "--metrics", "lead-time"},
			expectErr: false,
			validate: func(cfg *Config) bool {
				return cfg.CSVPath == "test.csv" && 
				       cfg.IsMetricsReport() &&
				       cfg.MetricsType == "lead-time"
			},
		},
		{
			name:      "Date range with last days",
			args:      []string{"cmd", "--csv", "test.csv", "--type", "epic", "--last", "7"},
			expectErr: false,
			validate: func(cfg *Config) bool {
				now := time.Now()
				startDiff := now.Sub(cfg.StartDate).Hours() / 24
				return startDiff > 6.9 && startDiff < 7.1 // ~7 days
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset flag CommandLine
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Set command line args for this test
			os.Args = tc.args

			// Call ParseFlags
			cfg, err := ParseFlags()

			// Check for expected error
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err != nil)
				return
			}

			// Skip validation if we expected an error
			if tc.expectErr {
				return
			}

			// Validate config if a validation function was provided
			if tc.validate != nil {
				if !tc.validate(cfg) {
					t.Errorf("Config validation failed for test: %s", tc.name)
				}
			}
		})
	}
}

func TestIsMetricsReport(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected bool
	}{
		{
			name:     "Regular report",
			config:   Config{MetricsType: ""},
			expected: false,
		},
		{
			name:     "Metrics report",
			config:   Config{MetricsType: "lead-time"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.IsMetricsReport() != tt.expected {
				t.Errorf("IsMetricsReport() = %v, want %v", tt.config.IsMetricsReport(), tt.expected)
			}
		})
	}
}

func TestGetDateRange(t *testing.T) {
	start := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC)
	
	config := Config{
		StartDate: start,
		EndDate:   end,
	}
	
	gotStart, gotEnd := config.GetDateRange()
	
	if !gotStart.Equal(start) || !gotEnd.Equal(end) {
		t.Errorf("GetDateRange() = %v, %v, want %v, %v", gotStart, gotEnd, start, end)
	}
}