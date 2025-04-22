// internal/parser/csv_parser_test.go
package parser

import (
	"os"
	"testing"
)

func TestCSVParser_Parse(t *testing.T) {
	// Create a temporary CSV file for testing
	tempFile, err := os.CreateTemp("", "test-kanban-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data to the file
	testCSV := `id,name,type,estimate,is_completed,completed_at,owners,epic,team,product_area
1,Task 1,Feature,3,TRUE,2024/05/07 10:30:00,john@example.com,Epic 1,Team A,Backend
2,Task 2,Bug,1,TRUE,2024/05/08 15:45:00,jane@example.com,Epic 1,Team A,Frontend
3,Task 3,Task,5,FALSE,,john@example.com;jane@example.com,Epic 2,Team B,Backend
`
	if _, err := tempFile.Write([]byte(testCSV)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Parse the CSV file
	parser := NewCSVParser(tempFile.Name())
	items, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Verify the results
	if len(items) != 3 {
		t.Errorf("Parse() returned %d items, want 3", len(items))
	}

	// Check the first item
	if items[0].ID != "1" || items[0].Name != "Task 1" || items[0].Type != "Feature" {
		t.Errorf("Parse() first item = %+v, want ID=1, Name=Task 1, Type=Feature", items[0])
	}

	// Check completed status
	if !items[0].IsCompleted || items[2].IsCompleted {
		t.Errorf("Parse() completed status incorrect")
	}

	// Check owners parsing
	if len(items[2].Owners) != 2 {
		t.Errorf("Parse() owners = %v, want 2 owners", items[2].Owners)
	}
}