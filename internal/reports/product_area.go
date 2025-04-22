package reports

import (
	"fmt"
	"sort"

	"github.com/hannasdev/kanban-reports/internal/models"
)

// generateProductAreaReport creates a report of story points by product area
func (r *Reporter) generateProductAreaReport(items []models.KanbanItem) (string, error) {
	// Map to track points by product area
	areaPoints := make(map[string]float64)
	areaItems := make(map[string]int)
	
	// Calculate points by product area
	for _, item := range items {
		areaName := item.ProductArea
		if areaName == "" {
			areaName = "Uncategorized"
		}
		
		areaPoints[areaName] += item.Estimate
		areaItems[areaName]++
	}
	
	// Sort areas by points
	type areaStat struct {
		name      string
		points    float64
		itemCount int
	}
	
	var stats []areaStat
	for name, points := range areaPoints {
		stats = append(stats, areaStat{
			name:      name,
			points:    points,
			itemCount: areaItems[name],
		})
	}
	
	// Sort by points in descending order
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].points > stats[j].points
	})
	
	// Generate report string
	report := "Story Points by Product Area:\n\n"
	totalPoints := 0.0
	totalItems := 0
	
	for _, stat := range stats {
		report += fmt.Sprintf("%-30s %6.1f points  %3d items\n", 
			stat.name, stat.points, stat.itemCount)
		totalPoints += stat.points
		totalItems += stat.itemCount
	}
	
	report += fmt.Sprintf("\nTotal: %.1f points across %d items\n", totalPoints, totalItems)
	
	return report, nil
}