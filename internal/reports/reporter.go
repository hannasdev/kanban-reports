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
    var reportContent string
    var err error

    switch reportType {
    case ReportTypeContributor:
        reportContent, err = r.generateContributorReport(filteredItems)
    case ReportTypeEpic:
        reportContent, err = r.generateEpicReport(filteredItems)
    case ReportTypeProductArea:
        reportContent, err = r.generateProductAreaReport(filteredItems)
    case ReportTypeTeam:
        reportContent, err = r.generateTeamReport(filteredItems)
    default:
        return "", fmt.Errorf("unknown report type: %s", reportType)
    }

		if err != nil {
        return "", err
    }

		 // Add date range information to the report
    reportWithDateInfo := r.addDateRangeInfo(reportContent, reportType, startDate, endDate)
    return reportWithDateInfo, nil
}

// addDateRangeInfo adds date range information to the beginning of the report
func (r *Reporter) addDateRangeInfo(report string, reportType ReportType, startDate, endDate time.Time) string {
    // Create header with report type and date information
    var header string
    reportTypeName := string(reportType)
    
    // Format the header with date range information
    if !startDate.IsZero() && !endDate.IsZero() {
        header = fmt.Sprintf("Report Type: %s\nDate Range: %s to %s\n\n", 
            reportTypeName, 
            startDate.Format("2006-01-02"), 
            endDate.Format("2006-01-02"))
    } else if !startDate.IsZero() {
        header = fmt.Sprintf("Report Type: %s\nFrom: %s\n\n", 
            reportTypeName, 
            startDate.Format("2006-01-02"))
    } else if !endDate.IsZero() {
        header = fmt.Sprintf("Report Type: %s\nTo: %s\n\n", 
            reportTypeName, 
            endDate.Format("2006-01-02"))
    } else {
        header = fmt.Sprintf("Report Type: %s\nDate Range: All Time\n\n", 
            reportTypeName)
    }
    
    // Add ad-hoc filtering information
    switch r.adHocFilter {
    case AdHocFilterExclude:
        header += "Filter: Excluding ad-hoc requests\n\n"
    case AdHocFilterOnly:
        header += "Filter: Only ad-hoc requests\n\n"
    }
    
    return header + report
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