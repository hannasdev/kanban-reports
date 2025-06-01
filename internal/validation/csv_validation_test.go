package validation

import (
	"os"
	"strings"
	"testing"
)

func TestValidateCSVPath(t *testing.T) {
	// Create a temporary valid CSV file
	validCSV, err := os.CreateTemp("", "valid-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(validCSV.Name())
	
	// Write valid CSV content
	validContent := `id,name,estimate,is_completed,completed_at
1,Test Task,3,TRUE,2024/05/01 10:00:00
`
	if _, err := validCSV.WriteString(validContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	validCSV.Close()
	
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test-dir-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create an empty file
	emptyFile, err := os.CreateTemp("", "empty-*.csv")
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}
	emptyFile.Close()
	defer os.Remove(emptyFile.Name())
	
	tests := []struct {
		name      string
		path      string
		wantError bool
		errorType string
	}{
		{
			name:      "Valid CSV file",
			path:      validCSV.Name(),
			wantError: false,
		},
		{
			name:      "Empty path",
			path:      "",
			wantError: true,
			errorType: "empty",
		},
		{
			name:      "Nonexistent file",
			path:      "/nonexistent/file.csv",
			wantError: true,
			errorType: "not_found",
		},
		{
			name:      "Directory instead of file",
			path:      tempDir,
			wantError: true,
			errorType: "is_directory",
		},
		{
			name:      "Empty file",
			path:      emptyFile.Name(),
			wantError: true,
			errorType: "empty_file",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCSVPath(tt.path)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateCSVPath(%q) expected error, got nil", tt.path)
					return
				}
				
				csvErr, ok := err.(CSVPathError)
				if !ok {
					t.Errorf("Expected CSVPathError, got %T", err)
					return
				}
				
				if csvErr.Type != tt.errorType {
					t.Errorf("Expected error type %q, got %q", tt.errorType, csvErr.Type)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateCSVPath(%q) expected no error, got %v", tt.path, err)
				}
			}
		})
	}
}

func TestSuggestCSVFiles(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test-suggestions-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create some test files
	testFiles := []struct {
		name   string
		isCSV  bool
	}{
		{"data.csv", true},
		{"export.csv", true},
		{"readme.txt", true},  // .txt files are also suggested
		{"image.png", false},
		{"document.pdf", false},
		{"script.sh", false},
	}
	
	expectedSuggestions := 0
	for _, file := range testFiles {
		filePath := tempDir + "/" + file.name
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		if file.isCSV {
			expectedSuggestions++
		}
	}
	
	suggestions := SuggestCSVFiles(tempDir)
	
	if len(suggestions) != expectedSuggestions {
		t.Errorf("Expected %d suggestions, got %d", expectedSuggestions, len(suggestions))
	}
	
	// Check that all suggestions are valid paths
	for _, suggestion := range suggestions {
		if !strings.HasPrefix(suggestion, tempDir) {
			t.Errorf("Suggestion %q doesn't start with directory path %q", suggestion, tempDir)
		}
		
		ext := strings.ToLower(suggestion[strings.LastIndex(suggestion, "."):])
		if ext != ".csv" && ext != ".txt" {
			t.Errorf("Suggestion %q has unexpected extension %q", suggestion, ext)
		}
	}
}

func TestCSVPathError(t *testing.T) {
	err := CSVPathError{
		Path:    "/test/path",
		Type:    "test_type",
		Message: "test message",
	}
	
	if err.Error() != "test message" {
		t.Errorf("Expected error message 'test message', got %q", err.Error())
	}
}