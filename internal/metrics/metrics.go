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

// Generate generates metrics based on the specified type and time period
func (g *Generator) Generate(metricsType MetricsType, periodType PeriodType, startDate, endDate time.Time) (string, error) {
	// Filter items by completion date within range
	filteredItems := g.filterItemsByDateRange(startDate, endDate)
	
	if len(filteredItems) == 0 {
		return "No items completed in the specified date range.", nil
	}

	// Generate appropriate metrics based on type
	switch metricsType {
	case MetricsTypeLeadTime:
		return LeadTimeReport(filteredItems)
	case MetricsTypeThroughput:
		return ThroughputReport(filteredItems, string(periodType))
	case MetricsTypeFlow:
		return FlowEfficiencyReport(filteredItems)
	case MetricsTypeEstimation:
		return EstimationAccuracyReport(filteredItems)
	case MetricsTypeAge:
		return WorkItemAgeReport(filteredItems, time.Now())
	case MetricsTypeImprovement:
		return TeamImprovementReport(filteredItems)
	case MetricsTypeAll:
		return GenerateAllReports(filteredItems, string(periodType))
	default:
		return "", fmt.Errorf("unknown metrics type: %s", metricsType)
	}
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