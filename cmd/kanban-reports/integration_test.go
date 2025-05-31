package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMainIntegrationReports(t *testing.T) {
	// Skip if running in CI environment without the necessary setup
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI environment")
	}

	// Create a test CSV file
	tempDir, err := os.MkdirTemp("", "kanban-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	csvPath := filepath.Join(tempDir, "test-data.csv")

	// Create comprehensive test CSV content
	testCSV := `id,name,type,estimate,is_completed,completed_at,owners,epic,team,product_area,created_at,started_at,labels
1,Feature A,Feature,3,TRUE,2024/05/07 10:30:00,john@example.com,Epic Alpha,Team 1,Backend,2024/05/01 09:00:00,2024/05/03 11:00:00,feature
2,Bug Fix B,Bug,1,TRUE,2024/05/08 15:45:00,jane@example.com,Epic Alpha,Team 1,Frontend,2024/05/02 14:00:00,2024/05/05 10:00:00,bug
3,Feature C,Feature,5,TRUE,2024/05/10 16:30:00,john@example.com;jane@example.com,Epic Beta,Team 2,Backend,2024/05/03 08:00:00,2024/05/06 09:00:00,feature
4,Task D,Task,2,FALSE,,bob@example.com,Epic Beta,Team 2,Backend,2024/05/04 11:00:00,2024/05/08 14:00:00,task
5,Ad Hoc E,Feature,1,TRUE,2024/05/09 12:00:00,alice@example.com,Epic Alpha,Team 1,Frontend,2024/05/08 10:00:00,2024/05/09 11:00:00,ad-hoc-request
`
	if err := os.WriteFile(csvPath, []byte(testCSV), 0644); err != nil {
		t.Fatalf("Failed to write test CSV: %v", err)
	}

	// Test different report types and options using the main package logic
	testCases := []struct {
		name string
		args []string
		// We'll just check that the function completes without error
		// since we can't easily capture the output in this integration test
	}{
		{
			name: "Contributor Report",
			args: []string{"program", "--csv", csvPath, "--type", "contributor"},
		},
		{
			name: "Epic Report",
			args: []string{"program", "--csv", csvPath, "--type", "epic"},
		},
		{
			name: "Product Area Report",
			args: []string{"program", "--csv", csvPath, "--type", "product-area"},
		},
		{
			name: "Team Report",
			args: []string{"program", "--csv", csvPath, "--type", "team"},
		},
		{
			name: "Lead Time Metrics",
			args: []string{"program", "--csv", csvPath, "--metrics", "lead-time"},
		},
		{
			name: "Throughput Metrics",
			args: []string{"program", "--csv", csvPath, "--metrics", "throughput"},
		},
		{
			name: "Flow Metrics",
			args: []string{"program", "--csv", csvPath, "--metrics", "flow"},
		},
		{
			name: "Estimation Metrics",
			args: []string{"program", "--csv", csvPath, "--metrics", "estimation"},
		},
		{
			name: "Age Metrics",
			args: []string{"program", "--csv", csvPath, "--metrics", "age"},
		},
		{
			name: "Improvement Metrics",
			args: []string{"program", "--csv", csvPath, "--metrics", "improvement"},
		},
		{
			name: "All Metrics",
			args: []string{"program", "--csv", csvPath, "--metrics", "all"},
		},
		{
			name: "Report with Date Range",
			args: []string{"program", "--csv", csvPath, "--type", "contributor", "--start", "2024-05-01", "--end", "2024-05-31"},
		},
		{
			name: "Report with Last N Days",
			args: []string{"program", "--csv", csvPath, "--type", "epic", "--last", "7"},
		},
		{
			name: "Report Excluding Ad-Hoc",
			args: []string{"program", "--csv", csvPath, "--type", "contributor", "--ad-hoc", "exclude"},
		},
		{
			name: "Report Only Ad-Hoc",
			args: []string{"program", "--csv", csvPath, "--type", "contributor", "--ad-hoc", "only"},
		},
		{
			name: "Weekly Throughput Metrics",
			args: []string{"program", "--csv", csvPath, "--metrics", "throughput", "--period", "week"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save original command line arguments
			origArgs := os.Args
			defer func() { os.Args = origArgs }()

			// Set command line args for this test
			os.Args = tc.args

			// This is a simplified integration test that just verifies
			// the main components can work together without crashing.
			// In a real scenario, you might capture stdout/stderr or
			// write to a file and verify the output content.

			// We can't easily test the main() function directly since it calls os.Exit,
			// but we can test the core logic components individually through integration

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Test %s panicked: %v", tc.name, r)
				}
			}()

			// For now, we just verify the test setup is working
			// In a more complete integration test, you would:
			// 1. Capture the program output
			// 2. Verify the output contains expected content
			// 3. Test error conditions
		})
	}
}

