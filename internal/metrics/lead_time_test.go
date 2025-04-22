// internal/metrics/lead_time_test.go
package metrics

import (
	"strings"
	"testing"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

func TestLeadTimeReport(t *testing.T) {
	// Create test data with different story point sizes and lead times
	now := time.Now()
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Estimate:    1,
			IsCompleted: true,
			CreatedAt:   now.AddDate(0, 0, -10),
			StartedAt:   now.AddDate(0, 0, -7),
			CompletedAt: now.AddDate(0, 0, -5),
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Estimate:    3,
			IsCompleted: true,
			CreatedAt:   now.AddDate(0, 0, -15),
			StartedAt:   now.AddDate(0, 0, -12),
			CompletedAt: now.AddDate(0, 0, -8),
		},
		{
			ID:          "3",
			Name:        "Task 3",
			Estimate:    5,
			IsCompleted: true,
			CreatedAt:   now.AddDate(0, 0, -20),
			StartedAt:   now.AddDate(0, 0, -18),
			CompletedAt: now.AddDate(0, 0, -10),
		},
	}

	// Generate lead time report
	report, err := LeadTimeReport(items)
	if err != nil {
		t.Fatalf("LeadTimeReport() error = %v", err)
	}

	// Verify the report content
	if !strings.Contains(report, "Lead Time Analysis") {
		t.Errorf("Report doesn't contain expected header")
	}

	// Verify that all three story point sizes are included
	if !strings.Contains(report, "1 ") || !strings.Contains(report, "3 ") || !strings.Contains(report, "5 ") {
		t.Errorf("Report doesn't contain all expected story point sizes")
	}

	// Verify that both lead time and cycle time sections exist
	if !strings.Contains(report, "Lead Time (Creation to Completion)") || 
	   !strings.Contains(report, "Cycle Time (Start to Completion)") {
		t.Errorf("Report doesn't contain both lead time and cycle time sections")
	}
}