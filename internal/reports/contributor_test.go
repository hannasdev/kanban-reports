// internal/reports/contributor_test.go
package reports

import (
	"strings"
	"testing"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

func TestGenerateContributorReport(t *testing.T) {
	// Create test data
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Owners:      []string{"john@example.com"},
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    3,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Owners:      []string{"jane@example.com"},
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    2,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			Owners:      []string{"john@example.com", "jane@example.com"},
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    4,
		},
	}

	// Create reporter and generate report
	reporter := NewReporter(items)
	report, err := reporter.generateContributorReport(items)
	if err != nil {
		t.Fatalf("generateContributorReport() error = %v", err)
	}

	// Verify the report content
	if !strings.Contains(report, "Story Points by Contributor") {
		t.Errorf("Report doesn't contain expected header")
	}

	// Points should be: john = 3 + (4/2) = 5, jane = 2 + (4/2) = 4
	if !strings.Contains(report, "john@example.com") || !strings.Contains(report, "jane@example.com") {
		t.Errorf("Report doesn't contain expected contributors")
	}

	// Verify total points
	if !strings.Contains(report, "Total: 9.0 points") {
		t.Errorf("Report doesn't contain correct total points")
	}
}