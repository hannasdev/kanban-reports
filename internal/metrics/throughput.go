package metrics

import (
	"fmt"
	"sort"

	"github.com/hannasdev/kanban-reports/internal/models"
)

// ThroughputReport shows items and points completed per time period
func ThroughputReport(items []models.KanbanItem, periodType string) (string, error) {
	// Group items by time period (week or month)
	periodFormat := "2006-01"
	periodName := "Month"
	if periodType == "week" {
		periodFormat = "2006-W02" // ISO week format
		periodName = "Week"
	}
	
	throughputByPeriod := make(map[string]struct{
		Count int
		Points float64
		Types map[string]int
	})
	
	for _, item := range items {
		if item.IsCompleted && !item.CompletedAt.IsZero() {
			period := item.CompletedAt.Format(periodFormat)
			
			periodData := throughputByPeriod[period]
			periodData.Count++
			periodData.Points += item.Estimate
			
			// Initialize types map if needed
			if periodData.Types == nil {
				periodData.Types = make(map[string]int)
			}
			
			// Count by type
			itemType := item.Type
			if itemType == "" {
				itemType = "Unspecified"
			}
			periodData.Types[itemType]++
			
			throughputByPeriod[period] = periodData
		}
	}
	
	// Sort periods chronologically
	var periods []string
	for period := range throughputByPeriod {
		periods = append(periods, period)
	}
	sort.Strings(periods)
	
	report := fmt.Sprintf("# Throughput Analysis by %s\n\n", periodName)
	
	// Add explanatory text
	report += "## What is Throughput?\n\n"
	report += "Throughput measures how many items your team completes in a given time period. It represents your delivery capacity and is a key metric for planning and forecasting.\n\n"
	report += "- Higher numbers indicate greater delivery capacity\n"
	report += "- Consistent throughput suggests a stable, predictable process\n"
	report += "- Declining throughput may indicate impediments or reduced capacity\n"
	report += "- Rising throughput may indicate process improvements or increased capacity\n\n"
	report += "## How to use this data:\n"
	report += "- Use average throughput to forecast future delivery capabilities\n"
	report += "- Look for trends or patterns in delivery capacity\n"
	report += "- Compare throughput across different time periods to identify improvements or issues\n"
	report += "- Analyze the balance between different types of work (features, bugs, etc.)\n\n"
	
	report += fmt.Sprintf("%s | Items Completed | Story Points | Avg Points/Item\n", periodName)
	report += "-------|----------------|-------------|---------------\n"
	
	for _, period := range periods {
		data := throughputByPeriod[period]
		avgPointsPerItem := 0.0
		if data.Count > 0 {
			avgPointsPerItem = data.Points / float64(data.Count)
		}
		
		report += fmt.Sprintf("%s | %15d | %11.1f | %14.1f\n", 
			period, data.Count, data.Points, avgPointsPerItem)
	}
	
	// Add breakdown by type
	report += "\n## Breakdown by Item Type\n\n"
	
	// Get all unique types across all periods
	allTypes := make(map[string]bool)
	for _, period := range periods {
		for itemType := range throughputByPeriod[period].Types {
			allTypes[itemType] = true
		}
	}
	
	// Convert to sorted slice
	var typesList []string
	for itemType := range allTypes {
		typesList = append(typesList, itemType)
	}
	sort.Strings(typesList)
	
	// Create header with all types
	report += periodName
	for _, itemType := range typesList {
		report += fmt.Sprintf(" | %s", itemType)
	}
	report += " | Total\n"
	
	// Add separator
	report += "-------"
	for range typesList {
		report += "|-------"
	}
	report += "|-------\n"
	
	// Add rows for each period
	for _, period := range periods {
		data := throughputByPeriod[period]
		report += period
		
		periodTotal := 0
		for _, itemType := range typesList {
			count := data.Types[itemType]
			report += fmt.Sprintf(" | %5d", count)
			periodTotal += count
		}
		
		report += fmt.Sprintf(" | %5d\n", periodTotal)
	}
	
	return report, nil
}