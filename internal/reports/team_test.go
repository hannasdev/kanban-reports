package reports

import (
	"strings"
	"testing"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

func TestGenerateTeamReport(t *testing.T) {
	// Create test data
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Team:        "Team Alpha",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    3,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Team:        "Team Beta",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    2,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			Team:        "Team Alpha",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    4,
		},
		{
			ID:          "4",
			Name:        "Task 4",
			Team:        "", // No team
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    1,
		},
	}

	// Create reporter and generate report
	reporter := NewReporter(items)
	report, err := reporter.generateTeamReport(items)
	if err != nil {
		t.Fatalf("generateTeamReport() error = %v", err)
	}

	// Verify the report content
	if !strings.Contains(report, "Story Points by Team") {
		t.Errorf("Report doesn't contain expected header")
	}

	// Team Alpha should have 3 + 4 = 7 points
	if !strings.Contains(report, "Team Alpha") {
		t.Errorf("Report doesn't contain Team Alpha")
	}

	// Team Beta should have 2 points
	if !strings.Contains(report, "Team Beta") {
		t.Errorf("Report doesn't contain Team Beta")
	}

	// Item without team should be categorized as "No Team"
	if !strings.Contains(report, "No Team") {
		t.Errorf("Report doesn't contain 'No Team' category")
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

func TestGenerateTeamReport_EmptyItems(t *testing.T) {
	// Test with empty items slice
	items := []models.KanbanItem{}

	reporter := NewReporter(items)
	report, err := reporter.generateTeamReport(items)
	if err != nil {
		t.Fatalf("generateTeamReport() error = %v", err)
	}

	// Should still have header and total
	if !strings.Contains(report, "Story Points by Team") {
		t.Errorf("Report doesn't contain expected header")
	}

	if !strings.Contains(report, "Total: 0.0 points across 0 items") {
		t.Errorf("Report doesn't contain correct totals for empty items")
	}
}

func TestGenerateTeamReport_SortingByPoints(t *testing.T) {
	// Create test data with different point values to test sorting
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Team:        "Team Small",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    2,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Team:        "Team Large",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    10,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			Team:        "Team Medium",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    5,
		},
	}

	reporter := NewReporter(items)
	report, err := reporter.generateTeamReport(items)
	if err != nil {
		t.Fatalf("generateTeamReport() error = %v", err)
	}

	// Find positions of each team in the report
	largePos := strings.Index(report, "Team Large")
	mediumPos := strings.Index(report, "Team Medium")
	smallPos := strings.Index(report, "Team Small")

	// Team Large (10 points) should come first
	// Team Medium (5 points) should come second
	// Team Small (2 points) should come last
	if largePos == -1 || mediumPos == -1 || smallPos == -1 {
		t.Fatalf("Not all teams found in report")
	}

	if largePos > mediumPos {
		t.Errorf("Team Large should appear before Team Medium (sorting by points descending)")
	}

	if mediumPos > smallPos {
		t.Errorf("Team Medium should appear before Team Small (sorting by points descending)")
	}
}

func TestGenerateTeamReport_MultipleItemsSameTeam(t *testing.T) {
	// Test aggregation of multiple items for the same team
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Team:        "Team Gamma",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    1.5,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Team:        "Team Gamma",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    2.5,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			Team:        "Team Gamma",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    3.0,
		},
	}

	reporter := NewReporter(items)
	report, err := reporter.generateTeamReport(items)
	if err != nil {
		t.Fatalf("generateTeamReport() error = %v", err)
	}

	// Should aggregate to 7.0 points and 3 items
	if !strings.Contains(report, "7.0 points") {
		t.Errorf("Report doesn't contain correct aggregated points")
	}

	if !strings.Contains(report, "3 items") {
		t.Errorf("Report doesn't contain correct item count")
	}
}

func TestGenerateTeamReport_SpecialCharactersInTeamNames(t *testing.T) {
	// Test with team names that have special characters or spaces
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Team:        "Team Alpha-1",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    5,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Team:        "Team Beta & Gamma",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    3,
		},
		{
			ID:          "3",
			Name:        "Task 3",
			Team:        "Team.Delta",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    2,
		},
	}

	reporter := NewReporter(items)
	report, err := reporter.generateTeamReport(items)
	if err != nil {
		t.Fatalf("generateTeamReport() error = %v", err)
	}

	// Should handle special characters correctly
	if !strings.Contains(report, "Team Alpha-1") {
		t.Errorf("Report doesn't contain team name with hyphen")
	}

	if !strings.Contains(report, "Team Beta & Gamma") {
		t.Errorf("Report doesn't contain team name with ampersand")
	}

	if !strings.Contains(report, "Team.Delta") {
		t.Errorf("Report doesn't contain team name with dot")
	}
}

func TestGenerateTeamReport_ZeroEstimates(t *testing.T) {
	// Test with items that have zero estimates
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Team:        "Team Zero",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    0,
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Team:        "Team Zero",
			IsCompleted: true,
			CompletedAt: time.Now(),
			Estimate:    5,
		},
	}

	reporter := NewReporter(items)
	report, err := reporter.generateTeamReport(items)
	if err != nil {
		t.Fatalf("generateTeamReport() error = %v", err)
	}

	// Should handle zero estimates correctly
	if !strings.Contains(report, "5.0 points") {
		t.Errorf("Report doesn't contain correct points total")
	}

	if !strings.Contains(report, "2 items") {
		t.Errorf("Report doesn't contain correct item count")
	}
}