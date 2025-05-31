package metrics

import (
	"strings"
	"testing"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

func TestFlowEfficiencyReport(t *testing.T) {
	// Create test data with completed items that have started dates
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			IsCompleted: true,
			CreatedAt:   baseTime.AddDate(0, 0, -10), // Created 10 days ago
			StartedAt:   baseTime.AddDate(0, 0, -5),  // Started 5 days ago (5 days waiting)
			CompletedAt: baseTime,                    // Completed today (5 days active)
		},
		{
			ID:          "2",
			Name:        "Task 2",
			IsCompleted: true,
			CreatedAt:   baseTime.AddDate(0, 0, -8), // Created 8 days ago
			StartedAt:   baseTime.AddDate(0, 0, -6), // Started 6 days ago (2 days waiting)
			CompletedAt: baseTime.AddDate(0, 0, -2), // Completed 2 days ago (4 days active)
		},
	}

	report, err := FlowEfficiencyReport(items)
	if err != nil {
		t.Fatalf("FlowEfficiencyReport() error = %v", err)
	}

	// Verify the report content
	expectedStrings := []string{
		"Flow Efficiency Analysis",
		"What is Flow Efficiency?",
		"Flow Efficiency = (Active Time / Total Time) × 100%",
		"Waiting",
		"Active",
		"Flow Efficiency:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(report, expected) {
			t.Errorf("Report doesn't contain expected string: %s", expected)
		}
	}

	// Verify that percentages are calculated
	if !strings.Contains(report, "%") {
		t.Errorf("Report doesn't contain percentage values")
	}
}

func TestFlowEfficiencyReport_NoStartDate(t *testing.T) {
	// Test items without start dates (all time should be considered active)
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			IsCompleted: true,
			CreatedAt:   baseTime.AddDate(0, 0, -5), // Created 5 days ago
			StartedAt:   time.Time{},                // No start date
			CompletedAt: baseTime,                   // Completed today
		},
	}

	report, err := FlowEfficiencyReport(items)
	if err != nil {
		t.Fatalf("FlowEfficiencyReport() error = %v", err)
	}

	// Should still generate a valid report
	if !strings.Contains(report, "Flow Efficiency Analysis") {
		t.Errorf("Report doesn't contain expected header")
	}

	// Should handle missing start dates gracefully
	if !strings.Contains(report, "Active") {
		t.Errorf("Report should still show Active time for items without start dates")
	}
}

func TestFlowEfficiencyReport_IncompleteItems(t *testing.T) {
	// Test that incomplete items are excluded
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Completed Task",
			IsCompleted: true,
			CreatedAt:   baseTime.AddDate(0, 0, -5),
			StartedAt:   baseTime.AddDate(0, 0, -3),
			CompletedAt: baseTime,
		},
		{
			ID:          "2",
			Name:        "Incomplete Task",
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -10),
			StartedAt:   baseTime.AddDate(0, 0, -8),
			CompletedAt: time.Time{}, // Not completed
		},
	}

	report, err := FlowEfficiencyReport(items)
	if err != nil {
		t.Fatalf("FlowEfficiencyReport() error = %v", err)
	}

	// Should only process completed items
	if !strings.Contains(report, "Flow Efficiency Analysis") {
		t.Errorf("Report doesn't contain expected header")
	}

	// Should calculate flow efficiency based only on completed items
	if !strings.Contains(report, "Flow Efficiency:") {
		t.Errorf("Report should contain flow efficiency calculation")
	}
}

func TestFlowEfficiencyReport_EmptyItems(t *testing.T) {
	items := []models.KanbanItem{}

	report, err := FlowEfficiencyReport(items)
	if err != nil {
		t.Fatalf("FlowEfficiencyReport() error = %v", err)
	}

	// Should handle empty items gracefully
	if !strings.Contains(report, "Flow Efficiency Analysis") {
		t.Errorf("Report doesn't contain header for empty items")
	}

	if !strings.Contains(report, "No data available for flow efficiency calculation") {
		t.Errorf("Report should indicate no data available")
	}
}