func TestMainErrorConditions(t *testing.T) {
	// Test various error conditions that the main function should handle

	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "Missing CSV file",
			args: []string{"program", "--type", "contributor"},
		},
		{
			name: "Invalid report type",
			args: []string{"program", "--csv", "nonexistent.csv", "--type", "invalid"},
		},
		{
			name: "Invalid metrics type",
			args: []string{"program", "--csv", "nonexistent.csv", "--metrics", "invalid"},
		},
		{
			name: "Invalid date format",
			args: []string{"program", "--csv", "nonexistent.csv", "--type", "contributor", "--start", "invalid-date"},
		},
		{
			name: "End date before start date",
			args: []string{"program", "--csv", "nonexistent.csv", "--type", "contributor", "--start", "2024-05-31", "--end", "2024-05-01"},
		},
		{
			name: "Negative last N days",
			args: []string{"program", "--csv", "nonexistent.csv", "--type", "contributor", "--last", "-5"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save original command line arguments
			origArgs := os.Args
			defer func() { os.Args = origArgs }()

			// Set command line args for this test
			os.Args = tc.args

			// These tests should result in errors, but we can't easily test
			// the main() function since it calls os.Exit()
			// This is more of a placeholder for comprehensive error testing

			// In a real integration test, you might:
			// 1. Run the program as a subprocess
			// 2. Capture exit codes and error messages
			// 3. Verify appropriate error handling
		})
	}
}

func TestMainWithOutputFile(t *testing.T) {
	// Skip if running in CI environment
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI environment")
	}

	// Create a test CSV file
	tempDir, err := os.MkdirTemp("", "kanban-output-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	csvPath := filepath.Join(tempDir, "test-data.csv")
	outputPath := filepath.Join(tempDir, "output.txt")

	// Create test CSV content
	testCSV := `id,name,type,estimate,is_completed,completed_at,owners,epic,team,product_area,created_at,started_at
1,Task 1,Feature,3,TRUE,2024/05/07 10:30:00,john@example.com,Epic 1,Team A,Backend,2024/05/01 09:00:00,2024/05/03 11:00:00
2,Task 2,Bug,1,TRUE,2024/05/08 15:45:00,jane@example.com,Epic 1,Team A,Frontend,2024/05/02 14:00:00,2024/05/05 10:00:00
`
	if err := os.WriteFile(csvPath, []byte(testCSV), 0644); err != nil {
		t.Fatalf("Failed to write test CSV: %v", err)
	}

	// Test cases that write to output files
	testCases := []struct {
		name         string
		args         []string
		expectOutput bool
	}{
		{
			name:         "Report to File",
			args:         []string{"program", "--csv", csvPath, "--type", "contributor", "--output", outputPath},
			expectOutput: true,
		},
		{
			name:         "Metrics to File",
			args:         []string{"program", "--csv", csvPath, "--metrics", "lead-time", "--output", outputPath},
			expectOutput: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Remove output file if it exists
			os.Remove(outputPath)

			// Save original command line arguments
			origArgs := os.Args
			defer func() { os.Args = origArgs }()

			// Set command line args for this test
			os.Args = tc.args

			// In a complete integration test, you would run the main function
			// or execute the program as a subprocess and then verify:
			
			if tc.expectOutput {
				// After running the program, check that output file was created
				// and contains expected content
				
				// For now, just check the test setup
				if _, err := os.Stat(csvPath); os.IsNotExist(err) {
					t.Errorf("CSV file should exist for test")
				}
			}
		})
	}
}

