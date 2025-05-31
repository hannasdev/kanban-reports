package metrics

import (
	"strings"
	"testing"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

func TestThroughputReport(t *testing.T) {
	// Create test data with items completed in different months
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Type:        "Feature",
			IsCompleted: true,
			CompletedAt: baseTime.AddDate(0, -1, 0), // April 2024
			Estimate:    3,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Type:        "Bug",
			IsCompleted: true,
			CompletedAt: baseTime.AddDate(0, -1, 5), // April 2024
			Estimate:    2,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			Type:        "Feature",
			IsCompleted: true,
			CompletedAt: baseTime, // May 2024
			Estimate:    5,
		},
		{
			ID:          "4",
			Name:        "Task 4",
			Type:        "Task",
			IsCompleted: true,
			CompletedAt: baseTime.AddDate(0, 0, 3), // May 2024
			Estimate:    1,
		},
	}

	// Test monthly report
	report, err := ThroughputReport(items, "month")
	if err != nil {
		t.Fatalf("ThroughputReport() error = %v", err)
	}

	// Verify the report content
	expectedStrings := []string{
		"Throughput Analysis by Month",
		"What is Throughput?",
		"Items Completed",
		"Story Points",
		"Avg Points/Item",
		"2024-04", // April
		"2024-05", // May
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(report, expected) {
			t.Errorf("Report doesn't contain expected string: %s", expected)
		}
	}

	// Verify breakdown by type section
	if !strings.Contains(report, "Breakdown by Item Type") {
		t.Errorf("Report doesn't contain type breakdown section")
	}

	if !strings.Contains(report, "Feature") || !strings.Contains(report, "Bug") || !strings.Contains(report, "Task") {
		t.Errorf("Report doesn't contain all item types")
	}
}

func TestThroughputReport_WeeklyPeriod(t *testing.T) {
	// Create test data for weekly reporting
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC) // Wednesday
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Type:        "Feature",
			IsCompleted: true,
			CompletedAt: baseTime.AddDate(0, 0, -7), // Previous week
			Estimate:    3,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Type:        "Bug",
			IsCompleted: true,
			CompletedAt: baseTime, // Current week
			Estimate:    2,
		},
	}

	report, err := ThroughputReport(items, "week")
	if err != nil {
		t.Fatalf("ThroughputReport() error = %v", err)
	}

	// Should show weekly format
	if !strings.Contains(report, "Throughput Analysis by Week") {
		t.Errorf("Report doesn't contain weekly header")
	}

	// Should contain ISO week format (2024-W19, 2024-W20, etc.)
	if !strings.Contains(report, "2024-W") {
		t.Errorf("Report doesn't contain ISO week format")
	}
}

func TestThroughputReport_EmptyItems(t *testing.T) {
	items := []models.KanbanItem{}

	report, err := ThroughputReport(items, "month")
	if err != nil {
		t.Fatalf("ThroughputReport() error = %v", err)
	}

	// Should still contain header and explanatory text
	if !strings.Contains(report, "Throughput Analysis by Month") {
		t.Errorf("Report doesn't contain header for empty items")
	}

	if !strings.Contains(report, "What is Throughput?") {
		t.Errorf("Report doesn't contain explanatory text for empty items")
	}
}

func TestThroughputReport_IncompleteItems(t *testing.T) {
	// Create test data with incomplete items (should be excluded)
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Completed Task",
			Type:        "Feature",
			IsCompleted: true,
			CompletedAt: time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC),
			Estimate:    3,
		},
		{
			ID:          "2",
			Name:        "Incomplete Task",
			Type:        "Bug",
			IsCompleted: false,
			CompletedAt: time.Time{}, // No completion date
			Estimate:    2,
		},
	}

	report, err := ThroughputReport(items, "month")
	if err != nil {
		t.Fatalf("ThroughputReport() error = %v", err)
	}

	// Should only count completed items
	// The report should show 1 item and 3.0 points for May 2024
	lines := strings.Split(report, "\n")
	foundDataLine := false
	for _, line := range lines {
		if strings.Contains(line, "2024-05") {
			foundDataLine = true
			if !strings.Contains(line, "1") { // 1 item
				t.Errorf("Report should show 1 completed item, got line: %s", line)
			}
			if !strings.Contains(line, "3.0") { // 3.0 points
				t.Errorf("Report should show 3.0 points, got line: %s", line)
			}
			break
		}
	}

	if !foundDataLine {
		t.Errorf("Report doesn't contain data line for 2024-05")
	}
}

