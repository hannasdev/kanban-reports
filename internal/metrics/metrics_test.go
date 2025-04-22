// internal/metrics/metrics_test.go
package metrics

import (
	"strings"
	"testing"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
	"github.com/hannasdev/kanban-reports/internal/reports"
)

func TestNewGenerator(t *testing.T) {
	items := []models.KanbanItem{
		{ID: "1", Name: "Test Item"},
	}
	
	generator := NewGenerator(items)
	
	if generator == nil {
		t.Fatal("NewGenerator() returned nil")
	}
	
	if len(generator.items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(generator.items))
	}
	
	if generator.adHocFilter != reports.AdHocFilterInclude {
		t.Errorf("Expected default adHocFilter to be AdHocFilterInclude, got %v", generator.adHocFilter)
	}
}

func TestWithAdHocFilter(t *testing.T) {
	generator := NewGenerator(nil)
	
	// Test each filter type
	filters := []reports.AdHocFilterType{
		reports.AdHocFilterInclude,
		reports.AdHocFilterExclude,
		reports.AdHocFilterOnly,
	}
	
	for _, filter := range filters {
		result := generator.WithAdHocFilter(filter)
		
		// Should return the same instance (fluent interface)
		if result != generator {
			t.Errorf("WithAdHocFilter() didn't return the same generator instance")
		}
		
		if generator.adHocFilter != filter {
			t.Errorf("Expected adHocFilter to be %v, got %v", filter, generator.adHocFilter)
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
	
	generator := NewGenerator(items)
	
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
			filtered := generator.filterItemsByDateRange(tt.startDate, tt.endDate)
			if len(filtered) != tt.expected {
				t.Errorf("filterItemsByDateRange() returned %d items, expected %d", len(filtered), tt.expected)
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	// Create test items
	now := time.Now()
	baseTime := now.AddDate(0, 0, -30) // 30 days ago
	
	items := []models.KanbanItem{
		{
			ID:          "1",
			Name:        "Task 1",
			Estimate:    3,
			IsCompleted: true,
			CreatedAt:   baseTime.AddDate(0, 0, -10),
			StartedAt:   baseTime.AddDate(0, 0, -7),
			CompletedAt: baseTime.AddDate(0, 0, -5),
			Type:        "Feature",
		},
		{
			ID:          "2",
			Name:        "Task 2",
			Estimate:    1,
			IsCompleted: true,
			CreatedAt:   baseTime.AddDate(0, 0, -8),
			StartedAt:   baseTime.AddDate(0, 0, -5),
			CompletedAt: baseTime.AddDate(0, 0, -3),
			Type:        "Bug",
		},
		{
			ID:          "3",
			Name:        "Task 3",
			Estimate:    5,
			IsCompleted: true,
			CreatedAt:   baseTime.AddDate(0, 0, -15),
			StartedAt:   baseTime.AddDate(0, 0, -12),
			CompletedAt: baseTime.AddDate(0, 0, -2),
			Type:        "Feature",
		},
	}
	
	generator := NewGenerator(items)
	
	tests := []struct {
		name        string
		metricsType MetricsType
		periodType  PeriodType
		contains    []string
	}{
		{
			name:        "Lead Time Report",
			metricsType: MetricsTypeLeadTime,
			periodType:  PeriodTypeMonth,
			contains:    []string{"Lead Time Analysis", "Creation to Completion", "Start to Completion"},
		},
		{
			name:        "Throughput Report",
			metricsType: MetricsTypeThroughput,
			periodType:  PeriodTypeMonth,
			contains:    []string{"Throughput Analysis", "Items Completed", "Story Points"},
		},
		{
			name:        "Flow Efficiency Report",
			metricsType: MetricsTypeFlow,
			periodType:  PeriodTypeMonth,
			contains:    []string{"Flow Efficiency Analysis", "Waiting", "Active"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report, err := generator.Generate(tt.metricsType, tt.periodType, time.Time{}, time.Time{})
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			
			for _, str := range tt.contains {
				if !strings.Contains(report, str) {
					t.Errorf("Report doesn't contain expected string: %s", str)
				}
			}
		})
	}
}

func TestAddDateRangeInfo(t *testing.T) {
	generator := NewGenerator(nil)
	
	// Sample metrics content
	metricsContent := "# Lead Time Analysis by Story Point Size (in days)\n\n## Lead Time (Creation to Completion)\n"
	
	tests := []struct {
		name       string
		metricsType MetricsType
		periodType PeriodType
		startDate  time.Time
		endDate    time.Time
		adHocFilter reports.AdHocFilterType
		expected   []string
	}{
		{
			name:       "Full date range",
			metricsType: MetricsTypeLeadTime,
			periodType: PeriodTypeMonth,
			startDate:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC),
			adHocFilter: reports.AdHocFilterInclude,
			expected:   []string{
				"Metrics Type: lead-time",
				"Period Type: month",
				"Date Range: 2024-05-01 to 2024-05-31",
				"# Lead Time Analysis",
			},
		},
		{
			name:       "Only start date",
			metricsType: MetricsTypeThroughput,
			periodType: PeriodTypeWeek,
			startDate:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Time{},
			adHocFilter: reports.AdHocFilterInclude,
			expected:   []string{
				"Metrics Type: throughput",
				"Period Type: week",
				"From: 2024-05-01",
				"# Lead Time Analysis",
			},
		},
		{
			name:       "Only end date",
			metricsType: MetricsTypeFlow,
			periodType: PeriodTypeMonth,
			startDate:  time.Time{},
			endDate:    time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC),
			adHocFilter: reports.AdHocFilterInclude,
			expected:   []string{
				"Metrics Type: flow",
				"Period Type: month",
				"To: 2024-05-31",
				"# Lead Time Analysis",
			},
		},
		{
			name:       "No date range",
			metricsType: MetricsTypeEstimation,
			periodType: PeriodTypeMonth,
			startDate:  time.Time{},
			endDate:    time.Time{},
			adHocFilter: reports.AdHocFilterInclude,
			expected:   []string{
				"Metrics Type: estimation",
				"Period Type: month",
				"Date Range: All Time",
				"# Lead Time Analysis",
			},
		},
		{
			name:       "With ad-hoc exclude filter",
			metricsType: MetricsTypeLeadTime,
			periodType: PeriodTypeMonth,
			startDate:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC),
			adHocFilter: reports.AdHocFilterExclude,
			expected:   []string{
				"Metrics Type: lead-time",
				"Period Type: month",
				"Date Range: 2024-05-01 to 2024-05-31",
				"Filter: Excluding ad-hoc requests",
				"# Lead Time Analysis",
			},
		},
		{
			name:       "With ad-hoc only filter",
			metricsType: MetricsTypeLeadTime,
			periodType: PeriodTypeMonth,
			startDate:  time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC),
			endDate:    time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC),
			adHocFilter: reports.AdHocFilterOnly,
			expected:   []string{
				"Metrics Type: lead-time",
				"Period Type: month",
				"Date Range: 2024-05-01 to 2024-05-31",
				"Filter: Only ad-hoc requests",
				"# Lead Time Analysis",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generator.adHocFilter = tt.adHocFilter
			result := generator.addDateRangeInfo(metricsContent, tt.metricsType, tt.periodType, tt.startDate, tt.endDate)
			
			for _, expectedStr := range tt.expected {
				if !strings.Contains(result, expectedStr) {
					t.Errorf("Expected metrics to contain: %s\nGot: %s", expectedStr, result)
				}
			}
		})
	}
}