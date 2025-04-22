package metrics

import (
	"fmt"

	"github.com/hannasdev/kanban-reports/internal/models"
)

// LeadTimeReport shows how long items take from creation to completion
func LeadTimeReport(items []models.KanbanItem) (string, error) {
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
	
	// Add explanatory text
	report += "## What is Lead Time?\n\n"
	report += "Lead Time measures how long it takes for work to go from creation to completion. It's the total elapsed time that a customer waits for their request to be delivered.\n\n"
	report += "- **Lead Time**: Time from when an item is created to when it's completed (includes both waiting and active time)\n"
	report += "- **Cycle Time**: Time from when work actively starts on an item to when it's completed (active time only)\n\n"
	report += "Lower values indicate faster delivery. Higher story point items typically have longer lead and cycle times.\n\n"
	report += "## How to use this data:\n"
	report += "- Compare different sized items to validate your estimation system\n"
	report += "- Use these values to set realistic delivery expectations with stakeholders\n"
	report += "- Track these metrics over time to identify process improvements\n\n"
	
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