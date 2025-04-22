package metrics

import (
	"fmt"
	"sort"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

// WorkItemAgeReport shows how long current items have been in each state
func WorkItemAgeReport(items []models.KanbanItem, asOf time.Time) (string, error) {
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