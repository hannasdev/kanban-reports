package filtering

import (
	"strings"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
	"github.com/hannasdev/kanban-reports/pkg/types" // Use types instead of reports
)

// IsAdHocRequest checks if an item is an ad-hoc request (has "ad-hoc-request" label)
func IsAdHocRequest(item models.KanbanItem) bool {
	for _, label := range item.Labels {
		if strings.ToLower(label) == "ad-hoc-request" {
			return true
		}
	}
	return false
}

// FilterItemsByDateRange returns items filtered by the given date range and filter criteria
func FilterItemsByDateRange(
	items []models.KanbanItem, 
	startDate, endDate time.Time, 
	filterField models.FilterField, 
	adHocFilter types.AdHocFilterType, // Updated type
) []models.KanbanItem {
	var filtered []models.KanbanItem
	
	for _, item := range items {
		// Get the appropriate date field using the FilterField's method
		itemDate, hasDate := filterField.GetItemDate(item)
		
		// Skip items with no date in the requested field
		if !hasDate {
			continue
		}
		
		// Check if date is within range
		if (startDate.IsZero() || !itemDate.Before(startDate)) &&
		   (endDate.IsZero() || !itemDate.After(endDate)) {
			
			// Apply ad-hoc request filter
			isAdHoc := IsAdHocRequest(item)
			
			switch adHocFilter {
			case types.AdHocFilterInclude: // Updated constant
				filtered = append(filtered, item)
			case types.AdHocFilterExclude: // Updated constant
				if !isAdHoc {
					filtered = append(filtered, item)
				}
			case types.AdHocFilterOnly: // Updated constant
				if isAdHoc {
					filtered = append(filtered, item)
				}
			}
		}
	}
	
	return filtered
}