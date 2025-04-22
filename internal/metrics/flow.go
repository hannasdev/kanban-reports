package metrics

import (
	"fmt"

	"github.com/hannasdev/kanban-reports/internal/models"
)

// FlowEfficiencyReport analyzes time spent in each state
func FlowEfficiencyReport(items []models.KanbanItem) (string, error) {
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
	
	// Add explanatory text
	report += "## What is Flow Efficiency?\n\n"
	report += "Flow efficiency measures the percentage of time items spend being actively worked on versus waiting. It's a key metric in Lean and Kanban methodologies.\n\n"
	report += "**Flow Efficiency = (Active Time / Total Time) × 100%**\n\n"
	report += "- **Waiting time**: Time from creation until work begins (queue time)\n"
	report += "- **Active time**: Time from when work begins until completion (processing time)\n"
	report += "- **Total time**: Sum of waiting and active time (lead time)\n\n"
	report += "## Interpretation:\n"
	report += "- **Low efficiency (10-30%)**: Common in many organizations. Items spend most of their time waiting.\n"
	report += "- **Medium efficiency (30-50%)**: Generally considered good for knowledge work.\n"
	report += "- **High efficiency (>50%)**: Excellent! Your team has minimal waiting time.\n\n"
	report += "## How to improve flow efficiency:\n"
	report += "- Limit work in progress (WIP)\n"
	report += "- Reduce handoffs between teams\n"
	report += "- Eliminate bottlenecks\n"
	report += "- Implement pull systems\n"
	report += "- Reduce batch sizes\n\n"
	
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
	} else {
		report += "No data available for flow efficiency calculation.\n"
	}
	
	// Additional advanced flow analysis could go here
	// For example, analyzing flow efficiency by story point size or type
	
	return report, nil
}