func TestThroughputReport_UnspecifiedType(t *testing.T) {
	// Test items with empty type field
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task with type",
			Type:        "Feature",
			IsCompleted: true,
			CompletedAt: time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC),
			Estimate:    3,
		},
		{
			ID:          "2",
			Name:        "Task without type",
			Type:        "", // Empty type
			IsCompleted: true,
			CompletedAt: time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC),
			Estimate:    2,
		},
	}

	report, err := ThroughputReport(items, "month")
	if err != nil {
		t.Fatalf("ThroughputReport() error = %v", err)
	}

	// Should categorize empty type as "Unspecified"
	if !strings.Contains(report, "Unspecified") {
		t.Errorf("Report doesn't contain 'Unspecified' category for empty types")
	}

	if !strings.Contains(report, "Feature") {
		t.Errorf("Report doesn't contain 'Feature' type")
	}
}

func TestThroughputReport_CalculateAveragePointsPerItem(t *testing.T) {
	// Test average points per item calculation
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Type:        "Feature",
			IsCompleted: true,
			CompletedAt: time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC),
			Estimate:    4,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Type:        "Bug",
			IsCompleted: true,
			CompletedAt: time.Date(2024, 5, 16, 12, 0, 0, 0, time.UTC),
			Estimate:    2,
		},
	}

	report, err := ThroughputReport(items, "month")
	if err != nil {
		t.Fatalf("ThroughputReport() error = %v", err)
	}

	// Should show 2 items, 6.0 total points, 3.0 average points per item
	lines := strings.Split(report, "\n")
	foundDataLine := false
	for _, line := range lines {
		if strings.Contains(line, "2024-05") {
			foundDataLine = true
			if !strings.Contains(line, "2") { // 2 items
				t.Errorf("Report should show 2 items, got line: %s", line)
			}
			if !strings.Contains(line, "6.0") { // 6.0 total points
				t.Errorf("Report should show 6.0 points, got line: %s", line)
			}
			if !strings.Contains(line, "3.0") { // 3.0 average per item
				t.Errorf("Report should show 3.0 average points per item, got line: %s", line)
			}
			break
		}
	}

	if !foundDataLine {
		t.Errorf("Report doesn't contain data line for 2024-05")
	}
}

func TestThroughputReport_ChronologicalSorting(t *testing.T) {
	// Test that periods are sorted chronologically
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Type:        "Feature",
			IsCompleted: true,
			CompletedAt: time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC), // June (later)
			Estimate:    3,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Type:        "Bug",
			IsCompleted: true,
			CompletedAt: time.Date(2024, 4, 15, 12, 0, 0, 0, time.UTC), // April (earlier)
			Estimate:    2,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			Type:        "Task",
			IsCompleted: true,
			CompletedAt: time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC), // May (middle)
			Estimate:    1,
		},
	}

	report, err := ThroughputReport(items, "month")
	if err != nil {
		t.Fatalf("ThroughputReport() error = %v", err)
	}

	// Find positions of each month in the report
	april := strings.Index(report, "2024-04")
	may := strings.Index(report, "2024-05")
	june := strings.Index(report, "2024-06")

	if april == -1 || may == -1 || june == -1 {
		t.Fatalf("Not all months found in report")
	}

	// April should come before May, May should come before June
	if april > may {
		t.Errorf("April should appear before May in chronological order")
	}

	if may > june {
		t.Errorf("May should appear before June in chronological order")
	}
}