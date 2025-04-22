package reports

import (
	"fmt"
	"sort"

	"github.com/hannasdev/kanban-reports/internal/models"
)

// generateEpicReport creates a report of story points by epic
func (r *Reporter) generateEpicReport(items []models.KanbanItem) (string, error) {
	// Map to track points by epic
	epicPoints := make(map[string]float64)
	epicItems := make(map[string]int)
	
	// Calculate points by epic
	for _, item := range items {
		epicName := item.Epic
		if epicName == "" {
			epicName = "No Epic"
		}
		
		epicPoints[epicName] += item.Estimate
		epicItems[epicName]++
	}
	
	// Sort epics by points
	type epicStat struct {
		name      string
		points    float64
		itemCount int
	}
	
	var stats []epicStat
	for name, points := range epicPoints {
		stats = append(stats, epicStat{
			name:      name,
			points:    points,
			itemCount: epicItems[name],
		})
	}
	
	// Sort by points in descending order
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].points > stats[j].points
	})
	
	// Generate report string
	report := "Story Points by Epic:\n\n"
	totalPoints := 0.0
	totalItems := 0
	
	for _, stat := range stats {
		report += fmt.Sprintf("%-50s %6.1f points  %3d items\n", 
			stat.name, stat.points, stat.itemCount)
		totalPoints += stat.points
		totalItems += stat.itemCount
	}
	
	report += fmt.Sprintf("\nTotal: %.1f points across %d items\n", totalPoints, totalItems)
	
	return report, nil
}