package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CSVPathError represents different types of CSV path validation errors
type CSVPathError struct {
	Path    string
	Type    string
	Message string
}

func (e CSVPathError) Error() string {
	return e.Message
}

// ValidateCSVPath performs comprehensive validation of a CSV file path
func ValidateCSVPath(path string) error {
	// Check if path is empty
	if strings.TrimSpace(path) == "" {
		return CSVPathError{
			Path:    path,
			Type:    "empty",
			Message: "CSV file path cannot be empty",
		}
	}

	// Clean the path
	cleanPath := filepath.Clean(path)

	// Check if path exists
	info, err := os.Stat(cleanPath)
	if os.IsNotExist(err) {
		return CSVPathError{
			Path:    cleanPath,
			Type:    "not_found",
			Message: fmt.Sprintf("File '%s' does not exist", cleanPath),
		}
	}
	if err != nil {
		return CSVPathError{
			Path:    cleanPath,
			Type:    "access_error",
			Message: fmt.Sprintf("Cannot access '%s': %v", cleanPath, err),
		}
	}

	// Check if it's a directory
	if info.IsDir() {
		return CSVPathError{
			Path:    cleanPath,
			Type:    "is_directory",
			Message: fmt.Sprintf("'%s' is a directory, not a file. Please specify a CSV file, e.g., '%s/data.csv'", cleanPath, cleanPath),
		}
	}

	// Check if file is readable
	file, err := os.Open(cleanPath)
	if err != nil {
		return CSVPathError{
			Path:    cleanPath,
			Type:    "not_readable",
			Message: fmt.Sprintf("Cannot read file '%s': %v", cleanPath, err),
		}
	}
	file.Close()

	// Check file extension (warning, not error)
	ext := strings.ToLower(filepath.Ext(cleanPath))
	if ext != ".csv" && ext != ".txt" {
		// This is just a warning, not an error
		fmt.Printf("⚠️  Warning: File '%s' doesn't have a .csv or .txt extension. Proceeding anyway...\n", cleanPath)
	}

	// Check file size (warn if empty)
	if info.Size() == 0 {
		return CSVPathError{
			Path:    cleanPath,
			Type:    "empty_file",
			Message: fmt.Sprintf("File '%s' is empty", cleanPath),
		}
	}

	// Check if file looks like it might be a CSV (basic check)
	if err := validateCSVFormat(cleanPath); err != nil {
		return CSVPathError{
			Path:    cleanPath,
			Type:    "invalid_format",
			Message: fmt.Sprintf("File '%s' doesn't appear to be a valid CSV: %v", cleanPath, err),
		}
	}

	return nil
}

// validateCSVFormat does a basic check to see if the file looks like a CSV
func validateCSVFormat(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read first few bytes to check for common issues
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && n == 0 {
		return fmt.Errorf("cannot read file content")
	}

	content := string(buffer[:n])

	// Check for binary file (contains null bytes)
	if strings.Contains(content, "\x00") {
		return fmt.Errorf("file appears to be binary, not text")
	}

	// Check if it has some delimiter characters (basic heuristic)
	hasCommas := strings.Contains(content, ",")
	hasTabs := strings.Contains(content, "\t")
	hasSemicolons := strings.Contains(content, ";")

	if !hasCommas && !hasTabs && !hasSemicolons {
		return fmt.Errorf("file doesn't contain common CSV delimiters (comma, tab, or semicolon)")
	}

	// Check for headers that look like kanban data
	firstLine := strings.Split(content, "\n")[0]
	if firstLine != "" {
		lowerFirstLine := strings.ToLower(firstLine)
		// Look for some expected column names
		expectedColumns := []string{"id", "name", "estimate", "completed"}
		foundColumns := 0
		for _, col := range expectedColumns {
			if strings.Contains(lowerFirstLine, col) {
				foundColumns++
			}
		}
		
		if foundColumns == 0 {
			fmt.Printf("⚠️  Warning: File doesn't contain expected kanban columns (id, name, estimate, etc.). Proceeding anyway...\n")
		}
	}

	return nil
}

// SuggestCSVFiles suggests CSV files in a directory if user provided a directory
func SuggestCSVFiles(dirPath string) []string {
	var suggestions []string
	
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return suggestions
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			ext := strings.ToLower(filepath.Ext(name))
			if ext == ".csv" || ext == ".txt" {
				suggestions = append(suggestions, filepath.Join(dirPath, name))
			}
		}
	}

	return suggestions
}