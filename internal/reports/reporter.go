package reports

import (
	"fmt"
	"strings"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

// Reporter handles generation of different reports
type Reporter struct {
    items      []models.KanbanItem
    adHocFilter AdHocFilterType
}

// NewReporter creates a new reporter with the given items
func NewReporter(items []models.KanbanItem) *Reporter {
    return &Reporter{
        items:      items,
        adHocFilter: AdHocFilterInclude,
    }
}

// WithAdHocFilter sets the ad-hoc request filter
func (r *Reporter) WithAdHocFilter(filter AdHocFilterType) *Reporter {
    r.adHocFilter = filter
    return r
}

// GenerateReport generates a report based on the specified type and time period
func (r *Reporter) GenerateReport(reportType ReportType, startDate, endDate time.Time) (string, error) {
    // Filter items by completion date within range
    filteredItems := r.filterItemsByDateRange(startDate, endDate)
    
    if len(filteredItems) == 0 {
        return "No items completed in the specified date range.", nil
    }

    // Generate appropriate report based on type
    switch reportType {
    case ReportTypeContributor:
        return r.generateContributorReport(filteredItems)
    case ReportTypeEpic:
        return r.generateEpicReport(filteredItems)
    case ReportTypeProductArea:
        return r.generateProductAreaReport(filteredItems)
    case ReportTypeTeam:
        return r.generateTeamReport(filteredItems)
    default:
        return "", fmt.Errorf("unknown report type: %s", reportType)
    }
}

// filterItemsByDateRange returns items completed within the given date range
func (r *Reporter) filterItemsByDateRange(startDate, endDate time.Time) []models.KanbanItem {
    var filtered []models.KanbanItem
    
    for _, item := range r.items {
        // Only include completed items
        if !item.IsCompleted || item.CompletedAt.IsZero() {
            continue
        }
        
        // Check if completion date is within range
        if (startDate.IsZero() || !item.CompletedAt.Before(startDate)) &&
           (endDate.IsZero() || !item.CompletedAt.After(endDate)) {
            
            // Apply ad-hoc request filter
            isAdHoc := r.isAdHocRequest(item)
            
            if (r.adHocFilter == AdHocFilterInclude) ||
               (r.adHocFilter == AdHocFilterExclude && !isAdHoc) ||
               (r.adHocFilter == AdHocFilterOnly && isAdHoc) {
                filtered = append(filtered, item)
            }
        }
    }
    
    return filtered
}

// isAdHocRequest checks if an item is an ad-hoc request (has "ad-hoc-request" label)
func (r *Reporter) isAdHocRequest(item models.KanbanItem) bool {
    for _, label := range item.Labels {
        if strings.ToLower(label) == "ad-hoc-request" {
            return true
        }
    }
    return false
}