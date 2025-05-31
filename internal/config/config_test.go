package config

import (
	"flag"
	"os"
	"strings"
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

func TestParseFlags_ErrorHandling(t *testing.T) {
	// Save original command line arguments and restore after test
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	oldFlagCommandLine := flag.CommandLine
	defer func() { flag.CommandLine = oldFlagCommandLine }()

	testCases := []struct {
		name      string
		args      []string
		expectErr bool
		errorMsg  string
	}{
		{
			name:      "Invalid report type",
			args:      []string{"cmd", "--csv", "test.csv", "--type", "invalid-type"},
			expectErr: true,
			errorMsg:  "invalid report type",
		},
		{
			name:      "Invalid metrics type", 
			args:      []string{"cmd", "--csv", "test.csv", "--metrics", "invalid-metrics"},
			expectErr: true,
			errorMsg:  "invalid report type", // Note: config.go reuses report type error
		},
		{
			name:      "Invalid period type",
			args:      []string{"cmd", "--csv", "test.csv", "--metrics", "lead-time", "--period", "invalid"},
			expectErr: true,
			errorMsg:  "invalid period type",
		},
		{
			name:      "Invalid ad-hoc filter",
			args:      []string{"cmd", "--csv", "test.csv", "--type", "contributor", "--ad-hoc", "invalid"},
			expectErr: true,
			errorMsg:  "invalid ad-hoc filter type",
		},
		{
			name:      "Invalid filter field",
			args:      []string{"cmd", "--csv", "test.csv", "--type", "contributor", "--filter-field", "invalid"},
			expectErr: true,
			errorMsg:  "invalid filter field",
		},
		{
			name:      "Malformed start date",
			args:      []string{"cmd", "--csv", "test.csv", "--type", "contributor", "--start", "not-a-date"},
			expectErr: true,
			errorMsg:  "error parsing start date",
		},
		{
			name:      "Malformed end date",
			args:      []string{"cmd", "--csv", "test.csv", "--type", "contributor", "--end", "2024-13-45"}, // Invalid month/day
			expectErr: true,
			errorMsg:  "error parsing end date",
		},
		{
			name:      "End date before start date",
			args:      []string{"cmd", "--csv", "test.csv", "--type", "contributor", "--start", "2024-05-31", "--end", "2024-05-01"},
			expectErr: true,
			errorMsg:  "invalid date range",
		},
		{
			name:      "Negative last N days",
			args:      []string{"cmd", "--csv", "test.csv", "--type", "contributor", "--last", "-5"},
			expectErr: true,
			errorMsg:  "last N days must be a positive number",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset flag CommandLine
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Set command line args for this test
			os.Args = tc.args

			// Call ParseFlags
			_, err := ParseFlags()

			// Check for expected error
			if (err != nil) != tc.expectErr {
				t.Errorf("Expected error: %v, got: %v", tc.expectErr, err != nil)
				return
			}

			if tc.expectErr && err != nil {
				if !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("Expected error containing '%s', got: %s", tc.errorMsg, err.Error())
				}
			}
		})
	}
}

func TestParseFlags_DefaultBehavior(t *testing.T) {
	// Test default values and behavior when optional flags are not provided
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	oldFlagCommandLine := flag.CommandLine
	defer func() { flag.CommandLine = oldFlagCommandLine }()

	testCases := []struct {
		name     string
		args     []string
		validate func(*Config) bool
	}{
		{
			name: "Default delimiter should be auto",
			args: []string{"cmd", "--csv", "test.csv", "--type", "contributor"},
			validate: func(cfg *Config) bool {
				return cfg.Delimiter.AutoDetect == true && cfg.Delimiter.Name == "auto"
			},
		},
		{
			name: "Default ad-hoc filter should be include",
			args: []string{"cmd", "--csv", "test.csv", "--type", "contributor"},
			validate: func(cfg *Config) bool {
				return cfg.AdHocFilter == "include"
			},
		},
		{
			name: "Default filter field should be completed_at",
			args: []string{"cmd", "--csv", "test.csv", "--type", "contributor"},
			validate: func(cfg *Config) bool {
				return cfg.FilterField == "completed_at"
			},
		},
		{
			name: "Default period type should be month",
			args: []string{"cmd", "--csv", "test.csv", "--metrics", "lead-time"},
			validate: func(cfg *Config) bool {
				return cfg.PeriodType == "month"
			},
		},
		{
			name: "Invalid delimiter falls back to auto",
			args: []string{"cmd", "--csv", "test.csv", "--type", "contributor", "--delimiter", "invalid"},
			validate: func(cfg *Config) bool {
				return cfg.Delimiter.AutoDetect == true
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
			if err != nil {
				t.Fatalf("ParseFlags() error = %v", err)
			}

			if !tc.validate(cfg) {
				t.Errorf("Validation failed for test: %s", tc.name)
			}
		})
	}
}

func TestParseFlags_EdgeCaseBehavior(t *testing.T) {
	// Test behavior with edge case inputs
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	oldFlagCommandLine := flag.CommandLine
	defer func() { flag.CommandLine = oldFlagCommandLine }()

	testCases := []struct {
		name     string
		args     []string
		validate func(*Config) bool
	}{
		{
			name: "Both type and metrics specified - metrics should win",
			args: []string{"cmd", "--csv", "test.csv", "--type", "contributor", "--metrics", "lead-time"},
			validate: func(cfg *Config) bool {
				return cfg.IsMetricsReport() && cfg.MetricsType == "lead-time"
			},
		},
		{
			name: "Last N days takes precedence over explicit dates",
			args: []string{"cmd", "--csv", "test.csv", "--type", "contributor", "--start", "2024-05-01", "--end", "2024-05-31", "--last", "7"},
			validate: func(cfg *Config) bool {
				return cfg.LastNDays == 7 && !cfg.StartDate.IsZero() && !cfg.EndDate.IsZero()
			},
		},
		{
			name: "Zero last N days should not override explicit dates",
			args: []string{"cmd", "--csv", "test.csv", "--type", "contributor", "--start", "2024-05-01", "--last", "0"},
			validate: func(cfg *Config) bool {
				return cfg.LastNDays == 0 && !cfg.StartDate.IsZero()
			},
		},
		{
			name: "End date should be set to end of day",
			args: []string{"cmd", "--csv", "test.csv", "--type", "contributor", "--end", "2024-05-31"},
			validate: func(cfg *Config) bool {
				// Should be 2024-05-31 23:59:59
				return cfg.EndDate.Hour() == 23 && cfg.EndDate.Minute() == 59 && cfg.EndDate.Second() == 59
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
			if err != nil {
				t.Fatalf("ParseFlags() error = %v", err)
			}

			if !tc.validate(cfg) {
				t.Errorf("Validation failed for test: %s", tc.name)
			}
		})
	}
}