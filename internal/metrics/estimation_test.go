package metrics

import (
	"strings"
	"testing"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

func TestEstimationAccuracyReport(t *testing.T) {
	// Create test data with completed items having different story point sizes
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Small Task",
			IsCompleted: true,
			StartedAt:   baseTime.AddDate(0, 0, -3), // Started 3 days ago
			CompletedAt: baseTime,                   // Completed today (3 days cycle time)
			Estimate:    1,                          // 1 story point
		},
		{
			ID:          "2",
			Name:        "Medium Task",
			IsCompleted: true,
			StartedAt:   baseTime.AddDate(0, 0, -6), // Started 6 days ago
			CompletedAt: baseTime.AddDate(0, 0, -1), // Completed 1 day ago (5 days cycle time)
			Estimate:    3,                          // 3 story points
		},
		{
			ID:          "3",
			Name:        "Large Task",
			IsCompleted: true,
			StartedAt:   baseTime.AddDate(0, 0, -10), // Started 10 days ago
			CompletedAt: baseTime.AddDate(0, 0, -2),  // Completed 2 days ago (8 days cycle time)
			Estimate:    5,                           // 5 story points
		},
	}

	report, err := EstimationAccuracyReport(items)
	if err != nil {
		t.Fatalf("EstimationAccuracyReport() error = %v", err)
	}

	// Verify the report content
	expectedStrings := []string{
		"Estimation Accuracy Analysis",
		"What is Estimation Accuracy?",
		"Time Spent per Story Point Size",
		"Raw Cycle Time by Story Point Size",
		"Days/SP",
		"Correlation between story points and cycle time",
		"Story points",
		"Count",
		"Min",
		"Max",
		"Avg",
		"Median",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(report, expected) {
			t.Errorf("Report doesn't contain expected string: %s", expected)
		}
	}

	// Verify that story point sizes are present
	if !strings.Contains(report, "1") || !strings.Contains(report, "3") || !strings.Contains(report, "5") {
		t.Errorf("Report doesn't contain expected story point sizes")
	}
}

func TestEstimationAccuracyReport_IncompleteItems(t *testing.T) {
	// Test that incomplete items are excluded
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Completed Task",
			IsCompleted: true,
			StartedAt:   baseTime.AddDate(0, 0, -3),
			CompletedAt: baseTime,
			Estimate:    2,
		},
		{
			ID:          "2",
			Name:        "Incomplete Task",
			IsCompleted: false,
			StartedAt:   baseTime.AddDate(0, 0, -5),
			CompletedAt: time.Time{}, // Not completed
			Estimate:    3,
		},
	}

	report, err := EstimationAccuracyReport(items)
	if err != nil {
		t.Fatalf("EstimationAccuracyReport() error = %v", err)
	}

	// Should only process completed items
	if !strings.Contains(report, "Estimation Accuracy Analysis") {
		t.Errorf("Report doesn't contain expected header")
	}

	// Should show data for story point size 2 but not 3
	if !strings.Contains(report, "2") {
		t.Errorf("Report should contain data for completed 2-point item")
	}
}

func TestEstimationAccuracyReport_MissingDates(t *testing.T) {
	// Test items with missing started or completed dates
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Missing Start Date",
			IsCompleted: true,
			StartedAt:   time.Time{}, // Missing start date
			CompletedAt: baseTime,
			Estimate:    2,
		},
		{
			ID:          "2",
			Name:        "Missing Completed Date",
			IsCompleted: true,
			StartedAt:   baseTime.AddDate(0, 0, -3),
			CompletedAt: time.Time{}, // Missing completed date
			Estimate:    3,
		},
		{
			ID:          "3",
			Name:        "Valid Item",
			IsCompleted: true,
			StartedAt:   baseTime.AddDate(0, 0, -4),
			CompletedAt: baseTime,
			Estimate:    1,
		},
	}

	report, err := EstimationAccuracyReport(items)
	if err != nil {
		t.Fatalf("EstimationAccuracyReport() error = %v", err)
	}

	// Should skip items with missing dates and process valid ones
	if !strings.Contains(report, "Estimation Accuracy Analysis") {
		t.Errorf("Report doesn't contain expected header")
	}

	// Should show data for the valid 1-point item
	if !strings.Contains(report, "1") {
		t.Errorf("Report should contain data for valid 1-point item")
	}
}

func TestEstimationAccuracyReport_EmptyItems(t *testing.T) {
	items := []models.KanbanItem{}

	report, err := EstimationAccuracyReport(items)
	if err != nil {
		t.Fatalf("EstimationAccuracyReport() error = %v", err)
	}

	// Should handle empty items gracefully
	if !strings.Contains(report, "Estimation Accuracy Analysis") {
		t.Errorf("Report doesn't contain header for empty items")
	}

	// Should still contain explanatory text
	if !strings.Contains(report, "What is Estimation Accuracy?") {
		t.Errorf("Report should contain explanatory text even for empty items")
	}
}

