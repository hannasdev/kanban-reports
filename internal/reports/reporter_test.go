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
			contains:   []string{"Report Type: contributor", "Date Range:", "Story Points by Contributor", "john@example.com", "jane@example.com"},
		},
		{
			name:       "Epic Report",
			reportType: ReportTypeEpic,
			contains:   []string{"Report Type: epic", "Date Range:", "Story Points by Epic", "Epic 1"},
		},
		{
			name:       "Product Area Report",
			reportType: ReportTypeProductArea,
			contains:   []string{"Report Type: product-area", "Date Range:", "Story Points by Product Area", "Backend", "Frontend"},
		},
		{
			name:       "Team Report",
			reportType: ReportTypeTeam,
			contains:   []string{"Report Type: team", "Date Range:", "Story Points by Team", "Team A"},
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

func TestAddDateRangeInfo(t *testing.T) {
	reporter := NewReporter(nil)
	
	// Sample report content
	reportContent := "Story Points by Contributor:\n\njohn@example.com          5.0 points  2 items\n"
	
	tests := []struct {
		name       string
		reportType ReportType
		startDate  time.Time
		endDate    time.Time
		adHocFilter AdHocFilterType
		expected   []string
	}{
		{
			name:       "Full date range",
			reportType: ReportTypeContributor,
			startDate:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC),
			adHocFilter: AdHocFilterInclude,
			expected:   []string{
				"Report Type: contributor",
				"Date Range: 2024-05-01 to 2024-05-31",
				"Story Points by Contributor:",
			},
		},
		{
			name:       "Only start date",
			reportType: ReportTypeEpic,
			startDate:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Time{},
			adHocFilter: AdHocFilterInclude,
			expected:   []string{
				"Report Type: epic",
				"From: 2024-05-01",
				"Story Points by Contributor:",
			},
		},
		{
			name:       "Only end date",
			reportType: ReportTypeTeam,
			startDate:  time.Time{},
			endDate:    time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC),
			adHocFilter: AdHocFilterInclude,
			expected:   []string{
				"Report Type: team",
				"To: 2024-05-31",
				"Story Points by Contributor:",
			},
		},
		{
			name:       "No date range",
			reportType: ReportTypeProductArea,
			startDate:  time.Time{},
			endDate:    time.Time{},
			adHocFilter: AdHocFilterInclude,
			expected:   []string{
				"Report Type: product-area",
				"Date Range: All Time",
				"Story Points by Contributor:",
			},
		},
		{
			name:       "With ad-hoc exclude filter",
			reportType: ReportTypeContributor,
			startDate:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC),
			adHocFilter: AdHocFilterExclude,
			expected:   []string{
				"Report Type: contributor",
				"Date Range: 2024-05-01 to 2024-05-31",
				"Filter: Excluding ad-hoc requests",
				"Story Points by Contributor:",
			},
		},
		{
			name:       "With ad-hoc only filter",
			reportType: ReportTypeContributor,
			startDate:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC),
			adHocFilter: AdHocFilterOnly,
			expected:   []string{
				"Report Type: contributor",
				"Date Range: 2024-05-01 to 2024-05-31",
				"Filter: Only ad-hoc requests",
				"Story Points by Contributor:",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter.adHocFilter = tt.adHocFilter
			result := reporter.addDateRangeInfo(reportContent, tt.reportType, tt.startDate, tt.endDate)
			
			for _, expectedStr := range tt.expected {
				if !strings.Contains(result, expectedStr) {
					t.Errorf("Expected report to contain: %s\nGot: %s", expectedStr, result)
				}
			}
		})
	}
}

func TestGenerateReportWithDateRange(t *testing.T) {
	// Create test items
	now := time.Now()
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Owners:      []string{"john@example.com"},
			IsCompleted: true,
			CompletedAt: now.AddDate(0, 0, -5), // 5 days ago
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
			CompletedAt: now.AddDate(0, 0, -10), // 10 days ago
			Estimate:    2,
			Epic:        "Epic 1",
			Team:        "Team A",
			ProductArea: "Frontend",
		},
	}
	
	reporter := NewReporter(items)
	
	// Test with date range
	startDate := now.AddDate(0, 0, -7) // 7 days ago
	endDate := now
	
	report, err := reporter.GenerateReport(ReportTypeContributor, startDate, endDate)
	if err != nil {
		t.Fatalf("GenerateReport() error = %v", err)
	}
	
	// Check that the report contains date range information
	expectedStrings := []string{
		"Report Type: contributor",
		"Date Range:",
		"john@example.com", // Should only include items from last 7 days
	}
	
	for _, expected := range expectedStrings {
		if !strings.Contains(report, expected) {
			t.Errorf("Report doesn't contain expected string: %s", expected)
		}
	}
	
	// Check that items outside date range are excluded
	if strings.Contains(report, "jane@example.com") {
		t.Errorf("Report includes items outside the date range")
	}
}