package reports

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

// MetricsType defines the type of metrics to generate
type MetricsType string

const (
	// MetricsTypeLeadTime generates lead time analysis by story point size
	MetricsTypeLeadTime MetricsType = "lead-time"
	// MetricsTypeThroughput generates throughput analysis over time
	MetricsTypeThroughput MetricsType = "throughput"
	// MetricsTypeFlow generates flow efficiency analysis
	MetricsTypeFlow MetricsType = "flow"
	// MetricsTypeEstimation generates estimation accuracy analysis
	MetricsTypeEstimation MetricsType = "estimation"
	// MetricsTypeAge generates current work item age analysis
	MetricsTypeAge MetricsType = "age"
	// MetricsTypeImprovement generates month-over-month improvement metrics
	MetricsTypeImprovement MetricsType = "improvement"
	// MetricsTypeAll generates all metrics reports
	MetricsTypeAll MetricsType = "all"
)

// PeriodType defines the time period for grouping metrics
type PeriodType string

const (
	// PeriodTypeWeek groups metrics by week
	PeriodTypeWeek PeriodType = "week"
	// PeriodTypeMonth groups metrics by month
	PeriodTypeMonth PeriodType = "month"
)

// Standard story point sizes for grouping
var standardPointSizes = []float64{1, 2, 3, 5, 8, 13, 21}

