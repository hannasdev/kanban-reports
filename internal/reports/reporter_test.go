package reports

import (
	"strings"
	"testing"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

func TestNewReporter(t *testing.T) {
	items := []models.KanbanItem{
		{ID: "1", Name: "Test Item"},
	}
	
	reporter := NewReporter(items)
	
	if reporter == nil {
		t.Fatal("NewReporter() returned nil")
	}
	
	if len(reporter.items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(reporter.items))
	}
	
	if reporter.adHocFilter != AdHocFilterInclude {
		t.Errorf("Expected default adHocFilter to be AdHocFilterInclude, got %v", reporter.adHocFilter)
	}
}

func TestWithAdHocFilter(t *testing.T) {
	reporter := NewReporter(nil)
	
	// Test each filter type
	filters := []AdHocFilterType{
		AdHocFilterInclude,
		AdHocFilterExclude,
		AdHocFilterOnly,
	}
	
	for _, filter := range filters {
		result := reporter.WithAdHocFilter(filter)
		
		// Should return the same instance (fluent interface)
		if result != reporter {
			t.Errorf("WithAdHocFilter() didn't return the same reporter instance")
		}
		
		if reporter.adHocFilter != filter {
			t.Errorf("Expected adHocFilter to be %v, got %v", filter, reporter.adHocFilter)
		}
	}
}

func TestFilterItemsByDateRange(t *testing.T) {
	// Create base time
	baseTime := time.Date(2024, 5, 15, 12, 0, 0, 0, time.UTC)
	
	// Create test items with different completion dates
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Item 1",
			IsCompleted: true,
			CompletedAt: baseTime.AddDate(0, 0, -10), // 10 days ago
		},
		{
			ID:          "2",
			Name:        "Item 2",
			IsCompleted: true,
			CompletedAt: baseTime.AddDate(0, 0, -5), // 5 days ago
		},
		{
			ID:          "3",
			Name:        "Item 3",
			IsCompleted: true,
			CompletedAt: baseTime, // today
		},
		{
			ID:          "4",
			Name:        "Incomplete Item",
			IsCompleted: false,
		},
	}
	
	reporter := NewReporter(items)
	
	tests := []struct {
		name      string
		startDate time.Time
		endDate   time.Time
		expected  int
	}{
		{
			name:      "No date range (all completed items)",
			startDate: time.Time{},
			endDate:   time.Time{},
			expected:  3,
		},
		{
			name:      "Only items completed today",
			startDate: baseTime.AddDate(0, 0, -1),
			endDate:   baseTime.AddDate(0, 0, 1),
			expected:  1,
		},
		{
			name:      "Items completed in last week",
			startDate: baseTime.AddDate(0, 0, -7),
			endDate:   baseTime,
			expected:  2,
		},
		{
			name:      "No items in range",
			startDate: baseTime.AddDate(0, 0, 1),
			endDate:   baseTime.AddDate(0, 0, 2),
			expected:  0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := reporter.filterItemsByDateRange(tt.startDate, tt.endDate)
			if len(filtered) != tt.expected {
				t.Errorf("filterItemsByDateRange() returned %d items, expected %d", len(filtered), tt.expected)
			}
		})
	}
}

func TestIsAdHocRequest(t *testing.T) {
	reporter := NewReporter(nil)
	
	tests := []struct {
		name     string
		item     models.KanbanItem
		expected bool
	}{
		{
			name: "Item with ad-hoc-request label",
			item: models.KanbanItem{
				ID:     "1",
				Labels: []string{"ad-hoc-request", "priority"},
			},
			expected: true,
		},
		{
			name: "Item with AD-HOC-REQUEST label (case insensitive)",
			item: models.KanbanItem{
				ID:     "2",
				Labels: []string{"AD-HOC-REQUEST"},
			},
			expected: true,
		},
		{
			name: "Item without ad-hoc-request label",
			item: models.KanbanItem{
				ID:     "3",
				Labels: []string{"feature", "priority"},
			},
			expected: false,
		},
		{
			name: "Item with no labels",
			item: models.KanbanItem{
				ID:     "4",
				Labels: []string{},
			},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reporter.isAdHocRequest(tt.item)
			if result != tt.expected {
				t.Errorf("isAdHocRequest() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestGenerateReport(t *testing.T) {
	// Create test items
	now := time.Now()
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Owners:      []string{"john@example.com"},
			IsCompleted: true,
			CompletedAt: now,
			Estimate:    3,
			Epic:        "Epic 1",
			Team:        "Team A",
			ProductArea: "Backend",
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Owners:      []string{"jane@example.com"},
			IsCompleted: true,
			CompletedAt: now,
			Estimate:    2,
			Epic:        "Epic 1",
			Team:        "Team A",
			ProductArea: "Frontend",
		},
	}
	
	reporter := NewReporter(items)
	
	tests := []struct {
		name       string
		reportType ReportType
		contains   []string
	}{
		{
			name:       "Contributor Report",
			reportType: ReportTypeContributor,
			contains:   []string{"Story Points by Contributor", "john@example.com", "jane@example.com"},
		},
		{
			name:       "Epic Report",
			reportType: ReportTypeEpic,
			contains:   []string{"Story Points by Epic", "Epic 1"},
		},
		{
			name:       "Product Area Report",
			reportType: ReportTypeProductArea,
			contains:   []string{"Story Points by Product Area", "Backend", "Frontend"},
		},
		{
			name:       "Team Report",
			reportType: ReportTypeTeam,
			contains:   []string{"Story Points by Team", "Team A"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := reporter.GenerateReport(tt.reportType, time.Time{}, time.Time{})
			if err != nil {
				t.Fatalf("GenerateReport() error = %v", err)
			}
			
			for _, str := range tt.contains {
				if !strings.Contains(report, str) {
					t.Errorf("Report doesn't contain expected string: %s", str)
				}
			}
		})
	}
	
	// Test invalid report type
	_, err := reporter.GenerateReport("invalid-type", time.Time{}, time.Time{})
	if err == nil {
		t.Errorf("GenerateReport() with invalid type should return error")
	}
}

