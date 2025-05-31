package reports

import (
	"strings"
	"testing"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

func TestGenerateProductAreaReport(t *testing.T) {
	// Create test data
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			ProductArea: "Backend",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    3,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			ProductArea: "Frontend",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    2,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			ProductArea: "Backend",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    4,
		},
		{
			ID:          "4",
			Name:        "Task 4",
			ProductArea: "", // No product area
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    1,
		},
	}

	// Create reporter and generate report
	reporter := NewReporter(items)
	report, err := reporter.generateProductAreaReport(items)
	if err != nil {
		t.Fatalf("generateProductAreaReport() error = %v", err)
	}

	// Verify the report content
	if !strings.Contains(report, "Story Points by Product Area") {
		t.Errorf("Report doesn't contain expected header")
	}

	// Backend should have 3 + 4 = 7 points
	if !strings.Contains(report, "Backend") {
		t.Errorf("Report doesn't contain Backend")
	}

	// Frontend should have 2 points
	if !strings.Contains(report, "Frontend") {
		t.Errorf("Report doesn't contain Frontend")
	}

	// Item without product area should be categorized as "Uncategorized"
	if !strings.Contains(report, "Uncategorized") {
		t.Errorf("Report doesn't contain 'Uncategorized' category")
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

func TestGenerateProductAreaReport_EmptyItems(t *testing.T) {
	// Test with empty items slice
	items := []models.KanbanItem{}

	reporter := NewReporter(items)
	report, err := reporter.generateProductAreaReport(items)
	if err != nil {
		t.Fatalf("generateProductAreaReport() error = %v", err)
	}

	// Should still have header and total
	if !strings.Contains(report, "Story Points by Product Area") {
		t.Errorf("Report doesn't contain expected header")
	}

	if !strings.Contains(report, "Total: 0.0 points across 0 items") {
		t.Errorf("Report doesn't contain correct totals for empty items")
	}
}

func TestGenerateProductAreaReport_SortingByPoints(t *testing.T) {
	// Create test data with different point values to test sorting
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			ProductArea: "Mobile",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    2,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			ProductArea: "Infrastructure",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    10,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			ProductArea: "API",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    5,
		},
	}

	reporter := NewReporter(items)
	report, err := reporter.generateProductAreaReport(items)
	if err != nil {
		t.Fatalf("generateProductAreaReport() error = %v", err)
	}

	// Find positions of each product area in the report
	infraPos := strings.Index(report, "Infrastructure")
	apiPos := strings.Index(report, "API")
	mobilePos := strings.Index(report, "Mobile")

	// Infrastructure (10 points) should come first
	// API (5 points) should come second
	// Mobile (2 points) should come last
	if infraPos == -1 || apiPos == -1 || mobilePos == -1 {
		t.Fatalf("Not all product areas found in report")
	}

	if infraPos > apiPos {
		t.Errorf("Infrastructure should appear before API (sorting by points descending)")
	}

	if apiPos > mobilePos {
		t.Errorf("API should appear before Mobile (sorting by points descending)")
	}
}

func TestGenerateProductAreaReport_MultipleItemsSameArea(t *testing.T) {
	// Test aggregation of multiple items in the same product area
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			ProductArea: "Analytics",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    1.5,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			ProductArea: "Analytics",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    2.5,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			ProductArea: "Analytics",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    3.0,
		},
	}

	reporter := NewReporter(items)
	report, err := reporter.generateProductAreaReport(items)
	if err != nil {
		t.Fatalf("generateProductAreaReport() error = %v", err)
	}

	// Should aggregate to 7.0 points and 3 items
	if !strings.Contains(report, "7.0 points") {
		t.Errorf("Report doesn't contain correct aggregated points")
	}

	if !strings.Contains(report, "3 items") {
		t.Errorf("Report doesn't contain correct item count")
	}
}

func TestGenerateProductAreaReport_SpecialCharactersInNames(t *testing.T) {
	// Test with product areas that have special characters or spaces
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			ProductArea: "Machine Learning & AI",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    5,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			ProductArea: "Data Science/Analytics",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    3,
		},
	}

	reporter := NewReporter(items)
	report, err := reporter.generateProductAreaReport(items)
	if err != nil {
		t.Fatalf("generateProductAreaReport() error = %v", err)
	}

	// Should handle special characters correctly
	if !strings.Contains(report, "Machine Learning & AI") {
		t.Errorf("Report doesn't contain product area with special characters")
	}

	if !strings.Contains(report, "Data Science/Analytics") {
		t.Errorf("Report doesn't contain product area with slash")
	}
}