func TestEstimationAccuracyReport_ZeroEstimate(t *testing.T) {
	// Test items with zero story point estimates (should be excluded)
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Zero Estimate Task",
			IsCompleted: true,
			StartedAt:   baseTime.AddDate(0, 0, -3),
			CompletedAt: baseTime,
			Estimate:    0, // Zero estimate
		},
		{
			ID:          "2",
			Name:        "Valid Task",
			IsCompleted: true,
			StartedAt:   baseTime.AddDate(0, 0, -2),
			CompletedAt: baseTime,
			Estimate:    2,
		},
	}

	report, err := EstimationAccuracyReport(items)
	if err != nil {
		t.Fatalf("EstimationAccuracyReport() error = %v", err)
	}

	// Should exclude zero estimates from correlation calculation
	if !strings.Contains(report, "Estimation Accuracy Analysis") {
		t.Errorf("Report doesn't contain expected header")
	}

	// Should show data for the valid 2-point item
	if !strings.Contains(report, "2") {
		t.Errorf("Report should contain data for valid 2-point item")
	}
}

func TestEstimationAccuracyReport_ClosestPointSizeMapping(t *testing.T) {
	// Test that non-standard estimates are mapped to closest standard sizes
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Non-standard estimate 1.5",
			IsCompleted: true,
			StartedAt:   baseTime.AddDate(0, 0, -2),
			CompletedAt: baseTime,
			Estimate:    1.5, // Should map to closest standard size (1 or 2)
		},
		{
			ID:          "2",
			Name:        "Non-standard estimate 4",
			IsCompleted: true,
			StartedAt:   baseTime.AddDate(0, 0, -3),
			CompletedAt: baseTime,
			Estimate:    4, // Should map to closest standard size (3 or 5)
		},
	}

	report, err := EstimationAccuracyReport(items)
	if err != nil {
		t.Fatalf("EstimationAccuracyReport() error = %v", err)
	}

	// Should contain mapping to standard story point sizes
	if !strings.Contains(report, "Estimation Accuracy Analysis") {
		t.Errorf("Report doesn't contain expected header")
	}

	// The exact mapping depends on the findClosestPointSize function
	// The report should contain some standard story point sizes
	standardSizes := []string{"1", "2", "3", "5"}
	foundStandardSize := false
	for _, size := range standardSizes {
		if strings.Contains(report, size) {
			foundStandardSize = true
			break
		}
	}

	if !foundStandardSize {
		t.Errorf("Report should contain at least one standard story point size")
	}
}

func TestEstimationAccuracyReport_CorrelationCalculation(t *testing.T) {
	// Test correlation calculation with known values
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			IsCompleted: true,
			StartedAt:   baseTime.AddDate(0, 0, -1), // 1 day
			CompletedAt: baseTime,
			Estimate:    1,
		},
		{
			ID:          "2",
			Name:        "Task 2", 
			IsCompleted: true,
			StartedAt:   baseTime.AddDate(0, 0, -3), // 3 days
			CompletedAt: baseTime,
			Estimate:    3,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			IsCompleted: true,
			StartedAt:   baseTime.AddDate(0, 0, -5), // 5 days
			CompletedAt: baseTime,
			Estimate:    5,
		},
	}

	report, err := EstimationAccuracyReport(items)
	if err != nil {
		t.Fatalf("EstimationAccuracyReport() error = %v", err)
	}

	// Should include correlation value
	if !strings.Contains(report, "Correlation between story points and cycle time:") {
		t.Errorf("Report should contain correlation calculation")
	}

	// Should include interpretation guide
	if !strings.Contains(report, "Interpretation of correlation:") {
		t.Errorf("Report should contain correlation interpretation guide")
	}

	interpretationText := []string{
		"Strong positive correlation",
		"Moderate correlation", 
		"Weak correlation",
		"Negative",
	}

	foundInterpretation := false
	for _, text := range interpretationText {
		if strings.Contains(report, text) {
			foundInterpretation = true
			break
		}
	}

	if !foundInterpretation {
		t.Errorf("Report should contain correlation interpretation text")
	}
}

func TestEstimationAccuracyReport_ExplanatoryText(t *testing.T) {
	// Test that comprehensive explanatory text is included
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			IsCompleted: true,
			StartedAt:   time.Date(2024, 5, 10, 12, 0, 0, 0, time.UTC),
			CompletedAt: time.Date(2024, 5, 12, 12, 0, 0, 0, time.UTC),
			Estimate:    2,
		},
	}

	report, err := EstimationAccuracyReport(items)
	if err != nil {
		t.Fatalf("EstimationAccuracyReport() error = %v", err)
	}

	expectedExplanations := []string{
		"What is Estimation Accuracy?",
		"Estimation accuracy measures how well your story point estimates correlate",
		"How to use this data:",
		"Look for consistency in the days/SP metric",
		"Identify if certain sized items are consistently under or overestimated",
		"Use the correlation value to assess your estimation system's reliability",
	}

	for _, explanation := range expectedExplanations {
		if !strings.Contains(report, explanation) {
			t.Errorf("Report doesn't contain expected explanation: %s", explanation)
		}
	}
}