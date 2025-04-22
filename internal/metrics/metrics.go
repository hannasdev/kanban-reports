package metrics

import (
	"fmt"
	"strings"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
	"github.com/hannasdev/kanban-reports/internal/reports"
)

// Generator handles the generation of metrics
type Generator struct {
	items       []models.KanbanItem
	adHocFilter reports.AdHocFilterType
}

// NewGenerator creates a new metrics generator
func NewGenerator(items []models.KanbanItem) *Generator {
	return &Generator{
		items:       items,
		adHocFilter: reports.AdHocFilterInclude,
	}
}

// WithAdHocFilter sets the ad-hoc request filter
func (g *Generator) WithAdHocFilter(filter reports.AdHocFilterType) *Generator {
	g.adHocFilter = filter
	return g
}

// filterItemsByDateRange returns items completed within the given date range
func (g *Generator) filterItemsByDateRange(startDate, endDate time.Time) []models.KanbanItem {
	var filtered []models.KanbanItem
	
	for _, item := range g.items {
		// Only include completed items
		if !item.IsCompleted || item.CompletedAt.IsZero() {
			continue
		}
		
		// Check if completion date is within range
		if (startDate.IsZero() || !item.CompletedAt.Before(startDate)) &&
		   (endDate.IsZero() || !item.CompletedAt.After(endDate)) {
			
			// Apply ad-hoc request filter
			isAdHoc := g.isAdHocRequest(item)
			
			if (g.adHocFilter == reports.AdHocFilterInclude) ||
			   (g.adHocFilter == reports.AdHocFilterExclude && !isAdHoc) ||
			   (g.adHocFilter == reports.AdHocFilterOnly && isAdHoc) {
				filtered = append(filtered, item)
			}
		}
	}
	
	return filtered
}

// isAdHocRequest checks if an item is an ad-hoc request (has "ad-hoc-request" label)
func (g *Generator) isAdHocRequest(item models.KanbanItem) bool {
	for _, label := range item.Labels {
		if strings.ToLower(label) == "ad-hoc-request" {
			return true
		}
	}
	return false
}

// addDateRangeInfo adds date range information to the beginning of the metrics report
func (g *Generator) addDateRangeInfo(report string, metricsType MetricsType, periodType PeriodType, startDate, endDate time.Time) string {
	// Create header with metrics type and date information
	var header string
	
	// Format the header with date range information
	if !startDate.IsZero() && !endDate.IsZero() {
		header = fmt.Sprintf("Metrics Type: %s\nPeriod Type: %s\nDate Range: %s to %s\n\n", 
			metricsType, 
			periodType,
			startDate.Format("2006-01-02"), 
			endDate.Format("2006-01-02"))
	} else if !startDate.IsZero() {
		header = fmt.Sprintf("Metrics Type: %s\nPeriod Type: %s\nFrom: %s\n\n", 
			metricsType,
			periodType, 
			startDate.Format("2006-01-02"))
	} else if !endDate.IsZero() {
		header = fmt.Sprintf("Metrics Type: %s\nPeriod Type: %s\nTo: %s\n\n", 
			metricsType,
			periodType, 
			endDate.Format("2006-01-02"))
	} else {
		header = fmt.Sprintf("Metrics Type: %s\nPeriod Type: %s\nDate Range: All Time\n\n", 
			metricsType,
			periodType)
	}
	
	// Add ad-hoc filtering information
	switch g.adHocFilter {
	case reports.AdHocFilterExclude:
		header += "Filter: Excluding ad-hoc requests\n\n"
	case reports.AdHocFilterOnly:
		header += "Filter: Only ad-hoc requests\n\n"
	}
	
	return header + report
}

// Generate generates metrics based on the specified type and time period
func (g *Generator) Generate(metricsType MetricsType, periodType PeriodType, startDate, endDate time.Time) (string, error) {
	// Filter items by completion date within range
	filteredItems := g.filterItemsByDateRange(startDate, endDate)
	
	if len(filteredItems) == 0 {
		return "No items completed in the specified date range.", nil
	}

	// Generate appropriate metrics based on type
	var metricsContent string
	var err error

	switch metricsType {
	case MetricsTypeLeadTime:
		metricsContent, err = LeadTimeReport(filteredItems)
	case MetricsTypeThroughput:
		metricsContent, err = ThroughputReport(filteredItems, string(periodType))
	case MetricsTypeFlow:
		metricsContent, err = FlowEfficiencyReport(filteredItems)
	case MetricsTypeEstimation:
		metricsContent, err = EstimationAccuracyReport(filteredItems)
	case MetricsTypeAge:
		metricsContent, err = WorkItemAgeReport(filteredItems, time.Now())
	case MetricsTypeImprovement:
		metricsContent, err = TeamImprovementReport(filteredItems)
	case MetricsTypeAll:
		metricsContent, err = GenerateAllReports(filteredItems, string(periodType))
	default:
		return "", fmt.Errorf("unknown metrics type: %s", metricsType)
	}

	if err != nil {
		return "", err
	}

	// Add date range information to the metrics report
	reportWithDateInfo := g.addDateRangeInfo(metricsContent, metricsType, periodType, startDate, endDate)
	return reportWithDateInfo, nil
}

// GenerateAllReports generates all types of metrics reports
func GenerateAllReports(items []models.KanbanItem, periodType string) (string, error) {
	// Generate all reports and combine them
	reports := []string{}
	
	leadTime, err := LeadTimeReport(items)
	if err == nil {
		reports = append(reports, leadTime)
	}
	
	throughput, err := ThroughputReport(items, periodType)
	if err == nil {
		reports = append(reports, throughput)
	}
	
	flow, err := FlowEfficiencyReport(items)
	if err == nil {
		reports = append(reports, flow)
	}
	
	estimation, err := EstimationAccuracyReport(items)
	if err == nil {
		reports = append(reports, estimation)
	}
	
	age, err := WorkItemAgeReport(items, time.Now())
	if err == nil {
		reports = append(reports, age)
	}
	
	improvement, err := TeamImprovementReport(items)
	if err == nil {
		reports = append(reports, improvement)
	}
	
	return combineReports(reports), nil
}

// combineReports combines multiple report strings with separators
func combineReports(reports []string) string {
	combined := ""
	separator := "\n\n" + strings.Repeat("=", 80) + "\n\n"
	
	for i, report := range reports {
		combined += report
		if i < len(reports)-1 {
			combined += separator
		}
	}
	
	return combined
}