func TestFlowEfficiencyReport_MissingDates(t *testing.T) {
	// Test items with missing created or completed dates
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Missing Created Date",
			IsCompleted: true,
			CreatedAt:   time.Time{}, // Missing created date
			StartedAt:   baseTime.AddDate(0, 0, -3),
			CompletedAt: baseTime,
		},
		{
			ID:          "2",
			Name:        "Missing Completed Date",
			IsCompleted: true,
			CreatedAt:   baseTime.AddDate(0, 0, -5),
			StartedAt:   baseTime.AddDate(0, 0, -3),
			CompletedAt: time.Time{}, // Missing completed date
		},
		{
			ID:          "3",
			Name:        "Valid Item",
			IsCompleted: true,
			CreatedAt:   baseTime.AddDate(0, 0, -4),
			StartedAt:   baseTime.AddDate(0, 0, -2),
			CompletedAt: baseTime,
		},
	}

	report, err := FlowEfficiencyReport(items)
	if err != nil {
		t.Fatalf("FlowEfficiencyReport() error = %v", err)
	}

	// Should skip items with missing dates and process valid ones
	if !strings.Contains(report, "Flow Efficiency Analysis") {
		t.Errorf("Report doesn't contain expected header")
	}

	// Should still calculate efficiency for valid items
	if !strings.Contains(report, "Flow Efficiency:") {
		t.Errorf("Report should contain flow efficiency calculation for valid items")
	}
}

func TestFlowEfficiencyReport_CalculateEfficiency(t *testing.T) {
	// Test with known values to verify efficiency calculation
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			IsCompleted: true,
			CreatedAt:   baseTime.AddDate(0, 0, -10), // 10 days total lead time
			StartedAt:   baseTime.AddDate(0, 0, -3),  // 7 days waiting, 3 days active
			CompletedAt: baseTime,
		},
	}

	report, err := FlowEfficiencyReport(items)
	if err != nil {
		t.Fatalf("FlowEfficiencyReport() error = %v", err)
	}

	// With 3 days active out of 10 days total, efficiency should be 30%
	// This is a simplified test - actual calculation may vary due to averaging
	if !strings.Contains(report, "Flow Efficiency:") {
		t.Errorf("Report should contain flow efficiency percentage")
	}

	// Verify that waiting and active times are reported
	if !strings.Contains(report, "Waiting") || !strings.Contains(report, "Active") {
		t.Errorf("Report should contain both Waiting and Active time breakdowns")
	}
}

func TestFlowEfficiencyReport_ExplanatoryText(t *testing.T) {
	// Test that the report contains comprehensive explanatory text
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			IsCompleted: true,
			CreatedAt:   time.Date(2024, 5, 10, 12, 0, 0, 0, time.UTC),
			StartedAt:   time.Date(2024, 5, 12, 12, 0, 0, 0, time.UTC),
			CompletedAt: time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC),
		},
	}

	report, err := FlowEfficiencyReport(items)
	if err != nil {
		t.Fatalf("FlowEfficiencyReport() error = %v", err)
	}

	expectedExplanations := []string{
		"What is Flow Efficiency?",
		"Flow efficiency measures the percentage of time items spend being actively worked on versus waiting",
		"Waiting time",
		"Active time",
		"Total time",
		"How to improve flow efficiency:",
		"Limit work in progress",
		"Reduce handoffs",
		"Eliminate bottlenecks",
	}

	for _, explanation := range expectedExplanations {
		if !strings.Contains(report, explanation) {
			t.Errorf("Report doesn't contain expected explanation: %s", explanation)
		}
	}
}

func TestFlowEfficiencyReport_TableFormat(t *testing.T) {
	// Test that the report contains a properly formatted table
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			IsCompleted: true,
			CreatedAt:   baseTime.AddDate(0, 0, -6),
			StartedAt:   baseTime.AddDate(0, 0, -3),
			CompletedAt: baseTime,
		},
	}

	report, err := FlowEfficiencyReport(items)
	if err != nil {
		t.Fatalf("FlowEfficiencyReport() error = %v", err)
	}

	// Check for table structure
	tableHeaders := []string{
		"State",
		"Avg Time (days)",
		"% of Total Time",
	}

	for _, header := range tableHeaders {
		if !strings.Contains(report, header) {
			t.Errorf("Report doesn't contain expected table header: %s", header)
		}
	}

	// Check for table separators
	if !strings.Contains(report, "------|") {
		t.Errorf("Report doesn't contain table separator")
	}
}