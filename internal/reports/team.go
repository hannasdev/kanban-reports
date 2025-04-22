package reports

import (
	"fmt"
	"sort"

	"github.com/hannasdev/kanban-reports/internal/models"
)

// generateTeamReport creates a report of story points by team
func (r *Reporter) generateTeamReport(items []models.KanbanItem) (string, error) {
	// Map to track points by team
	teamPoints := make(map[string]float64)
	teamItems := make(map[string]int)
	
	// Calculate points by team
	for _, item := range items {
		teamName := item.Team
		if teamName == "" {
			teamName = "No Team"
		}
		
		teamPoints[teamName] += item.Estimate
		teamItems[teamName]++
	}
	
	// Sort teams by points
	type teamStat struct {
		name      string
		points    float64
		itemCount int
	}
	
	var stats []teamStat
	for name, points := range teamPoints {
		stats = append(stats, teamStat{
			name:      name,
			points:    points,
			itemCount: teamItems[name],
		})
	}
	
	// Sort by points in descending order
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].points > stats[j].points
	})
	
	// Generate report string
	report := "Story Points by Team:\n\n"
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