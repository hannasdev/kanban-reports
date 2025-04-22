package dateutil

import (
	"fmt"
	"time"
)

// GetStartOfPeriod returns the start date of a period (week or month)
func GetStartOfPeriod(date time.Time, periodType string) time.Time {
    if periodType == "week" {
        // Get start of week (Sunday)
        weekday := date.Weekday()
        return date.AddDate(0, 0, -int(weekday))
    } else {
        // Get start of month
        return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
    }
}

// FormatPeriod formats a date according to period type (week or month)
func FormatPeriod(date time.Time, periodType string) string {
    if periodType == "week" {
        // Format as ISO week: 2024-W02
        year, week := date.ISOWeek()
        return fmt.Sprintf("%d-W%02d", year, week)
    } else {
        // Format as year-month: 2024-01
        return date.Format("2006-01")
    }
}