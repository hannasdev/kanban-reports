package reports

import (
	"fmt"
	"sort"

	"github.com/hannasdev/kanban-reports/internal/models"
)

// generateContributorReport creates a report of story points by contributor
func (r *Reporter) generateContributorReport(items []models.KanbanItem) (string, error) {
    // Map to track points by contributor
    contributorPoints := make(map[string]float64)
    contributorItems := make(map[string]int)
    
    // Calculate points by contributor
    for _, item := range items {
        // If no owners, credit to "Unassigned"
        if len(item.Owners) == 0 {
            contributorPoints["Unassigned"] += item.Estimate
            contributorItems["Unassigned"]++
            continue
        }
        
        // Distribute points equally among owners
        pointsPerOwner := item.Estimate / float64(len(item.Owners))
        for _, owner := range item.Owners {
            contributorPoints[owner] += pointsPerOwner
            contributorItems[owner]++
        }
    }
    
    // Sort contributors by points
    type contributorStat struct {
        name       string
        points     float64
        itemCount  int
    }
    
    var stats []contributorStat
    for name, points := range contributorPoints {
        stats = append(stats, contributorStat{
            name:      name,
            points:    points,
            itemCount: contributorItems[name],
        })
    }
    
    // Sort by points in descending order
    sort.Slice(stats, func(i, j int) bool {
        return stats[i].points > stats[j].points
    })
    
    // Generate report string
    report := "Story Points by Contributor:\n\n"
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