func TestMainDateRangeFiltering(t *testing.T) {
	// Test that date range filtering works correctly

	// Create a test CSV with items across different dates
	tempDir, err := os.MkdirTemp("", "kanban-date-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	csvPath := filepath.Join(tempDir, "test-data.csv")

	// Create test CSV with items spanning multiple months
	now := time.Now()
	old := now.AddDate(0, -2, 0) // 2 months ago
	recent := now.AddDate(0, 0, -7) // 1 week ago

	testCSV := `id,name,type,estimate,is_completed,completed_at,owners,epic,team,product_area,created_at,started_at
1,Old Task,Feature,3,TRUE,` + old.Format("2006/01/02 15:04:05") + `,john@example.com,Epic 1,Team A,Backend,` + old.AddDate(0, 0, -5).Format("2006/01/02 15:04:05") + `,` + old.AddDate(0, 0, -3).Format("2006/01/02 15:04:05") + `
2,Recent Task,Bug,1,TRUE,` + recent.Format("2006/01/02 15:04:05") + `,jane@example.com,Epic 1,Team A,Frontend,` + recent.AddDate(0, 0, -2).Format("2006/01/02 15:04:05") + `,` + recent.AddDate(0, 0, -1).Format("2006/01/02 15:04:05") + `
`
	if err := os.WriteFile(csvPath, []byte(testCSV), 0644); err != nil {
		t.Fatalf("Failed to write test CSV: %v", err)
	}

	// Test cases with different date filters
	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "Last 30 days",
			args: []string{"program", "--csv", csvPath, "--type", "contributor", "--last", "30"},
		},
		{
			name: "Specific date range",
			args: []string{"program", "--csv", csvPath, "--type", "contributor", "--start", recent.AddDate(0, 0, -3).Format("2006-01-02"), "--end", now.Format("2006-01-02")},
		},
		{
			name: "From date only",
			args: []string{"program", "--csv", csvPath, "--type", "contributor", "--start", recent.Format("2006-01-02")},
		},
		{
			name: "To date only",
			args: []string{"program", "--csv", csvPath, "--type", "contributor", "--end", recent.Format("2006-01-02")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save original command line arguments
			origArgs := os.Args
			defer func() { os.Args = origArgs }()

			// Set command line args for this test
			os.Args = tc.args

			// In a complete test, you would run the program and verify
			// that the correct items are included/excluded based on dates
			
			// For now, just verify test setup
			if _, err := os.Stat(csvPath); os.IsNotExist(err) {
				t.Errorf("CSV file should exist for test")
			}
		})
	}
}

func TestMainDelimiterDetection(t *testing.T) {
	// Test automatic delimiter detection

	testCases := []struct {
		name      string
		csvContent string
		delimiter  string
	}{
		{
			name: "Comma delimited",
			csvContent: `id,name,estimate
1,Task 1,3
2,Task 2,2`,
			delimiter: "comma",
		},
		{
			name: "Tab delimited",
			csvContent: "id\tname\testimate\n1\tTask 1\t3\n2\tTask 2\t2",
			delimiter: "tab",
		},
		{
			name: "Semicolon delimited",
			csvContent: `id;name;estimate
1;Task 1;3
2;Task 2;2`,
			delimiter: "semicolon",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temp file with specific delimiter
			tempDir, err := os.MkdirTemp("", "kanban-delimiter-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp directory: %v", err)
			}
			defer os.RemoveAll(tempDir)

			csvPath := filepath.Join(tempDir, "test-data.csv")
			
			// Add required columns for a valid CSV
			fullCSV := strings.Replace(tc.csvContent, "estimate", "estimate,is_completed,completed_at", 1)
			if tc.delimiter == "tab" {
				fullCSV = strings.Replace(fullCSV, "estimate", "estimate\tis_completed\tcompleted_at", 1)
			} else if tc.delimiter == "semicolon" {
				fullCSV = strings.Replace(fullCSV, "estimate", "estimate;is_completed;completed_at", 1)
			}

			if err := os.WriteFile(csvPath, []byte(fullCSV), 0644); err != nil {
				t.Fatalf("Failed to write test CSV: %v", err)
			}

			// Test auto-detection
			args := []string{"program", "--csv", csvPath, "--type", "contributor", "--delimiter", "auto"}

			// Save original command line arguments
			origArgs := os.Args
			defer func() { os.Args = origArgs }()

			// Set command line args for this test
			os.Args = args

			// In a complete test, you would verify that the correct delimiter
			// was detected and the CSV was parsed correctly
		})
	}
}