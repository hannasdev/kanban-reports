package filtering

import (
	"testing"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
	"github.com/hannasdev/kanban-reports/pkg/types"
)

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
			Labels:      []string{"ad-hoc-request"},
		},
		{
			ID:          "4",
			Name:        "Incomplete Item",
			IsCompleted: false,
		},
	}
	
	tests := []struct {
		name       string
		startDate  time.Time
		endDate    time.Time
		adHocFilter types.AdHocFilterType
		expected   int
	}{
		{
			name:       "No date range (all completed items)",
			startDate:  time.Time{},
			endDate:    time.Time{},
			adHocFilter: types.AdHocFilterInclude,
			expected:   3,
		},
		{
			name:       "Only items completed today",
			startDate:  baseTime.AddDate(0, 0, -1),
			endDate:    baseTime.AddDate(0, 0, 1),
			adHocFilter: types.AdHocFilterInclude,
			expected:   1,
		},
		{
			name:       "Items completed in last week",
			startDate:  baseTime.AddDate(0, 0, -7),
			endDate:    baseTime,
			adHocFilter: types.AdHocFilterInclude,
			expected:   2,
		},
		{
			name:       "No items in range",
			startDate:  baseTime.AddDate(0, 0, 1),
			endDate:    baseTime.AddDate(0, 0, 2),
			adHocFilter: types.AdHocFilterInclude,
			expected:   0,
		},
		{
			name:       "Exclude ad-hoc requests",
			startDate:  time.Time{},
			endDate:    time.Time{},
			adHocFilter: types.AdHocFilterExclude,
			expected:   2, // 3 completed items minus 1 ad-hoc
		},
		{
			name:       "Only ad-hoc requests",
			startDate:  time.Time{},
			endDate:    time.Time{},
			adHocFilter: types.AdHocFilterOnly,
			expected:   1, // Only item 3 has the ad-hoc label
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := FilterItemsByDateRange(items, tt.startDate, tt.endDate, models.FilterFieldCompletedAt, tt.adHocFilter)
			if len(filtered) != tt.expected {
				t.Errorf("FilterItemsByDateRange() returned %d items, expected %d", len(filtered), tt.expected)
			}
		})
	}
}

func TestIsAdHocRequest(t *testing.T) {
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
			result := IsAdHocRequest(tt.item)
			if result != tt.expected {
				t.Errorf("IsAdHocRequest() = %v, expected %v", result, tt.expected)
			}
		})
	}
}