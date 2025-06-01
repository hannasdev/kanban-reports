package config

import (
	"os"
	"testing"
)

// Helper function to create a temporary valid CSV file for tests
func createTestCSVFile(t *testing.T) string {
	t.Helper()
	
	tempFile, err := os.CreateTemp("", "test-kanban-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	
	testCSV := `id,name,estimate,is_completed,completed_at
1,Test Task,3,TRUE,2024/05/01 10:00:00
`
	if _, err := tempFile.Write([]byte(testCSV)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()
	
	// Register cleanup
	t.Cleanup(func() {
		os.Remove(tempFile.Name())
	})
	
	return tempFile.Name()
}