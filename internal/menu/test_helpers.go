package menu

import (
	"os"
	"testing"
)

// TestHelper provides utilities for menu testing
type TestHelper struct {
	tempFiles []string
}

// NewTestHelper creates a new test helper
func NewTestHelper() *TestHelper {
	return &TestHelper{
		tempFiles: make([]string, 0),
	}
}

// CreateTempCSV creates a temporary CSV file with test data
func (h *TestHelper) CreateTempCSV(t *testing.T, content string) string {
	t.Helper()
	
	if content == "" {
		content = `id,name,estimate,is_completed,completed_at
1,Test Task,3,TRUE,2024/05/01 10:00:00
2,Another Task,2,FALSE,
`
	}
	
	tmpFile, err := os.CreateTemp("", "test-kanban-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	
	tmpFile.Close()
	h.tempFiles = append(h.tempFiles, tmpFile.Name())
	return tmpFile.Name()
}

// CreateTempDir creates a temporary directory for testing directory scenarios
func (h *TestHelper) CreateTempDir(t *testing.T) string {
	t.Helper()
	
	tmpDir, err := os.MkdirTemp("", "test-kanban-dir-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	
	h.tempFiles = append(h.tempFiles, tmpDir)
	return tmpDir
}

// Cleanup removes all temporary files created during testing
func (h *TestHelper) Cleanup() {
	for _, file := range h.tempFiles {
		os.RemoveAll(file) // RemoveAll works for both files and directories
	}
	h.tempFiles = nil
}