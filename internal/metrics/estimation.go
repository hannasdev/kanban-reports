package metrics

import (
	"fmt"

	"github.com/hannasdev/kanban-reports/internal/models"
)

// EstimationAccuracyReport compares story point sizes to actual completion times
func EstimationAccuracyReport(items []models.KanbanItem) (string, error) {
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
	
	// Add explanatory text
	report += "## What is Estimation Accuracy?\n\n"
	report += "Estimation accuracy measures how well your story point estimates correlate with the actual time spent completing work. Ideally, there should be a consistent relationship between story points and completion time.\n\n"
	report += "This analysis shows:\n"
	report += "- How much time is spent per story point for different sized items\n"
	report += "- Whether your story points consistently scale (e.g., do 3-point stories take about 3× as long as 1-point stories?)\n"
	report += "- The correlation between estimates and actual completion times\n\n"
	report += "## How to use this data:\n"
	report += "- Look for consistency in the days/SP metric across different sizes\n"
	report += "- Identify if certain sized items are consistently under or overestimated\n"
	report += "- Use the correlation value to assess your estimation system's reliability\n"
	report += "- Consider calibrating story point values based on actual completion times\n\n"
	
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
		report += "\nInterpretation of correlation:\n"
		report += "- **0.7-1.0**: Strong positive correlation. Excellent estimation system.\n"
		report += "- **0.4-0.7**: Moderate correlation. Reasonably good estimates.\n"
		report += "- **0.0-0.4**: Weak correlation. Estimates may need improvement.\n"
		report += "- **Negative**: Inverse relationship. Estimation system needs significant review.\n"
	}
	
	return report, nil
}