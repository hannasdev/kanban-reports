package reports

import (
	"strings"
	"testing"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

func TestGenerateEpicReport(t *testing.T) {
	// Create test data
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Epic:        "Epic Alpha",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    3,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Epic:        "Epic Beta",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    2,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			Epic:        "Epic Alpha",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    4,
		},
		{
			ID:          "4",
			Name:        "Task 4",
			Epic:        "", // No epic
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    1,
		},
	}

	// Create reporter and generate report
	reporter := NewReporter(items)
	report, err := reporter.generateEpicReport(items)
	if err != nil {
		t.Fatalf("generateEpicReport() error = %v", err)
	}

	// Verify the report content
	if !strings.Contains(report, "Story Points by Epic") {
		t.Errorf("Report doesn't contain expected header")
	}

	// Epic Alpha should have 3 + 4 = 7 points
	if !strings.Contains(report, "Epic Alpha") {
		t.Errorf("Report doesn't contain Epic Alpha")
	}

	// Epic Beta should have 2 points
	if !strings.Contains(report, "Epic Beta") {
		t.Errorf("Report doesn't contain Epic Beta")
	}

	// Item without epic should be categorized as "No Epic"
	if !strings.Contains(report, "No Epic") {
		t.Errorf("Report doesn't contain 'No Epic' category")
	}

	// Verify total points
	if !strings.Contains(report, "Total: 10.0 points") {
		t.Errorf("Report doesn't contain correct total points")
	}

	// Verify total items
	if !strings.Contains(report, "across 4 items") {
		t.Errorf("Report doesn't contain correct total items")
	}
}

func TestGenerateEpicReport_EmptyItems(t *testing.T) {
	// Test with empty items slice
	items := []models.KanbanItem{}

	reporter := NewReporter(items)
	report, err := reporter.generateEpicReport(items)
	if err != nil {
		t.Fatalf("generateEpicReport() error = %v", err)
	}

	// Should still have header and total
	if !strings.Contains(report, "Story Points by Epic") {
		t.Errorf("Report doesn't contain expected header")
	}

	if !strings.Contains(report, "Total: 0.0 points across 0 items") {
		t.Errorf("Report doesn't contain correct totals for empty items")
	}
}

func TestGenerateEpicReport_SortingByPoints(t *testing.T) {
	// Create test data with different point values to test sorting
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Epic:        "Epic Small",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    2,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Epic:        "Epic Large",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    10,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			Epic:        "Epic Medium",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    5,
		},
	}

	reporter := NewReporter(items)
	report, err := reporter.generateEpicReport(items)
	if err != nil {
		t.Fatalf("generateEpicReport() error = %v", err)
	}

	// Find positions of each epic in the report
	largePos := strings.Index(report, "Epic Large")
	mediumPos := strings.Index(report, "Epic Medium")
	smallPos := strings.Index(report, "Epic Small")

	// Epic Large (10 points) should come first
	// Epic Medium (5 points) should come second
	// Epic Small (2 points) should come last
	if largePos == -1 || mediumPos == -1 || smallPos == -1 {
		t.Fatalf("Not all epics found in report")
	}

	if largePos > mediumPos {
		t.Errorf("Epic Large should appear before Epic Medium (sorting by points descending)")
	}

	if mediumPos > smallPos {
		t.Errorf("Epic Medium should appear before Epic Small (sorting by points descending)")
	}
}

func TestGenerateEpicReport_MultipleItemsSameEpic(t *testing.T) {
	// Test aggregation of multiple items in the same epic
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Epic:        "Epic Test",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    1.5,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Epic:        "Epic Test",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    2.5,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			Epic:        "Epic Test",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    3.0,
		},
	}

	reporter := NewReporter(items)
	report, err := reporter.generateEpicReport(items)
	if err != nil {
		t.Fatalf("generateEpicReport() error = %v", err)
	}

	// Should aggregate to 7.0 points and 3 items
	if !strings.Contains(report, "7.0 points") {
		t.Errorf("Report doesn't contain correct aggregated points")
	}

	if !strings.Contains(report, "3 items") {
		t.Errorf("Report doesn't contain correct item count")
	}
}