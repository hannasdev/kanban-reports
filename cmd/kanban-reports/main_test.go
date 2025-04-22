// cmd/kanban-reports/main_test.go
package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestEndToEnd(t *testing.T) {
	// Skip if running in CI environment without the necessary setup
	if os.Getenv("CI") != "" {
		t.Skip("Skipping integration test in CI environment")
	}

	// Create a test CSV file
	tempDir, err := os.MkdirTemp("", "kanban-test-*")
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
3,Task 3,Feature,5,TRUE,2024/05/10 16:30:00,john@example.com;jane@example.com,Epic 2,Team B,Backend,2024/05/03 08:00:00,2024/05/06 09:00:00
4,Task 4,Task,2,FALSE,,bob@example.com,Epic 2,Team B,Backend,2024/05/04 11:00:00,2024/05/08 14:00:00
`
	if err := os.WriteFile(csvPath, []byte(testCSV), 0644); err != nil {
		t.Fatalf("Failed to write test CSV: %v", err)
	}

	// Get the path to the compiled binary
	// Note: You need to build the binary first with `go build -o bin/kanban-reports ./cmd/kanban-reports`
	binaryPath := filepath.Join("..", "..", "bin", "kanban-reports")
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skip("Binary not found at " + binaryPath + ". Run 'go build -o bin/kanban-reports ./cmd/kanban-reports' first")
	}

	// Test cases for different report types and options
	testCases := []struct {
		name    string
		args    []string
		checks  []string
	}{
		{
			name: "Contributor Report",
			args: []string{"--csv", csvPath, "--type", "contributor", "--output", outputPath},
			checks: []string{
				"Story Points by Contributor",
				"john@example.com",
				"jane@example.com",
			},
		},
		{
			name: "Epic Report",
			args: []string{"--csv", csvPath, "--type", "epic", "--output", outputPath},
			checks: []string{
				"Story Points by Epic",
				"Epic 1",
				"Epic 2",
			},
		},
		{
			name: "Lead Time Metrics",
			args: []string{"--csv", csvPath, "--metrics", "lead-time", "--output", outputPath},
			checks: []string{
				"Lead Time Analysis",
				"Creation to Completion",
				"Start to Completion",
			},
		},
		{
			name: "Throughput Metrics",
			args: []string{"--csv", csvPath, "--metrics", "throughput", "--output", outputPath},
			checks: []string{
				"Throughput Analysis",
				"Items Completed",
				"Story Points",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Run the command
			cmd := exec.Command(binaryPath, tc.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Command failed: %v\nOutput: %s", err, output)
			}

			// Read the output file
			reportContent, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}

			// Check for expected content
			for _, check := range tc.checks {
				if !strings.Contains(string(reportContent), check) {
					t.Errorf("Report doesn't contain expected content: %s", check)
				}
			}
		})
	}
}