package metrics

import (
	"strings"
	"testing"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

func TestWorkItemAgeReport(t *testing.T) {
	// Create test data with incomplete items in different states
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Old In Progress Task",
			State:       "In Progress",
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -10), // Created 10 days ago
			StartedAt:   baseTime.AddDate(0, 0, -8),  // Started 8 days ago (8 days old)
		},
		{
			ID:          "2",
			Name:        "Recent In Progress Task",
			State:       "In Progress",
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -5), // Created 5 days ago
			StartedAt:   baseTime.AddDate(0, 0, -3), // Started 3 days ago (3 days old)
		},
		{
			ID:          "3",
			Name:        "Waiting Task",
			State:       "To Do",
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -7), // Created 7 days ago (no start date, so 7 days old)
			StartedAt:   time.Time{},                // Not started yet
		},
		{
			ID:          "4",
			Name:        "Completed Task",
			State:       "Done",
			IsCompleted: true, // Should be excluded
			CreatedAt:   baseTime.AddDate(0, 0, -15),
			StartedAt:   baseTime.AddDate(0, 0, -12),
			CompletedAt: baseTime.AddDate(0, 0, -2),
		},
	}

	report, err := WorkItemAgeReport(items, baseTime)
	if err != nil {
		t.Fatalf("WorkItemAgeReport() error = %v", err)
	}

	// Verify the report content
	expectedStrings := []string{
		"Current Work Item Age Analysis",
		"Age of incomplete items by state (in days)",
		"In Progress",
		"To Do",
		"Min:",
		"Max:",
		"Avg:",
		"Median:",
		"Oldest Items:",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(report, expected) {
			t.Errorf("Report doesn't contain expected string: %s", expected)
		}
	}

	// Should contain the task names
	if !strings.Contains(report, "Old In Progress Task") {
		t.Errorf("Report should contain 'Old In Progress Task'")
	}

	if !strings.Contains(report, "Recent In Progress Task") {
		t.Errorf("Report should contain 'Recent In Progress Task'")
	}

	if !strings.Contains(report, "Waiting Task") {
		t.Errorf("Report should contain 'Waiting Task'")
	}

	// Should NOT contain completed tasks
	if strings.Contains(report, "Completed Task") {
		t.Errorf("Report should not contain completed tasks")
	}
}

func TestWorkItemAgeReport_EmptyItems(t *testing.T) {
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	items := []models.KanbanItem{}

	report, err := WorkItemAgeReport(items, baseTime)
	if err != nil {
		t.Fatalf("WorkItemAgeReport() error = %v", err)
	}

	// Should handle empty items gracefully
	if !strings.Contains(report, "Current Work Item Age Analysis") {
		t.Errorf("Report doesn't contain header for empty items")
	}

	if !strings.Contains(report, "Age of incomplete items by state") {
		t.Errorf("Report should contain explanatory text even for empty items")
	}
}

func TestWorkItemAgeReport_OnlyCompletedItems(t *testing.T) {
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Completed Task 1",
			State:       "Done",
			IsCompleted: true,
			CreatedAt:   baseTime.AddDate(0, 0, -10),
			StartedAt:   baseTime.AddDate(0, 0, -8),
			CompletedAt: baseTime.AddDate(0, 0, -2),
		},
		{
			ID:          "2",
			Name:        "Completed Task 2",
			State:       "Done",
			IsCompleted: true,
			CreatedAt:   baseTime.AddDate(0, 0, -5),
			StartedAt:   baseTime.AddDate(0, 0, -3),
			CompletedAt: baseTime.AddDate(0, 0, -1),
		},
	}

	report, err := WorkItemAgeReport(items, baseTime)
	if err != nil {
		t.Fatalf("WorkItemAgeReport() error = %v", err)
	}

	// Should handle case where all items are completed
	if !strings.Contains(report, "Current Work Item Age Analysis") {
		t.Errorf("Report should contain header")
	}

	// Should not contain any task names since all are completed
	if strings.Contains(report, "Completed Task 1") || strings.Contains(report, "Completed Task 2") {
		t.Errorf("Report should not contain completed tasks")
	}
}

func TestWorkItemAgeReport_UnknownState(t *testing.T) {
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task with Empty State",
			State:       "", // Empty state
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -5),
			StartedAt:   baseTime.AddDate(0, 0, -3),
		},
		{
			ID:          "2",
			Name:        "Task with Custom State",
			State:       "Custom State",
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -4),
			StartedAt:   baseTime.AddDate(0, 0, -2),
		},
	}

	report, err := WorkItemAgeReport(items, baseTime)
	if err != nil {
		t.Fatalf("WorkItemAgeReport() error = %v", err)
	}

	// Should categorize empty state as "Unknown"
	if !strings.Contains(report, "Unknown") {
		t.Errorf("Report should contain 'Unknown' category for empty states")
	}

	// Should handle custom states
	if !strings.Contains(report, "Custom State") {
		t.Errorf("Report should contain custom state names")
	}
}

func TestWorkItemAgeReport_AgeCalculation(t *testing.T) {
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task with Start Date",
			State:       "In Progress",
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -10), // 10 days ago
			StartedAt:   baseTime.AddDate(0, 0, -6),  // 6 days ago (should use started date)
		},
		{
			ID:          "2",
			Name:        "Task without Start Date",
			State:       "To Do",
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -8), // 8 days ago (should use created date)
			StartedAt:   time.Time{},                // No start date
		},
	}

	report, err := WorkItemAgeReport(items, baseTime)
	if err != nil {
		t.Fatalf("WorkItemAgeReport() error = %v", err)
	}

	// Should calculate age correctly for both scenarios
	if !strings.Contains(report, "In Progress") {
		t.Errorf("Report should contain 'In Progress' state")
	}

	if !strings.Contains(report, "To Do") {
		t.Errorf("Report should contain 'To Do' state")
	}

	// Should show statistical information
	if !strings.Contains(report, "Min:") || !strings.Contains(report, "Max:") {
		t.Errorf("Report should contain statistical information")
	}
}

