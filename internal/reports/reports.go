package reports

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/hannasdev/kanban-reports/internal/models"
)

// ReportType defines the type of report to generate
type ReportType string

const (
	// ReportTypeContributor generates report by contributor
	ReportTypeContributor ReportType = "contributor"
	// ReportTypeEpic generates report by epic
	ReportTypeEpic ReportType = "epic"
	// ReportTypeProductArea generates report by product area
	ReportTypeProductArea ReportType = "product-area"
	// ReportTypeTeam generates report by team
	ReportTypeTeam ReportType = "team"
)

// AdHocFilterType defines how to handle ad-hoc requests
type AdHocFilterType string

const (
	// AdHocFilterInclude includes all items (default)
	AdHocFilterInclude AdHocFilterType = "include"
	// AdHocFilterExclude excludes ad-hoc requests
	AdHocFilterExclude AdHocFilterType = "exclude"
	// AdHocFilterOnly shows only ad-hoc requests
	AdHocFilterOnly AdHocFilterType = "only"
)

// Reporter handles generation of different reports
type Reporter struct {
	items []models.KanbanItem
	adHocFilter AdHocFilterType
}

// NewReporter creates a new reporter with the given items
func NewReporter(items []models.KanbanItem) *Reporter {
	return &Reporter{
		items: items,
		adHocFilter: AdHocFilterInclude,
	}
}

// WithAdHocFilter sets the ad-hoc request filter
func (r *Reporter) WithAdHocFilter(filter AdHocFilterType) *Reporter {
	r.adHocFilter = filter
	return r
}

// GenerateReport generates a report based on the specified type and time period
func (r *Reporter) GenerateReport(reportType ReportType, startDate, endDate time.Time) (string, error) {
	// Filter items by completion date within range
	filteredItems := r.filterItemsByDateRange(startDate, endDate)
	
	if len(filteredItems) == 0 {
		return "No items completed in the specified date range.", nil
	}

	// Generate appropriate report based on type
	switch reportType {
	case ReportTypeContributor:
		return r.generateContributorReport(filteredItems)
	case ReportTypeEpic:
		return r.generateEpicReport(filteredItems)
	case ReportTypeProductArea:
		return r.generateProductAreaReport(filteredItems)
	case ReportTypeTeam:
		return r.generateTeamReport(filteredItems)
	default:
		return "", fmt.Errorf("unknown report type: %s", reportType)
	}
}

// filterItemsByDateRange returns items completed within the given date range
func (r *Reporter) filterItemsByDateRange(startDate, endDate time.Time) []models.KanbanItem {
	var filtered []models.KanbanItem
	
	for _, item := range r.items {
		// Only include completed items
		if !item.IsCompleted || item.CompletedAt.IsZero() {
			continue
		}
		
		// Check if completion date is within range
		if (startDate.IsZero() || !item.CompletedAt.Before(startDate)) &&
		   (endDate.IsZero() || !item.CompletedAt.After(endDate)) {
			
			// Apply ad-hoc request filter
			isAdHoc := r.isAdHocRequest(item)
			
			if (r.adHocFilter == AdHocFilterInclude) ||
			   (r.adHocFilter == AdHocFilterExclude && !isAdHoc) ||
			   (r.adHocFilter == AdHocFilterOnly && isAdHoc) {
				filtered = append(filtered, item)
			}
		}
	}
	
	return filtered
}

// isAdHocRequest checks if an item is an ad-hoc request (has "ad-hoc-request" label)
func (r *Reporter) isAdHocRequest(item models.KanbanItem) bool {
	for _, label := range item.Labels {
		if strings.ToLower(label) == "ad-hoc-request" {
			return true
		}
	}
	return false
}

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