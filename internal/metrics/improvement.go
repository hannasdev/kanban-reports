package metrics

import (
	"fmt"
	"sort"

	"github.com/hannasdev/kanban-reports/internal/models"
)

// TeamImprovementReport shows how metrics change month over month
func TeamImprovementReport(items []models.KanbanItem) (string, error) {
	// Group items by month
	itemsByMonth := make(map[string][]models.KanbanItem)
	
	for _, item := range items {
		if item.IsCompleted && !item.CompletedAt.IsZero() {
			month := item.CompletedAt.Format("2006-01")
			itemsByMonth[month] = append(itemsByMonth[month], item)
		}
	}
	
	// Sort months
	var months []string
	for month := range itemsByMonth {
		months = append(months, month)
	}
	sort.Strings(months)
	
	// Calculate metrics for each month
	type monthlyMetrics struct {
		ItemCount int
		StoryPoints float64
		AvgLeadTime float64
		AvgCycleTime float64
		LeadTimeMedian float64
		CycleTimeMedian float64
	}
	
	metricsByMonth := make(map[string]monthlyMetrics)
	
	for _, month := range months {
		monthItems := itemsByMonth[month]
		metrics := monthlyMetrics{
			ItemCount: len(monthItems),
		}
		
		var leadTimes, cycleTimes []float64
		
		for _, item := range monthItems {
			metrics.StoryPoints += item.Estimate
			
			if !item.CreatedAt.IsZero() && !item.CompletedAt.IsZero() {
				leadTime := item.CompletedAt.Sub(item.CreatedAt).Hours() / 24
				leadTimes = append(leadTimes, leadTime)
			}
			
			if !item.StartedAt.IsZero() && !item.CompletedAt.IsZero() {
				cycleTime := item.CompletedAt.Sub(item.StartedAt).Hours() / 24
				cycleTimes = append(cycleTimes, cycleTime)
			}
		}
		
		// Calculate lead time statistics
		if len(leadTimes) > 0 {
			_, _, avg, median := calculateStats(leadTimes)
			metrics.AvgLeadTime = avg
			metrics.LeadTimeMedian = median
		}
		
		// Calculate cycle time statistics
		if len(cycleTimes) > 0 {
			_, _, avg, median := calculateStats(cycleTimes)
			metrics.AvgCycleTime = avg
			metrics.CycleTimeMedian = median
		}
		
		metricsByMonth[month] = metrics
	}
	
	// Generate report with month-over-month changes
	report := "# Team Improvement Metrics\n\n"
	
	// Add explanatory text
	report += "## What are Team Improvement Metrics?\n\n"
	report += "Team Improvement Metrics track how your team's performance changes over time across several key dimensions. This helps identify trends, improvements, and areas that need attention.\n\n"
	report += "The metrics tracked month-over-month include:\n"
	report += "- **Item Count**: Number of completed items\n"
	report += "- **Story Points**: Total points completed\n"
	report += "- **Lead Time**: Average time from creation to completion\n"
	report += "- **Cycle Time**: Average time from start to completion\n\n"
	report += "## How to use this data:\n"
	report += "- Look for trends in delivery capacity (items and points)\n"
	report += "- Track improvements in lead time and cycle time\n"
	report += "- Use delta (Δ) values to see percentage improvements\n"
	report += "- Celebrate improvements and investigate regressions\n"
	report += "- Set team goals based on historical performance\n\n"
	
	report += "Month | Items | Points | Avg Lead Time | Avg Cycle Time | Lead Time Δ | Cycle Time Δ\n"
	report += "------|-------|--------|---------------|----------------|------------|-------------\n"
	
	var prevMonth string
	for _, month := range months {
		metrics := metricsByMonth[month]
		
		leadTimeChange := ""
		cycleTimeChange := ""
		
		if prevMonth != "" {
			prevMetrics := metricsByMonth[prevMonth]
			
			leadTimeDiff := metrics.AvgLeadTime - prevMetrics.AvgLeadTime
			if prevMetrics.AvgLeadTime > 0 {
				leadTimeChange = fmt.Sprintf("%+.1f (%+.1f%%)", 
					leadTimeDiff, 
					(leadTimeDiff/prevMetrics.AvgLeadTime)*100)
			}
			
			cycleTimeDiff := metrics.AvgCycleTime - prevMetrics.AvgCycleTime
			if prevMetrics.AvgCycleTime > 0 {
				cycleTimeChange = fmt.Sprintf("%+.1f (%+.1f%%)", 
					cycleTimeDiff, 
					(cycleTimeDiff/prevMetrics.AvgCycleTime)*100)
			}
		}
		
		report += fmt.Sprintf("%s | %5d | %6.1f | %13.1f | %14.1f | %10s | %11s\n",
			month, 
			metrics.ItemCount, 
			metrics.StoryPoints, 
			metrics.AvgLeadTime, 
			metrics.AvgCycleTime,
			leadTimeChange,
			cycleTimeChange)
		
		prevMonth = month
	}
	
	// Add statistical analysis section
	report += "\n## Statistical Trends\n\n"
	report += "Month | Lead Time (Median) | Cycle Time (Median) | Items/Month | Points/Month\n"
	report += "------|-------------------|-------------------|------------|-------------\n"
	
	for _, month := range months {
		metrics := metricsByMonth[month]
		report += fmt.Sprintf("%s | %17.1f | %19.1f | %10d | %11.1f\n",
			month,
			metrics.LeadTimeMedian,
			metrics.CycleTimeMedian,
			metrics.ItemCount,
			metrics.StoryPoints)
	}
	
	return report, nil
}