func TestWorkItemAgeReport_DefaultAsOfTime(t *testing.T) {
	// Test with zero time (should use current time)
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Recent Task",
			State:       "In Progress",
			IsCompleted: false,
			CreatedAt:   time.Now().AddDate(0, 0, -1), // 1 day ago
			StartedAt:   time.Now().AddDate(0, 0, -1),
		},
	}

	report, err := WorkItemAgeReport(items, time.Time{}) // Zero time
	if err != nil {
		t.Fatalf("WorkItemAgeReport() error = %v", err)
	}

	// Should handle zero time by using current time
	if !strings.Contains(report, "Current Work Item Age Analysis") {
		t.Errorf("Report should contain header when using default time")
	}

	if !strings.Contains(report, "In Progress") {
		t.Errorf("Report should process items when using default time")
	}
}

func TestWorkItemAgeReport_SortingByAge(t *testing.T) {
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Newest Task",
			State:       "In Progress",
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -2), // 2 days ago
			StartedAt:   baseTime.AddDate(0, 0, -1), // 1 day old
		},
		{
			ID:          "2",
			Name:        "Oldest Task",
			State:       "In Progress",
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -10), // 10 days ago
			StartedAt:   baseTime.AddDate(0, 0, -8),  // 8 days old
		},
		{
			ID:          "3",
			Name:        "Middle Task",
			State:       "In Progress",
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -6), // 6 days ago
			StartedAt:   baseTime.AddDate(0, 0, -4), // 4 days old
		},
	}

	report, err := WorkItemAgeReport(items, baseTime)
	if err != nil {
		t.Fatalf("WorkItemAgeReport() error = %v", err)
	}

	// Find positions of each task in the "Oldest Items" section
	oldestPos := strings.Index(report, "Oldest Task")
	middlePos := strings.Index(report, "Middle Task")
	newestPos := strings.Index(report, "Newest Task")

	if oldestPos == -1 || middlePos == -1 || newestPos == -1 {
		t.Fatalf("Not all tasks found in report")
	}

	// Oldest task should appear first in the "Oldest Items" list
	if oldestPos > middlePos {
		t.Errorf("Oldest Task should appear before Middle Task (sorting by age descending)")
	}

	if middlePos > newestPos {
		t.Errorf("Middle Task should appear before Newest Task (sorting by age descending)")
	}
}

func TestWorkItemAgeReport_MultipleStates(t *testing.T) {
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Backlog Task",
			State:       "Backlog",
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -5),
			StartedAt:   time.Time{},
		},
		{
			ID:          "2",
			Name:        "In Progress Task",
			State:       "In Progress",
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -3),
			StartedAt:   baseTime.AddDate(0, 0, -2),
		},
		{
			ID:          "3",
			Name:        "Review Task",
			State:       "In Review",
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -4),
			StartedAt:   baseTime.AddDate(0, 0, -3),
		},
	}

	report, err := WorkItemAgeReport(items, baseTime)
	if err != nil {
		t.Fatalf("WorkItemAgeReport() error = %v", err)
	}

	// Should contain all different states
	states := []string{"Backlog", "In Progress", "In Review"}
	for _, state := range states {
		if !strings.Contains(report, state) {
			t.Errorf("Report should contain state: %s", state)
		}
	}

	// Each state should have its own section with statistics
	for _, state := range states {
		stateSection := strings.Index(report, state)
		if stateSection == -1 {
			continue
		}

		// Look for statistics after the state name
		reportAfterState := report[stateSection:]
		if !strings.Contains(reportAfterState, "Min:") {
			t.Errorf("State %s should have statistical information", state)
		}
	}
}

func TestWorkItemAgeReport_LimitOldestItems(t *testing.T) {
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	// Create more than 5 items to test the limit
	items := []models.KanbanItem{}
	for i := 0; i < 8; i++ {
		items = append(items, models.KanbanItem{
			ID:          string(rune('1' + i)),
			Name:        "Task " + string(rune('A'+i)),
			State:       "In Progress",
			IsCompleted: false,
			CreatedAt:   baseTime.AddDate(0, 0, -(i+1)),
			StartedAt:   baseTime.AddDate(0, 0, -(i+1)),
		})
	}

	report, err := WorkItemAgeReport(items, baseTime)
	if err != nil {
		t.Fatalf("WorkItemAgeReport() error = %v", err)
	}

	// Should limit to 5 oldest items in the "Oldest Items" section
	oldestSection := strings.Index(report, "Oldest Items:")
	if oldestSection == -1 {
		t.Fatalf("Report should contain 'Oldest Items' section")
	}

	// Count how many tasks are listed in the oldest items section
	oldestSectionText := report[oldestSection:]
	nextSectionStart := strings.Index(oldestSectionText, "\n\n")
	if nextSectionStart != -1 {
		oldestSectionText = oldestSectionText[:nextSectionStart]
	}

	taskCount := strings.Count(oldestSectionText, "Task ")
	if taskCount > 5 {
		t.Errorf("Oldest Items section should show at most 5 items, but shows %d", taskCount)
	}
}