// GenerateMetrics generates metrics based on the specified type and time period
func (r *Reporter) GenerateMetrics(metricsType MetricsType, periodType PeriodType, startDate, endDate time.Time) (string, error) {
	// Filter items by completion date within range
	filteredItems := r.filterItemsByDateRange(startDate, endDate)
	
	if len(filteredItems) == 0 {
		return "No items completed in the specified date range.", nil
	}

	// Generate appropriate metrics based on type
	switch metricsType {
	case MetricsTypeLeadTime:
		return r.leadTimeReport(filteredItems)
	case MetricsTypeThroughput:
		return r.throughputReport(filteredItems, string(periodType))
	case MetricsTypeFlow:
		return r.flowEfficiencyReport(filteredItems)
	case MetricsTypeEstimation:
		return r.estimationAccuracyReport(filteredItems)
	case MetricsTypeAge:
		return r.workItemAgeReport(filteredItems, time.Now())
	case MetricsTypeImprovement:
		return r.teamImprovementReport(filteredItems)
	case MetricsTypeAll:
		// Generate all reports and combine them
		reports := []string{}
		
		leadTime, err := r.leadTimeReport(filteredItems)
		if err == nil {
			reports = append(reports, leadTime)
		}
		
		throughput, err := r.throughputReport(filteredItems, string(periodType))
		if err == nil {
			reports = append(reports, throughput)
		}
		
		flow, err := r.flowEfficiencyReport(filteredItems)
		if err == nil {
			reports = append(reports, flow)
		}
		
		estimation, err := r.estimationAccuracyReport(filteredItems)
		if err == nil {
			reports = append(reports, estimation)
		}
		
		age, err := r.workItemAgeReport(filteredItems, time.Now())
		if err == nil {
			reports = append(reports, age)
		}
		
		improvement, err := r.teamImprovementReport(filteredItems)
		if err == nil {
			reports = append(reports, improvement)
		}
		
		return combineReports(reports), nil
		
	default:
		return "", fmt.Errorf("unknown metrics type: %s", metricsType)
	}
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

// LeadTimeReport shows how long items take from creation to completion
func (r *Reporter) leadTimeReport(items []models.KanbanItem) (string, error) {
	// Group by story point size
	leadTimesByPoints := make(map[float64][]float64)
	cycleTimesByPoints := make(map[float64][]float64)
	
	for _, item := range items {
		if item.IsCompleted && !item.CompletedAt.IsZero() {
			// Skip items with invalid dates
			if item.CreatedAt.IsZero() {
				continue
			}
			
			// Calculate lead time (created to completed)
			leadTime := item.CompletedAt.Sub(item.CreatedAt).Hours() / 24 // in days
			
			// Calculate cycle time (started to completed) if available
			if !item.StartedAt.IsZero() {
				cycleTime := item.CompletedAt.Sub(item.StartedAt).Hours() / 24 // in days
				
				// Find closest standard point size
				closestSize := findClosestPointSize(item.Estimate, standardPointSizes)
				cycleTimesByPoints[closestSize] = append(cycleTimesByPoints[closestSize], cycleTime)
			}
			
			// Find closest standard point size
			closestSize := findClosestPointSize(item.Estimate, standardPointSizes)
			leadTimesByPoints[closestSize] = append(leadTimesByPoints[closestSize], leadTime)
		}
	}
	
	// Calculate statistics for each point size
	report := "# Lead Time Analysis by Story Point Size (in days)\n\n"
	report += "## Lead Time (Creation to Completion)\n\n"
	report += "Story points | Count | Min | Max | Avg | Median\n"
	report += "-------------|-------|-----|-----|-----|-------\n"
	
	// Process all standard point sizes, even if we don't have data for some
	for _, size := range standardPointSizes {
		times := leadTimesByPoints[size]
		if len(times) == 0 {
			continue
		}
		
		min, max, avg, median := calculateStats(times)
		report += fmt.Sprintf("%12.0f | %5d | %3.1f | %3.1f | %3.1f | %6.1f\n", 
			size, len(times), min, max, avg, median)
	}
	
	// Add cycle time statistics
	report += "\n## Cycle Time (Start to Completion)\n\n"
	report += "Story points | Count | Min | Max | Avg | Median\n"
	report += "-------------|-------|-----|-----|-----|-------\n"
	
	for _, size := range standardPointSizes {
		times := cycleTimesByPoints[size]
		if len(times) == 0 {
			continue
		}
		
		min, max, avg, median := calculateStats(times)
		report += fmt.Sprintf("%12.0f | %5d | %3.1f | %3.1f | %3.1f | %6.1f\n", 
			size, len(times), min, max, avg, median)
	}
	
	return report, nil
}

// ThroughputReport shows items and points completed per time period
func (r *Reporter) throughputReport(items []models.KanbanItem, periodType string) (string, error) {
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

// FlowEfficiencyReport analyzes time spent in each state
func (r *Reporter) flowEfficiencyReport(items []models.KanbanItem) (string, error) {
	// Track time spent in each state
	stateTimeTotal := make(map[string]float64) // in days
	stateItemCount := make(map[string]int)
	
	for _, item := range items {
		if item.IsCompleted && !item.CompletedAt.IsZero() && !item.CreatedAt.IsZero() {
			// Simplified flow: Created -> Started -> Completed
			waitTime := 0.0
			activeTime := 0.0
			
			if !item.StartedAt.IsZero() {
				waitTime = item.StartedAt.Sub(item.CreatedAt).Hours() / 24
				activeTime = item.CompletedAt.Sub(item.StartedAt).Hours() / 24
			} else {
				// If no start time, consider all as active time
				activeTime = item.CompletedAt.Sub(item.CreatedAt).Hours() / 24
			}
			
			stateTimeTotal["Waiting"] += waitTime
			stateTimeTotal["Active"] += activeTime
			stateItemCount["Waiting"]++
			stateItemCount["Active"]++
		}
	}
	
	report := "# Flow Efficiency Analysis\n\n"
	report += "State | Avg Time (days) | % of Total Time\n"
	report += "------|-----------------|---------------\n"
	
	totalTime := stateTimeTotal["Waiting"] + stateTimeTotal["Active"]
	if totalTime > 0 {
		waitAvg := 0.0
		if stateItemCount["Waiting"] > 0 {
			waitAvg = stateTimeTotal["Waiting"] / float64(stateItemCount["Waiting"])
		}
		
		activeAvg := 0.0
		if stateItemCount["Active"] > 0 {
			activeAvg = stateTimeTotal["Active"] / float64(stateItemCount["Active"])
		}
		
		waitPercent := (stateTimeTotal["Waiting"] / totalTime) * 100
		activePercent := (stateTimeTotal["Active"] / totalTime) * 100
		
		report += fmt.Sprintf("Waiting | %15.1f | %13.1f%%\n", waitAvg, waitPercent)
		report += fmt.Sprintf("Active  | %15.1f | %13.1f%%\n", activeAvg, activePercent)
		report += fmt.Sprintf("\nFlow Efficiency: %.1f%%\n", activePercent)
		report += "\nNote: Flow efficiency is the percentage of time items spend in active states versus waiting states."
	}
	
	return report, nil
}

// EstimationAccuracyReport compares story point sizes to actual completion times
func (r *Reporter) estimationAccuracyReport(items []models.KanbanItem) (string, error) {
	// Map story points to actual cycle times
	cycleTimesByPoints := make(map[float64][]float64)
	
	for _, item := range items {
		if item.IsCompleted && !item.CompletedAt.IsZero() && !item.StartedAt.IsZero() {
			cycleTime := item.CompletedAt.Sub(item.StartedAt).Hours() / 24 // in days
			
			// Find closest standard point size
			closestSize := findClosestPointSize(item.Estimate, standardPointSizes)
			cycleTimesByPoints[closestSize] = append(cycleTimesByPoints[closestSize], cycleTime)
		}
	}
	
	// Calculate cycle time per story point for each size
	report := "# Estimation Accuracy Analysis\n\n"
	report += "## Time Spent per Story Point Size\n\n"
	report += "Story points | Count | Min Days/SP | Max Days/SP | Avg Days/SP | Median Days/SP\n"
	report += "-------------|-------|------------|------------|-------------|---------------\n"
	
	for _, size := range standardPointSizes {
		times := cycleTimesByPoints[size]
		if len(times) == 0 || size == 0 {
			continue
		}
		
		// Convert to days per story point
		daysPerSP := make([]float64, len(times))
		for i, t := range times {
			daysPerSP[i] = t / size
		}
		
		min, max, avg, median := calculateStats(daysPerSP)
		report += fmt.Sprintf("%12.0f | %5d | %10.1f | %10.1f | %11.1f | %15.1f\n", 
			size, len(times), min, max, avg, median)
	}
	
	// Add raw cycle time data for comparison
	report += "\n## Raw Cycle Time by Story Point Size\n\n"
	report += "Story points | Count | Min | Max | Avg | Median\n"
	report += "-------------|-------|-----|-----|-----|-------\n"
	
	for _, size := range standardPointSizes {
		times := cycleTimesByPoints[size]
		if len(times) == 0 {
			continue
		}
		
		min, max, avg, median := calculateStats(times)
		report += fmt.Sprintf("%12.0f | %5d | %3.1f | %3.1f | %3.1f | %6.1f\n", 
			size, len(times), min, max, avg, median)
	}
	
	// Calculate overall correlation between story points and cycle time
	var allPoints []float64
	var allTimes []float64
	
	for size, times := range cycleTimesByPoints {
		if size == 0 {
			continue
		}
		
		for _, t := range times {
			allPoints = append(allPoints, size)
			allTimes = append(allTimes, t)
		}
	}
	
	if len(allPoints) > 0 {
		correlation := calculateCorrelation(allPoints, allTimes)
		report += fmt.Sprintf("\nCorrelation between story points and cycle time: %.2f\n", correlation)
		report += "\nNote: A correlation closer to 1.0 indicates that higher point items consistently take more time."
	}
	
	return report, nil
}

// WorkItemAgeReport shows how long current items have been in each state
func (r *Reporter) workItemAgeReport(items []models.KanbanItem, asOf time.Time) (string, error) {
	if asOf.IsZero() {
		asOf = time.Now()
	}
	
	// Group items by state
	stateItems := make(map[string][]struct{
		Name string
		Age float64
	})
	
	for _, item := range items {
		if item.IsCompleted {
			continue // Skip completed items
		}
		
		var age float64
		if !item.StartedAt.IsZero() {
			age = asOf.Sub(item.StartedAt).Hours() / 24
		} else {
			age = asOf.Sub(item.CreatedAt).Hours() / 24
		}
		
		state := item.State
		if state == "" {
			state = "Unknown"
		}
		
		stateItems[state] = append(stateItems[state], struct{
			Name string
			Age float64
		}{item.Name, age})
	}
	
	// Sort states
	var states []string
	for state := range stateItems {
		states = append(states, state)
	}
	sort.Strings(states)
	
	// Generate report
	report := "# Current Work Item Age Analysis\n\n"
	report += "Age of incomplete items by state (in days):\n\n"
	
	for _, state := range states {
		items := stateItems[state]
		if len(items) == 0 {
			continue
		}
		
		report += fmt.Sprintf("## %s (%d items)\n\n", state, len(items))
		
		// Sort by age (descending)
		sort.Slice(items, func(i, j int) bool {
			return items[i].Age > items[j].Age
		})
		
		// Calculate statistics
		var ages []float64
		for _, item := range items {
			ages = append(ages, item.Age)
		}
		min, max, avg, median := calculateStats(ages)
		
		report += fmt.Sprintf("Min: %.1f, Max: %.1f, Avg: %.1f, Median: %.1f days\n\n", 
			min, max, avg, median)
		
		// Show oldest 5 items
		report += "Oldest Items:\n\n"
		for i, item := range items {
			if i >= 5 {
				break
			}
			report += fmt.Sprintf("- %s (%.1f days)\n", item.Name, item.Age)
		}
		report += "\n"
	}
	
	return report, nil
}

// TeamImprovementReport shows how metrics change month over month
func (r *Reporter) teamImprovementReport(items []models.KanbanItem) (string, error) {
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

// Helper function to find closest story point size from standard sizes
func findClosestPointSize(estimate float64, standardSizes []float64) float64 {
	if estimate == 0 {
		return 0
	}
	
	closest := standardSizes[0]
	minDiff := math.Abs(estimate - closest)
	
	for _, size := range standardSizes {
		diff := math.Abs(estimate - size)
		if diff < minDiff {
			minDiff = diff
			closest = size
		}
	}
	
	return closest
}

// Calculate statistical values from a set of data points
func calculateStats(values []float64) (min, max, avg, median float64) {
	if len(values) == 0 {
		return 0, 0, 0, 0
	}
	
	// Sort for min, max, median
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	
	min = sorted[0]
	max = sorted[len(sorted)-1]
	
	// Calculate average
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	avg = sum / float64(len(values))
	
	// Calculate median
	if len(sorted)%2 == 0 {
		// Even number of values
		middle1 := sorted[len(sorted)/2-1]
		middle2 := sorted[len(sorted)/2]
		median = (middle1 + middle2) / 2
	} else {
		// Odd number of values
		median = sorted[len(sorted)/2]
	}
	
	return min, max, avg, median
}

// Calculate correlation between two sets of values
func calculateCorrelation(x, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0
	}
	
	n := float64(len(x))
	
	// Calculate sums
	sumX := 0.0
	sumY := 0.0
	sumXY := 0.0
	sumX2 := 0.0
	sumY2 := 0.0
	
	for i := range x {
		sumX += x[i]
		sumY += y[i]
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
		sumY2 += y[i] * y[i]
	}
	
	// Calculate correlation coefficient
	numerator := n*sumXY - sumX*sumY
	denominator := math.Sqrt((n*sumX2 - sumX*sumX) * (n*sumY2 - sumY*sumY))
	
	if denominator == 0 {
		return 0
	}
	
	return numerator / denominator
}